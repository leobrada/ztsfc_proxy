# Project setup
PROJECT_NAME := ztsfc_proxy
PROJECT_REPO := ./cmd/ztsfc_proxy

# Build variables
BUILD_DIR := ./build
BINARY := $(BUILD_DIR)/$(PROJECT_NAME)

# Go related variables.
#GOBASE := $(shell pwd)
#GOBIN := $(GOBASE)/bin
GOPKG := $(.)

# Go files
#GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

# Redirect error output to a file, so we can show it in development mode.
#STDERR := /tmp/.$(PROJECT_NAME)-stderr.txt

# PID file will keep the process id of the server
PID := /tmp/.$(PROJECT_NAME).pid

# Make is verbose in Linux. Make it silent.
#MAKEFLAGS += --silent

## build: Build the binary.
build: go-get go-tidy go-build

## start: Start in development mode.
#start: go-get go-build go-run

## stop: Stop development mode.
#stop:
#	@-touch $(PID)
#	@-kill `cat $(PID)` 2> /dev/null || true#

### test: Run all tests.
#test: go-tes#t

### clean: Clean build files. Runs `go clean` internally.
#clean: go-clean

#go-run:
#	@echo "  >  Running $(PROJECT_NAME)..."
#	@-$(BINARY) 2>&1 & echo $$! > $(PID)
#	@cat $(PID) | sed "/^/s/^/  \>  PID: /"
#	@echo "  >  $(PROJECT_NAME) is running!"

#go-generate:
#	@echo "  >  Generating dependency files..."
#	@GOBIN=$(GOBIN) go generate $(generate)

go-get:
	@echo "> Checking for dependency updates..."
	@go get -v -u all

go-tidy:
	@echo "> Checking if there is any missing dependencies..."
	@go mod tidy

go-build:
	@echo "> Building binary..."
	@go build -o $(BINARY) $(PROJECT_REPO)

#go-install:
#	@echo "  >  Installing dependencies..."
#	@GOBIN=$(GOBIN) go install $(GOPKG)

#go-test:
#	@echo "  >  Testing..."
#	@GOBIN=$(GOBIN) go test ./...

#go-clean:
#	@echo "  >  Cleaning build cache"
#	@GOBIN=$(GOBIN) go clean

#.PHONY: build start stop test clean go-build go-run go-generate go-get go-install go-test go-clean

.PHONY: build go-get go-tidy go-build