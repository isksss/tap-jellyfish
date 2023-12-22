DIST_DIR=dist
WASMNAME=app.wasm
BINNAME=jellyfish

.PHONY: wasm
wasm:
	@GOOS=js GOARCH=wasm go build -o $(DIST_DIR)/$(WASMNAME) .

.PHONY: build
build:
	@go build -o $(DIST_DIR)/$(BINNAME) .

.PHONY: run
run: build
	@./$(DIST_DIR)/$(BINNAME)

.PHONY: debug
debug: build
	@DEBUG_JELLYFISH=true ./$(DIST_DIR)/$(BINNAME)