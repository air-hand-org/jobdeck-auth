#!/bin/sh
set -eu

compose_file="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)/compose.yaml"
project_name="jobdeck_auth_server_test_$$"
tmp_jwks="$(mktemp)"

: "${AUTH_IMAGE:?AUTH_IMAGE must be set}"

cleanup() {
    docker compose --project-name "$project_name" --file "$compose_file" down --volumes --remove-orphans
    rm -f "$tmp_jwks"
}
trap cleanup EXIT INT TERM

docker compose --project-name "$project_name" --file "$compose_file" up --detach

auth_port="$(docker compose --project-name "$project_name" --file "$compose_file" port auth 8080 | sed 's/.*://')"

attempt=1
while ! curl --fail --silent --show-error "http://127.0.0.1:${auth_port}/.well-known/jwks.json" > "$tmp_jwks"; do
    if [ "$attempt" -ge 30 ]; then
        docker compose --project-name "$project_name" --file "$compose_file" logs >&2
        exit 1
    fi
    attempt=$((attempt + 1))
    sleep 1
done

grep -q '"kid":"auth-server-integration-test"' "$tmp_jwks"
docker compose --project-name "$project_name" --file "$compose_file" run --rm grpcurl \
    -plaintext -d '{"service":""}' auth:8080 grpc.health.v1.Health/Check \
    | grep -q 'SERVING'

if rpc_output="$(docker compose --project-name "$project_name" --file "$compose_file" run --rm grpcurl \
    -plaintext -d '{"clientId":"unknown","clientSecret":"unknown","audience":"jobdeck-job-api"}' \
    auth:8080 jobdeck.auth.v1.AuthService/AuthenticateLocalClient 2>&1)"; then
    echo "auth RPC unexpectedly accepted invalid credentials" >&2
    exit 1
fi
printf '%s\n' "$rpc_output" | grep -q 'Unauthenticated'

echo "auth server integration test passed"
