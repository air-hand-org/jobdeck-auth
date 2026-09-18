#!/bin/sh
set -eu

compose_file="$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)/compose.yaml"
project_name="jobdeck_auth_proxy_test_$$"
tmp_jwks="$(mktemp)"

: "${AUTH_IMAGE:?AUTH_IMAGE must be set}"
: "${AUTH_PROXY_IMAGE:?AUTH_PROXY_IMAGE must be set}"

cleanup() {
    docker compose --project-name "$project_name" --file "$compose_file" down --volumes --remove-orphans
    rm -f "$tmp_jwks"
}
trap cleanup EXIT INT TERM

docker compose --project-name "$project_name" --file "$compose_file" up --detach

proxy_port="$(docker compose --project-name "$project_name" --file "$compose_file" port auth-proxy 18080 | sed 's/.*://')"

attempt=1
while ! curl --fail --silent --show-error "http://127.0.0.1:${proxy_port}/.well-known/jwks.json" > "$tmp_jwks"; do
    if [ "$attempt" -ge 30 ]; then
        docker compose --project-name "$project_name" --file "$compose_file" logs >&2
        exit 1
    fi
    attempt=$((attempt + 1))
    sleep 1
done

grep -q '"kid":"auth-proxy-integration-test"' "$tmp_jwks"
docker compose --project-name "$project_name" --file "$compose_file" run --rm grpcurl \
    -plaintext -d '{"service":""}' auth-proxy:18080 grpc.health.v1.Health/Check \
    | grep -q 'SERVING'

echo "auth proxy integration test passed"
