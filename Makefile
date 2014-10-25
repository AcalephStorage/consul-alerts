APP_NAME = consul-alerts

all: clean deps install-global

clean:
	@echo "--> Cleaning build"
	@rm -rf ./build

prepare:
	@mkdir -p build/bin/${PLATFORM}
	@mkdir -p build/test
	@mkdir -p build/doc
	@mkdir -p build/tar

deps:
	@echo "--> Getting Dependencies"
	@go get github.com/mattn/gom
	@go install github.com/mattn/gom
	@gom install

format:
	@echo "--> Formatting source code"
	@gom exec go fmt ./...

test: prepare
	@echo "--> Testing application"
	@gom exec go test -outputdir build/test ./...

build: test
	@echo "--> Building application"
	@gom exec go build -o build/bin/${APP_NAME} -v .

install-global: build
	@echo "--> Installing app"
	@gom exec cp build/bin/${APP_NAME} ${GOPATH}/bin/