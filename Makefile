.DEFAULT_GOAL := build

.PHONY: build lint build.docker run

build:
	tinygo build -o ./my-filter.wasm -scheduler=none -target=wasi -wasm-abi=generic ./main.go

build.docker:
	docker run --rm -it -w /tmp/proxy-wasm-go -v $(shell pwd):/tmp/proxy-wasm-go tinygo/tinygo-dev:latest \
		tinygo build -o /tmp/proxy-wasm-go/my-filter.wasm -scheduler=none -target=wasi \
		-wasm-abi=generic /tmp/proxy-wasm-go/main.go


lint:
	golangci-lint run --build-tags proxytest


run:
	docker-compose up