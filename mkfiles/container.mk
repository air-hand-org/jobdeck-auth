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

.PHONY: auth-proxy-build
auth-proxy-build:
	@cd $(CONTAINER_MK_DIR)../ \
	&& docker build -t jobdeck-auth-proxy:test ./containers/auth-proxy

.PHONY: auth-proxy-integration-test
auth-proxy-integration-test: container-build auth-proxy-build
	@cd $(CONTAINER_MK_DIR)../ \
	&& AUTH_IMAGE=jobdeck-auth:test AUTH_PROXY_IMAGE=jobdeck-auth-proxy:test sh ./tests/auth-proxy/run.sh
