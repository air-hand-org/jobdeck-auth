# syntax=docker/dockerfile:1

FROM golang:1.26.4-bookworm AS build
WORKDIR /src

# Cache module downloads independently from application source changes.
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod,id=jobdeck-go-mod \
    go mod download

# Keep unrelated repository changes out of the build cache key.
COPY cmd ./cmd
COPY internal ./internal
RUN --mount=type=cache,target=/go/pkg/mod,id=jobdeck-go-mod \
    --mount=type=cache,target=/root/.cache/go-build,id=jobdeck-go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/auth ./cmd/server
RUN --mount=type=cache,target=/go/pkg/mod,id=jobdeck-go-mod \
    --mount=type=cache,target=/root/.cache/go-build,id=jobdeck-go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/auth-migrator ./cmd/migrator
RUN --mount=type=cache,target=/go/pkg/mod,id=jobdeck-go-mod \
    --mount=type=cache,target=/root/.cache/go-build,id=jobdeck-go-build \
    CGO_ENABLED=0 go build -trimpath -ldflags='-s -w' -o /out/auth-entrypoint ./cmd/entrypoint

FROM gcr.io/distroless/base-debian12:nonroot
COPY LICENSE /LICENSE
COPY third_party/licenses/linux-amd64 /third_party
COPY --from=build /out/auth /auth
COPY --from=build /out/auth-migrator /auth-migrator
COPY --from=build /out/auth-entrypoint /auth-entrypoint
USER nonroot:nonroot
ENTRYPOINT ["/auth-entrypoint"]
