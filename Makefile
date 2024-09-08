#!/usr/bin/make -f

# Git information
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%h')
VERSION ?= $(shell echo $(shell git describe --tags --always) | sed 's/^v//')

# Path and directory variables
MKFILE_PATH := $(abspath $(lastword $(MAKEFILE_LIST)))
CURRENT_DIR := $(notdir $(patsubst %/,%,$(dir $(MKFILE_PATH))))

# Docker configuration
DOCKER := $(shell which docker)

# Binary configuration
SETTLUS_BINARY := settlusd

.PHONY: all

all: build

###############################################################################
###                                  Build                                  ###
###############################################################################

# Build tags processing
build_tags := netgo
ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

# Build tags to linker flags
build_tags_comma_sep := $(subst $(subst ,, ),,,$(build_tags))

# Linker flags setup
ldflags := -X github.com/cosmos/cosmos-sdk/version.Name=settlus \
           -X github.com/cosmos/cosmos-sdk/version.AppName=$(SETTLUS_BINARY) \
           -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
           -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
           -X github.com/settlus/chain/tools/interop-node/version.Name=interop-node \
           -X github.com/settlus/chain/tools/interop-node/version.Version=$(VERSION) \
           -X github.com/settlus/chain/tools/interop-node/version.Commit=$(COMMIT) \
           -X github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep) \
           -X github.com/settlus/chain/tools/interop-node/version.BuildTags=$(build_tags_comma_sep)

# Additional linker flags
ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

