#!/bin/sh
set -eu

: "${AUTH_GRPC_HOST:?AUTH_GRPC_HOST must be set}"
: "${AUTH_GRPC_PORT:?AUTH_GRPC_PORT must be set}"
: "${AUTH_JWKS_HOST:?AUTH_JWKS_HOST must be set}"
: "${AUTH_JWKS_PORT:?AUTH_JWKS_PORT must be set}"

PROXY_PORT="${PROXY_PORT:-8080}"

for port in "$PROXY_PORT" "$AUTH_GRPC_PORT" "$AUTH_JWKS_PORT"; do
    case "$port" in
        *[!0-9]* | '')
            echo "proxy and upstream ports must be numeric" >&2
            exit 1
            ;;
    esac
done

export PROXY_PORT

envsubst \
    '${PROXY_PORT} ${AUTH_GRPC_HOST} ${AUTH_GRPC_PORT} ${AUTH_JWKS_HOST} ${AUTH_JWKS_PORT}' \
    < /etc/envoy/envoy.yaml.template \
    > /tmp/envoy.yaml

exec envoy -c /tmp/envoy.yaml "$@"
