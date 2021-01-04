# Partially based on this article:
# https://le-gall.bzh/post/makefile-based-ci-chain-for-go/

# Set the default shell
SHELL := $(shell which bash)

# Set the env path
ENV = /usr/bin/env

# Enable Go module support
#export GO111MODULE=on
#export GOPROXY=

#export PATH := $(GOPATH)/bin:$(PATH)

TAG_COMMIT := $(shell git rev-list --abbrev-commit --tags --max-count=1)
TAG := $(shell git describe --abbrev=0 --tags ${TAG_COMMIT} 2>/dev/null || true)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")
VERSION := $(TAG:v%=%)
ifneq ($(COMMIT), $(TAG_COMMIT))
	VERSION := $(VERSION)-next-$(COMMIT)-$(DATE)
endif
ifeq ($(VERSION),)
	VERSION := $(COMMIT)-$(DATE)
endif
ifneq ($(shell git status --porcelain),)
	VERSION := $(VERSION)-dirty
endif

# BINARY_VERSION?=0.0.1
BINARY_OUTPUT?=hackchat
EXTRA_FLAGS?=-mod=vendor

# Run commands in a -c flag 
.SHELLFLAGS = -c

#.SILENT: ;               # no need for @
.ONESHELL: ;             # recipes execute in same shell
.NOTPARALLEL: ;          # wait for this target to finish
.EXPORT_ALL_VARIABLES: ; # send all vars to shell

#.PHONY: all install uninstall build test clean deps upgrade tidy dev
.PHONY: all # All targets are accessible for user
.DEFAULT: help # Running Make will run the help target

help: ## Show Help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

#all: deps build

install: ## Install the binary
	go install -v $(EXTRA_FLAGS) -ldflags "-X main.version=$(VERSION)"

uninstall: ## Uninstall the binary
	rm -f $(GOPATH)/bin/$(BINARY_OUTPUT)

build: ## Build the binary
	@echo "Building version $(VERSION)"
	go build -v $(EXTRA_FLAGS) -ldflags "-X main.version=$(VERSION)" -o $(BINARY_OUTPUT)

test: ## Run unit tests
	go test -v $(EXTRA_FLAGS) -race -coverprofile=coverage.txt -covermode=atomic ./...

clean: ## Cleanup temporary files
	go clean
	rm -f $(BINARY_OUTPUT)

deps: ## Install dependencies
	go build -v $(EXTRA_FLAGS) ./...

upgrade: ## Upgrade dependencies
	go get -d
	go mod vendor
	go mod tidy

tidy: ## Cleanup Go modules
	go mod tidy

dev: ## Start in hot reload mode
	go install github.com/cosmtrek/air@v1.15.1
	rm -fr *.log .tmp
	air -c .air.toml
