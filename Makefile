main_package_path = .
binary_name = gator

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## test: run tests
.PHONY: test
test:
	go test -v ./...

## audit: run quality checks
.PHONY: audit
audit: test
	go mod tidy -diff
	go vet ./...
	test -z "$$(gofmt -l .)"

## build: build the binary
.PHONY: build
build:
	mkdir -p /tmp/bin
	go build -o /tmp/bin/${binary_name} ${main_package_path}

## run: build and run the app
.PHONY: run
run: build
	/tmp/bin/${binary_name}

## tidy: format and tidy
.PHONY: tidy
tidy:
	go mod tidy
	go fmt ./...
