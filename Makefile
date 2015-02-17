APP_NAME = consul-alerts
VERSION = latest

all: clean build

clean:
	@echo "--> Cleaning build"
	@rm -rf ./build

prepare:
	@mkdir -p build/bin/{linux-386,linux-amd64,darwin-amd64}
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
	@go build -o build/bin/linux-386/${VERSION}/${APP_NAME} -v .

build-all: test
	@echo "--> Building all application"
	@echo "... linux-386"
	@GOOS=linux GOARCH=386 go build -o build/bin/linux-386/${VERSION}/${APP_NAME} -v .
	@echo "... linux-amd64"
	@GOOS=linux GOARCH=amd64 go build -o build/bin/linux-amd64/${VERSION}/${APP_NAME} -v .
	@echo "... darwin-amd64"
	@GOOS=darwin GOARCH=amd64 go build -o build/bin/darwin-amd64/${VERSION}/${APP_NAME} -v .

package: build-all
	@echo "--> Packaging application"
	@tar cf build/tar/${APP_NAME}-${VERSION}-linux-386.tar -C build/bin/linux-386/${VERSION} ${APP_NAME}
	@tar cf build/tar/${APP_NAME}-${VERSION}-linux-amd64.tar -C build/bin/linux-amd64/${VERSION} ${APP_NAME}
	@tar cf build/tar/${APP_NAME}-${VERSION}-darwin-amd64.tar -C build/bin/darwin-amd64/${VERSION} ${APP_NAME}

release: package
ifeq ($(VERSION) , latest)
	@echo "--> Removing Latest Version"
	@curl -s -X DELETE -u ${ACCESS_KEY} https://api.bintray.com/packages/darkcrux/generic/${APP_NAME}/versions/${VERSION}
	@echo
endif
	@echo "--> Releasing version: ${VERSION}"
	@curl -s -T "build/tar/${APP_NAME}-${VERSION}-linux-386.tar" -u "${ACCESS_KEY}" "https://api.bintray.com/content/darkcrux/generic/${APP_NAME}/${VERSION}/${APP_NAME}-${VERSION}-linux-386.tar"
	@echo "... linux-386"
	@curl -s -T "build/tar/${APP_NAME}-${VERSION}-linux-amd64.tar" -u "${ACCESS_KEY}" "https://api.bintray.com/content/darkcrux/generic/${APP_NAME}/${VERSION}/${APP_NAME}-${VERSION}-linux-amd64.tar"
	@echo "... linux-amd64"
	@curl -s -T "build/tar/${APP_NAME}-${VERSION}-darwin-amd64.tar" -u "${ACCESS_KEY}" "https://api.bintray.com/content/darkcrux/generic/${APP_NAME}/${VERSION}/${APP_NAME}-${VERSION}-darwin-amd64.tar"
	@echo "... darwin-amd64"
	@echo "--> Publishing version ${VERSION}"
	@curl -s -X POST -u ${ACCESS_KEY} https://api.bintray.com/content/darkcrux/generic/${APP_NAME}/${VERSION}/publish
	@echo 
