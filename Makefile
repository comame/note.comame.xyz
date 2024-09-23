.PHONY: build
build: static_files server

.PHONY: run
run: build
	export $$(cat ./.env) && ./out/server

.PHONY: server
server: static_files
	go build -o out/server

# Assets

.PHONY: static_files
static_files: wasm

.PHONY: wasm
wasm: clean
	cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" out/dist/wasm_exec.js
	GOOS=js GOARCH=wasm go build -o out/dist/goapp.wasm .

.PHONY: clean
clean:
	rm -rf out/*
	rm -f internal/server/out
	mkdir -p out/dist

# Tests

.PHONY: test_all
test_all: test_wasm
	go test ./...

.PHONY: test_wasm
test_wasm:
	GOOS=js GOARCH=wasm go test -exec="$$(go env GOPATH)/bin/wasmbrowsertest" ./...

.PHONY: test_coverage
test_coverage:
	go test ./... -cover -coverprofile=cover.out
	go tool cover -html cover.out -o cover.html