# Final build flags setup
BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif
ifneq (,$(findstring nooptimization,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -gcflags "all=-N -l"
endif

# Build targets
BUILD_TARGETS := build install
BUILDDIR ?= $(CURDIR)/build

build: BUILD_ARGS=-o $(BUILDDIR)/
$(BUILD_TARGETS): go.sum $(BUILDDIR)/
	@echo "Building..."
	@go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./... || (echo "Build failed"; exit 1)

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

test:
	@echo "Running tests..."
	@go test -mod=readonly $(shell go list ./... | grep -v tests/e2e)|| (echo "Tests failed"; exit 1)

test-e2e:
	@echo "Running e2e tests..."
	@go test ./tests/e2e -v || (echo "Tests failed"; exit 1)

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILDDIR)

.PHONY: build test clean

###############################################################################
###                                Linting                                  ###
###############################################################################

lint:
	@golangci-lint run ./...

.PHONY: lint

###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=0.11.6
protoImageName=ghcr.io/cosmos/proto-builder:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace --user 0 $(protoImageName)

protoLintVer=0.44.0
protoLinterImage=yoheimuta/protolint
protoLinter=$(DOCKER) run --rm -v "$(CURDIR):/workspace" --workdir /workspace --user 0 $(protoLinterImage):$(protoLintVer)

SWAGGER_DIR=./swagger-proto
THIRD_PARTY_DIR=$(SWAGGER_DIR)/third_party
DEPS_COSMOS_SDK_VERSION := $(shell cat go.sum | grep 'github.com/settlus/cosmos-sdk' | grep -v -e 'go.mod' | tail -n 1 | awk '{ print $$2; }')
DEPS_IBC_GO_VERSION := $(shell cat go.sum | grep 'github.com/cosmos/ibc-go' | grep -v -e 'go.mod' | tail -n 1 | awk '{ print $$2; }')
DEPS_COSMOS_PROTO := $(shell cat go.sum | grep 'github.com/cosmos/cosmos-proto' | grep -v -e 'go.mod' | tail -n 1 | awk '{ print $$2; }')
DEPS_COSMOS_GOGOPROTO := $(shell cat go.sum | grep 'github.com/cosmos/gogoproto' | grep -v -e 'go.mod' | tail -n 1 | awk '{ print $$2; }')
DEPS_COSMOS_ICS23 := go/$(shell cat go.sum | grep 'github.com/cosmos/ics23/go' | grep -v -e 'go.mod' | tail -n 1 | awk '{ print $$2; }')

proto-all: proto-format proto-lint proto-gen proto-swagger-gen

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./scripts/protocgen.sh

proto-format:
	@echo "Formatting Protobuf files"
	$(protoImage) find ./ -name *.proto -exec clang-format -i {} \;

proto-lint:
	@echo "Linting Protobuf files"
	@$(protoImage) buf lint --error-format=json
	@$(protoLinter) lint ./proto

proto-swagger-gen:
	@echo "Downloading Protobuf dependencies"
	@make proto-download-deps
	@echo "Generating Protobuf Swagger"
	$(protoImage) sh ./scripts/protoc-swagger-gen.sh

proto-download-deps:
	mkdir -p "$(THIRD_PARTY_DIR)/cosmos_tmp" && \
	cd "$(THIRD_PARTY_DIR)/cosmos_tmp" && \
	git init && \
	git remote add origin "https://github.com/settlus/cosmos-sdk.git" && \
	git config core.sparseCheckout true && \
	printf "proto\nthird_party\n" > .git/info/sparse-checkout && \
	git pull origin "$(DEPS_COSMOS_SDK_VERSION)" && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/cosmos_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/ibc_tmp" && \
	cd "$(THIRD_PARTY_DIR)/ibc_tmp" && \
	git init && \
	git remote add origin "https://github.com/cosmos/ibc-go.git" && \
	git config core.sparseCheckout true && \
	printf "proto\n" > .git/info/sparse-checkout && \
	git pull origin "$(DEPS_IBC_GO_VERSION)" && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/ibc_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/cosmos_proto_tmp" && \
	cd "$(THIRD_PARTY_DIR)/cosmos_proto_tmp" && \
	git init && \
	git remote add origin "https://github.com/cosmos/cosmos-proto.git" && \
	git config core.sparseCheckout true && \
	printf "proto\n" > .git/info/sparse-checkout && \
	git pull origin "$(DEPS_COSMOS_PROTO_VERSION)" && \
	rm -f ./proto/buf.* && \
	mv ./proto/* ..
	rm -rf "$(THIRD_PARTY_DIR)/cosmos_proto_tmp"

	mkdir -p "$(THIRD_PARTY_DIR)/gogoproto" && \
	curl -SSL "https://raw.githubusercontent.com/cosmos/gogoproto/$(DEPS_COSMOS_GOGOPROTO)/gogoproto/gogo.proto" > "$(THIRD_PARTY_DIR)/gogoproto/gogo.proto"

	mkdir -p "$(THIRD_PARTY_DIR)/google/api" && \
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/annotations.proto > "$(THIRD_PARTY_DIR)/google/api/annotations.proto"
	curl -sSL https://raw.githubusercontent.com/googleapis/googleapis/master/google/api/http.proto > "$(THIRD_PARTY_DIR)/google/api/http.proto"

	mkdir -p "$(THIRD_PARTY_DIR)/cosmos/ics23/v1" && \
	curl -sSL "https://raw.githubusercontent.com/cosmos/ics23/$(DEPS_COSMOS_ICS23)/proto/cosmos/ics23/v1/proofs.proto" > "$(THIRD_PARTY_DIR)/cosmos/ics23/v1/proofs.proto"

proto-clean:
	rm -rf $(SWAGGER_DIR) tmp-swagger-gen

.PHONY: proto-all proto-gen proto-format proto-lint proto-swagger-gen proto-download-deps proto-clean

###############################################################################
###                                Localnet                                 ###
###############################################################################

LOCALNET_SETUP_FILE=docker-compose.yml
LOCALNET_DOCKER_TAG=settlus/localnet

localnet-build:
	docker build --tag $(LOCALNET_DOCKER_TAG) $(CURDIR)/.

localnet-stop:
	docker-compose -f $(LOCALNET_SETUP_FILE) down -v
	rm -rf $(CURDIR)/.testnets

localnet-start: localnet-stop
	make build
	$(BUILDDIR)/settlusd testnet init-files --keyring-backend test --starting-ip-address 192.168.11.2
	docker-compose -f $(LOCALNET_SETUP_FILE) up -d

.PHONY: localnet-build localnet-start localnet-stop

###############################################################################
###                                Releasing                                ###
###############################################################################

PACKAGE_NAME:=github.com/settlus/chain
GOLANG_CROSS_VERSION  = v1.22
GOPATH ?= '$(HOME)/go'
release-dry-run:
	docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		-e LDFLAGS="$(ldflags)" \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-v ${GOPATH}/pkg:/go/pkg \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		--clean --skip=validate --skip=publish --snapshot

release:
	@if [ ! -f ".release-env" ]; then \
		echo "\033[91m.release-env is required for release\033[0m";\
		exit 1;\
	fi
	docker run \
		--rm \
		--privileged \
		-e CGO_ENABLED=1 \
		-e LDFLAGS="$(ldflags)" \
		--env-file .release-env \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --clean --skip=validate

.PHONY: release-dry-run release

###############################################################################
###                                    MISC                                 ###
###############################################################################
help:
	@echo "Available targets:"
	@echo "  build              - Build the project"
	@echo "  test               - Run tests"
	@echo "  e2e-test           - Run end-to-end tests"
	@echo "  clean              - Clean build artifacts"
	@echo "  lint               - Run linter"
	@echo "  proto-all          - Run all protobuf-related tasks"
	@echo "  proto-gen          - Generate Protobuf files"
	@echo "  proto-format       - Format Protobuf files"
	@echo "  proto-lint         - Lint Protobuf files"
	@echo "  proto-swagger-gen  - Generate Protobuf Swagger"
	@echo "  release-dry-run    - Dry run release"
	@echo "  release            - Release"
