APP_NAME = consul-alerts
VERSION = latest
BUILD_ARCHS=linux-386 linux-amd64 darwin-amd64 freebsd-amd64

all: clean build

clean:
	@echo "--> Cleaning build"
	@rm -rf ./build

prepare:
	@for arch in ${BUILD_ARCHS}; do \
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

build: test
	@echo "--> Building local application"
	@go build -o build/bin/`uname -s`-`uname -p`/${VERSION}/${APP_NAME} -v .

build-all: test
	@echo "--> Building all application"
	@for arch in ${BUILD_ARCHS}; do \
		echo "... $${arch}"; \
		GOOS=`echo $${arch} | cut -d '-' -f 1` \
		GOARCH=`echo $${arch} | cut -d '-' -f 2` \
		go build -o build/bin/$${arch}/${VERSION}/${APP_NAME} -v . ; \
	done

package: build-all
	@echo "--> Packaging application"
	@for arch in ${BUILD_ARCHS}; do \
		tar czf build/tar/${APP_NAME}-${VERSION}-$${arch}.tgz -C build/bin/$${arch}/${VERSION} ${APP_NAME} ; \
	done

release: package
ifeq ($(VERSION) , latest)
	@echo "--> Removing Latest Version"
	@curl -s -X DELETE -u ${ACCESS_KEY} https://api.bintray.com/packages/darkcrux/generic/${APP_NAME}/versions/${VERSION}
	@echo
endif
	@echo "--> Releasing version: ${VERSION}"
	@for arch in ${BUILD_ARCHS}; do \
		curl -s -T "build/tar/${APP_NAME}-${VERSION}-$${arch}.tgz" -u "${ACCESS_KEY}" "https://api.bintray.com/content/darkcrux/generic/${APP_NAME}/${VERSION}/${APP_NAME}-${VERSION}-$${arch}.tgz"; \
		echo "... $${arch}"; \
	done
	@echo "--> Publishing version ${VERSION}"
	@curl -s -X POST -u ${ACCESS_KEY} https://api.bintray.com/content/darkcrux/generic/${APP_NAME}/${VERSION}/publish
	@echo
	@echo "Github Release"
	@gh-release create uchiru/consul-alerts ${VERSION}
