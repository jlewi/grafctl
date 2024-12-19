ROOT := $(shell git rev-parse --show-toplevel)
PROJECT := chat-lewi

GIT_SHA := $(shell git rev-parse HEAD)
GIT_SHA_SHORT := $(shell git rev-parse --short HEAD)
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
# This gives an error no names if there are no tags but this gets fixed as soon
# as you create releases.
VERSION := $(shell git describe --tags)-$(GIT_SHA_SHORT)
LDFLAGS := -s -w \
        -X 'github.com/jlewi/grafctl/pkg/version.Date=$(DATE)' \
        -X 'github.com/jlewi/grafctl/pkg/version.Version=$(subst v,,$(VERSION))' \
        -X 'github.com/jlewi/grafctl/pkg/version.Commit=$(GIT_SHA)'

build: build-dir
	CGO_ENABLED=0 go build -o .build/grafctl -ldflags="$(LDFLAGS)" github.com/jlewi/grafctl

build-dir:
	mkdir -p .build

tidy:
	gofmt -s -w .
	goimports -w .
	

lint:
	# golangci-lint automatically searches up the root tree for configuration files.
	golangci-lint run

test:	
	GITHUB_ACTIONS=true go test -v ./...