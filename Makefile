#!/usr/bin/make -f

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%h')
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))
DOCKER := $(shell which docker)
SETTLUS_BINARY = settlusd

.PHONY: all chaind test lint clean proto-gen

all: build

###############################################################################
###                                  Build                                  ###
###############################################################################

VERSION ?= $(shell echo $(shell git describe --tags --always) | sed 's/^v//')

# process build tags

build_tags = netgo

ifeq (cleveldb,$(findstring cleveldb,$(COSMOS_BUILD_OPTIONS)))
  build_tags += gcc
endif
build_tags += $(BUILD_TAGS)
build_tags := $(strip $(build_tags))

# process linker flags

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=settlus \
          -X github.com/cosmos/cosmos-sdk/version.AppName=$(SETTLUS_BINARY) \
          -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
          -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \

# add build tags to linker flags
whitespace := $(subst ,, )
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))
ldflags += -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)"

ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  ldflags += -w -s
endif
ldflags += $(LDFLAGS)
ldflags := $(strip $(ldflags))

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'
# check for nostrip option
ifeq (,$(findstring nostrip,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -trimpath
endif

# check if no optimization option is passed
# used for remote debugging
ifneq (,$(findstring nooptimization,$(COSMOS_BUILD_OPTIONS)))
  BUILD_FLAGS += -gcflags "all=-N -l"
endif

BUILD_TARGETS := build install
BUILDDIR ?= $(CURDIR)/build
BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

build: BUILD_ARGS=-o $(BUILDDIR)/

build-linux-amd64:
	GOOS=linux GOARCH=amd64 CGO_ENABLED="1" LEDGER_ENABLED=false $(MAKE) build

build-linux-arm64:
	GOOS=linux GOARCH=arm64 LEDGER_ENABLED=false $(MAKE) build

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
	go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

test:
	go test -mod=readonly ./...

clean:
	@rm -rf $(BUILDDIR)

.PHONY: build build-linux-amd64 build-linux-arm64

###############################################################################
###                                Linting                                  ###
###############################################################################

lint:
	@golangci-lint run ./...

###############################################################################
###                                Protobuf                                 ###
###############################################################################

protoVer=v0.7
protoImageName=tendermintdev/sdk-proto-gen:$(protoVer)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

proto-all: proto-format proto-lint proto-gen

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./scripts/protocgen.sh

proto-format:
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

proto-lint:
	@$(protoImage) buf lint --error-format=json

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
	$(BUILDDIR)/chaind testnet init-files --keyring-backend test --starting-ip-address 192.168.11.2
	docker-compose -f $(LOCALNET_SETUP_FILE) up -d
