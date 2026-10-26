all: setup build

include ./mkfiles/*.mk

.PHONY: setup
setup:
	@aqua i && mise i

.PHONY: generate
generate:
	@sqlc generate

.PHONY: build
build: generate licenses
	@go generate ./... && go mod tidy && go build ./...

# Match the platform currently published by _build-images.yaml.
.PHONY: licenses
licenses:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go-licenses check ./cmd/server ./cmd/migrator ./cmd/entrypoint --ignore github.com/air-hand-org/jobdeck-auth --disallowed_types=forbidden,unknown
	rm -rf third_party/licenses/linux-amd64
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go-licenses save ./cmd/server ./cmd/migrator ./cmd/entrypoint --ignore github.com/air-hand-org/jobdeck-auth --save_path third_party/licenses/linux-amd64

.PHONY: fmt
fmt:
	@golangci-lint fmt ./...

.PHONY: lint
lint:
	@golangci-lint fmt ./... --diff && golangci-lint run ./...

.PHONY: check-clean
check-clean:
	@git status --short
	@test -z "$$(git status --porcelain)"

.PHONY: test
test:
	@go test -cover ./...
