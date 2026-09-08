all: setup build

include ./mkfiles/*.mk

.PHONY: setup
setup:
	@aqua i && mise i

.PHONY: generate
generate:
	@sqlc generate

.PHONY: build
build: generate
	@go generate ./... && go mod tidy && go build ./...

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
