CONTAINER_MK := $(lastword $(MAKEFILE_LIST))
CONTAINER_MK_DIR := $(dir $(abspath $(CONTAINER_MK)))
GITHUB_CONTENTS_TOKEN ?= $(shell gh auth token)

.PHONY: container-build
container-build:
	@cd $(CONTAINER_MK_DIR)../ \
	&& docker build \
		-t jobdeck-auth:test \
		--secret id=github_contents_token,env=GITHUB_CONTENTS_TOKEN \
		. \
	;

.PHONY: auth-server-integration-test
auth-server-integration-test:
	@cd $(CONTAINER_MK_DIR)../ \
	&& AUTH_IMAGE=jobdeck-auth:test sh ./tests/auth-server/run.sh

.PHONY: auth-server-build-and-test
auth-server-build-and-test: container-build auth-server-integration-test
