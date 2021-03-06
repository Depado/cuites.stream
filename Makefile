.DEFAULT_GOAL := build

export CGO_ENABLED=0
BINARY=cuites.stream
VERSION=$(shell git describe --abbrev=0 --tags 2> /dev/null || echo "0.1.0")
BUILD=$(shell git rev-parse HEAD 2> /dev/null || echo "undefined")
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Build=$(BUILD)"

.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build
	go build -o $(BINARY) $(LDFLAGS)

.PHONY: armbuild
armbuild: ## Build a version for ARM
	GOARCH=arm GOARM=5 go build -o $(BINARY) $(LDFLAGS)

.PHONY: install
install: ## Build and install
	go install $(LDFLAGS) 

.PHONY: test
test: ## Run the test suite
	CGO_ENABLED=1 go test ./...

.PHONY: clean
clean: ## Remove the binary
	if [ -f $(BINARY) ] ; then rm $(BINARY) ; fi