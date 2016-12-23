# variable definitions
NAME := inca3
DESC := Configuration grabber
VERSION := $(shell git describe --tags --always --dirty)
GOVERSION := $(shell go version)
BUILDTIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
BUILDDATE := $(shell date -u +"%B %d, %Y")
BUILDER := $(shell echo "`git config user.name` <`git config user.email`>")
PKG_RELEASE ?= 1
PROJECT_URL := "https://github.com/lfkeitel/$(NAME)"
LDFLAGS := -X 'main.version=$(VERSION)' \
			-X 'main.buildTime=$(BUILDTIME)' \
			-X 'main.builder=$(BUILDER)' \
			-X 'main.goversion=$(GOVERSION)'

.PHONY: all doc fmt alltests test coverage benchmark lint vet build dist

all: test build

# development tasks
doc:
	@godoc -http=:6060 -index

fmt:
	@go fmt $$(go list ./src/...)

alltests: test lint vet

test:
	@go test $$(go list ./src/...)

coverage:
	@go test -cover $$(go list ./src/...)

benchmark:
	@echo "Running tests..."
	@go test -bench=. $$(go list ./src/...)

# https://github.com/golang/lint
# go get github.com/golang/lint/golint
lint:
	@golint ./src/...

vet:
	@go vet $$(go list ./src/...)

build:
	GOBIN=$(PWD)/bin go install -v -ldflags "$(LDFLAGS)" ./cmd/inca3

dist: vet all
	@rm -rf ./dist
	@mkdir -p dist/inca
	@cp -R public dist/inca/
	@cp -R scripts dist/inca/

	@cp LICENSE dist/inca/
	@cp README.md dist/inca/

	@mkdir dist/inca/bin
	@cp bin/inca3 dist/inca/bin/inca

	(cd "dist"; tar -cz inca) > "dist/inca-dist-$(VERSION).tar.gz"

	@rm -rf dist/inca