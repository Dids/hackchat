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

BINARY_VERSION?=0.0.1
BINARY_OUTPUT?=hackchat
EXTRA_FLAGS?=-mod=vendor

# Run commands in a -c flag 
.SHELLFLAGS = -c

.SILENT: ;               # no need for @
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
	go install -v $(EXTRA_FLAGS) -ldflags "-X main.Version=$(BINARY_VERSION)"

uninstall: ## Uninstall the binary
	rm -f $(GOPATH)/bin/$(BINARY_OUTPUT)

build: ## Build the binary
	go build -v $(EXTRA_FLAGS) -ldflags "-X main.Version=$(BINARY_VERSION)" -o $(BINARY_OUTPUT)

test: ## Run unit tests
	go test -v $(EXTRA_FLAGS) -race -coverprofile=coverage.txt -covermode=atomic ./...

clean: ## Cleanup temporary files
	go clean
	rm -f $(BINARY_OUTPUT)

deps: ## Install dependencies
	go build -v $(EXTRA_FLAGS) ./...
	go install github.com/cosmtrek/air@v1.15.1

upgrade: ## Upgrade dependencies
	go get -d
	go mod vendor
	go mod tidy

tidy: ## Cleanup Go modules
	go mod tidy

dev: ## Start in hot reload mode
	rm -fr *.log .tmp
	air -c .air.toml
