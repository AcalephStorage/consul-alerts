APP_NAME = consul-alerts

all: clean build

clean:
	@echo "--> Cleaning build"
	@rm -rf ./build

prepare:
	@mkdir -p build/bin
	@mkdir -p build/test
	@mkdir -p build/doc
	@mkdir -p build/tar

format:
	@echo "--> Formatting source code"
	@go fmt ./...

test: prepare format
	@echo "--> Testing application"
	@go test -outputdir build/test ./...

OS   := $(shell uname -s)
ARCH := $(shell uname -p)

build:
	@echo "--> Building local application"
	@go build -o build/bin/$(OS)-$(ARCH)/${VERSION}/${APP_NAME} -v .

build-all:
	goreleaser build

export GITHUB_TOKEN :=

release:
	@echo "Github Release"
	goreleaser
