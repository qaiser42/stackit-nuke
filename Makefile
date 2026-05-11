SHELL := /bin/bash
BIN   := stackit-nuke
PKG   := github.com/qaiser42/stackit-nuke

VERSION ?= $(shell git describe --tags --dirty --always 2>/dev/null || echo "0.0.0-dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "dirty")
BRANCH  ?= $(shell git rev-parse --abbrev-ref HEAD 2>/dev/null || echo "dev")

LDFLAGS := -s -w \
	-X $(PKG)/pkg/common.SUMMARY=$(VERSION) \
	-X $(PKG)/pkg/common.COMMIT=$(COMMIT) \
	-X $(PKG)/pkg/common.BRANCH=$(BRANCH)

.PHONY: build test lint tidy run snapshot docs-serve docs-build clean

build:
	CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN) ./

test:
	go test -timeout 120s -race -cover ./...

lint:
	golangci-lint run

tidy:
	go mod tidy

run: build
	./$(BIN) --help

snapshot:
	goreleaser release --snapshot --clean

docs-serve:
	docker run --rm -p 8000:8000 -v "$(PWD):/docs" squidfunk/mkdocs-material

docs-build:
	docker run --rm -v "$(PWD):/docs" squidfunk/mkdocs-material build

clean:
	rm -rf $(BIN) releases/ public/
