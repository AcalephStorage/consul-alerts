APP_NAME = consul-alerts
VERSION = latest
BUILD_ARCHES=linux-386 linux-amd64 darwin-amd64 freebsd-amd64

all: clean build

clean:
	@echo "--> Cleaning build"
	@rm -rf ./build

prepare:
	@for arch in ${BUILD_ARCHES}; do \
		mkdir -p build/bin/$${arch}; \
	done
	@mkdir -p build/test
	@mkdir -p build/doc
	@mkdir -p build/tar

format:
	@echo "--> Formatting source code"
	@go fmt ./...

test: prepare format
	@echo "--> Testing application"
	@go test -outputdir build/test ./...

OS   := $(shell uname -s | tr [:upper:] [:lower:])
ARCH := $(shell uname -p | sed -e "s/x86_64/amd64/" | tr [:upper:] [:lower:])

build:
	@echo "--> Building local application"
	@go build -o build/bin/$(OS)-$(ARCH)/${VERSION}/${APP_NAME} -v .

build-all:
	@echo "--> Building all application"
	@for arch in ${BUILD_ARCHES}; do \
		echo "... $${arch}"; \
		GOOS=`echo $${arch} | cut -d '-' -f 1` \
		GOARCH=`echo $${arch} | cut -d '-' -f 2` \
		go build -o build/bin/$${arch}/${VERSION}/${APP_NAME} . ; \
	done

package: build-all
	@echo "--> Packaging application"
	@for arch in ${BUILD_ARCHES}; do \
		tar czf build/tar/${APP_NAME}-${VERSION}-$${arch}.tgz -C build/bin/$${arch}/${VERSION} ${APP_NAME} ; \
	done

release: package
	@echo "--> Releasing version: ${VERSION}"
ifneq ($(VERSION),latest)
	@echo "Github Release"
	@gh-release create YOwatari/consul-alerts ${VERSION}
endif
