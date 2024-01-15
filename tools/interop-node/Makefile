#!/usr/bin/make -f

BUILDDIR ?= $(CURDIR)/build

BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
COMMIT := $(shell git log -1 --format='%h')
mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
current_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))

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

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)' -trimpath

.PHONY: all build test lint clean

all: build

build: 
	@go build $(BUILD_FLAGS) -o $(BUILDDIR)/ ./...

test:
	go test -mod=readonly ./...

lint:
	@golangci-lint run ./...

clean:
	@rm -rf $(BUILDDIR)