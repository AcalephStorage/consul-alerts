APP_NAME = consul-alerts
VERSION = 0.1.0

all: clean package

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

test: prepare
	@echo "--> Testing application"
	@go test -outputdir build/test ./...

build: test
	@echo "--> Building application"
	@echo "... linux-386"
	@GOOS=linux GOARCH=386 go build -o build/bin/linux-386/${VERSION}/${APP_NAME} -v .
	@echo "... linux-amd64"
	@GOOS=linux GOARCH=amd64 go build -o build/bin/linux-amd64/${VERSION}/${APP_NAME} -v .
	@echo "... darwin-amd64"
	@GOOS=darwin GOARCH=amd64 go build -o build/bin/darwin-amd64/${VERSION}/${APP_NAME} -v .

package: build
	@echo "--> Packaging application"
	@tar cf build/tar/${APP_NAME}-${VERSION}-linux-386.tar -C build/bin/linux-386/${VERSION} ${APP_NAME}
	@tar cf build/tar/${APP_NAME}-${VERSION}-linux-amd64.tar -C build/bin/linux-amd64/${VERSION} ${APP_NAME}
	@tar cf build/tar/${APP_NAME}-${VERSION}-darwin-amd64.tar -C build/bin/darwin-amd64/${VERSION} ${APP_NAME}
