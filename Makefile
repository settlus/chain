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

.PHONY: all chaind test lint clean proto-gen

all: build

build: go.sum
	@go build -mod=readonly $(BUILD_FLAGS) -o $(BUILDDIR)/settlusd ./cmd/settlusd

test:
	go test -mod=readonly ./...

lint:
	@golangci-lint run ./...

clean:
	@rm -rf $(BUILDDIR)

# Protobuf generation

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

# Local network command
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
