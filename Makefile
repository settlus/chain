#!/usr/bin/make -f

BUILDDIR ?= $(CURDIR)/build

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%h')
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))
DOCKER := $(shell which docker)

# don't override user values
ifeq (,$(VERSION))
  VERSION := $(shell git describe --exact-match 2>/dev/null)
  # if VERSION is empty, then populate it with branch's name and raw commit hash
  ifeq (,$(VERSION))
    VERSION := $(BRANCH)-$(COMMIT)
  endif
endif

ldflags = -w -s
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

build_tags = $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

whitespace := $(subst ,, )
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))

# process linker flags
ldflags += -X github.com/cosmos/cosmos-sdk/version.Name=settlus \
	-X github.com/cosmos/cosmos-sdk/version.AppName=settlus \
	-X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
	-X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
	-X github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)' -trimpath

.PHONY: all settlusd test lint clean proto-gen

all: settlusd

settlusd: go.sum
	@go build -mod=readonly $(BUILD_FLAGS) -o $(BUILDDIR)/settlusd ./cmd/settlusd

test:
	go test -mod=readonly ./...

lint:
	@golangci-lint run ./...

clean:
	@rm -rf $(BUILDDIR)

###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=v0.7
protoImageName=tendermintdev/sdk-proto-gen:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace --user 0 $(protoImageName)

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
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint --error-format=json

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

proto-clean:
	rm -rf $(SWAGGER_DIR) tmp-swagger-gen

###############################################################################
###                                Localnet                                 ###
###############################################################################

LOCALNET_SETUP_FILE=docker-compose.yml
LOCALNET_DOCKER_TAG=settlus/localnet

localnet-build:
	docker build --build-arg GITHUB_TOKEN=${GITHUB_TOKEN} --tag $(LOCALNET_DOCKER_TAG) $(CURDIR)/.

localnet-stop:
	docker-compose -f $(LOCALNET_SETUP_FILE) down -v
	rm -rf $(CURDIR)/.testnets

localnet-start: localnet-stop
	make settlusd
	$(BUILDDIR)/settlusd testnet init-files --keyring-backend test --starting-ip-address 192.168.11.2
	docker-compose -f $(LOCALNET_SETUP_FILE) up -d
