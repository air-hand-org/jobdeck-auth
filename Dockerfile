# syntax=docker/dockerfile:1

FROM golang:1.26.4-bookworm AS build
WORKDIR /src

# Cache module downloads independently from application source changes.
COPY go.mod go.sum ./
RUN --mount=type=secret,id=github_contents_token \
    --mount=type=cache,target=/go/pkg/mod,id=jobdeck-go-mod \
    GOPRIVATE=github.com/air-hand-org/api-contracts \
    GIT_CONFIG_COUNT=1 \
    GIT_CONFIG_KEY_0="url.https://x-access-token:$(cat /run/secrets/github_contents_token)@github.com/.insteadOf" \
    GIT_CONFIG_VALUE_0="https://github.com/" \
    go mod download

# The image runs only these commands. Keep unrelated repository changes out of
# the build cache key.
COPY cmd/entrypoint ./cmd/entrypoint
COPY cmd/migrator ./cmd/migrator
COPY cmd/server ./cmd/server
COPY internal ./internal
RUN --mount=type=secret,id=github_contents_token \
    --mount=type=cache,target=/go/pkg/mod,id=jobdeck-go-mod \
    --mount=type=cache,target=/root/.cache/go-build,id=jobdeck-go-build \
    GOPRIVATE=github.com/air-hand-org/api-contracts \
    GIT_CONFIG_COUNT=1 \
    GIT_CONFIG_KEY_0="url.https://x-access-token:$(cat /run/secrets/github_contents_token)@github.com/.insteadOf" \
    GIT_CONFIG_VALUE_0="https://github.com/" \
    CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/auth ./cmd/server
RUN --mount=type=cache,target=/go/pkg/mod,id=jobdeck-go-mod \
    --mount=type=cache,target=/root/.cache/go-build,id=jobdeck-go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/auth-migrator ./cmd/migrator
RUN --mount=type=cache,target=/go/pkg/mod,id=jobdeck-go-mod \
    --mount=type=cache,target=/root/.cache/go-build,id=jobdeck-go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/auth-entrypoint ./cmd/entrypoint

FROM gcr.io/distroless/base-debian12:nonroot
COPY --from=build /out/auth /auth
COPY --from=build /out/auth-migrator /auth-migrator
COPY --from=build /out/auth-entrypoint /auth-entrypoint
USER nonroot:nonroot
ENTRYPOINT ["/auth-entrypoint"]
