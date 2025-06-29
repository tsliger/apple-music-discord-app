BINARY_NAME=go-am-discord-rpc
OUTPUT_DIR=./builds

BINARY_WINDOWS_AMD64=$(OUTPUT_DIR)/$(BINARY_NAME)-x86_64-pc-windows-msvc.exe
BINARY_DARWIN_ARM64=$(OUTPUT_DIR)/$(BINARY_NAME)-aarch64-apple-darwin
BINARY_DARWIN_AMD64=$(OUTPUT_DIR)/$(BINARY_NAME)-x86_64-apple-darwin

HOST_OS := $(shell go env GOOS)

all: build

build: prebuild
ifeq ($(HOST_OS),windows)
	@echo "Windows detected — building Windows binary first"
	@$(MAKE) windows
	@echo "Copying binary..."
	@cp $(BINARY_WINDOWS_AMD64) ../src-tauri/binaries/
	@rm -f $(BINARY_WINDOWS_AMD64)
else ifeq ($(HOST_OS),darwin)
	@echo "macOS detected — building all targets"
	@$(MAKE) darwin
	@cp $(BINARY_DARWIN_ARM64) ../src-tauri/binaries/
	@cp $(BINARY_DARWIN_AMD64) ../src-tauri/binaries/
	@rm -f $(BINARY_DARWIN_ARM64) $(BINARY_DARWIN_AMD64)
else
	@echo "Unsupported OS for Darwin cross-compilation — building Windows binary only"
	@$(MAKE) windows
endif

prebuild:
	@mkdir -p ../src-tauri/binaries
	@mkdir -p $(OUTPUT_DIR)

windows: $(BINARY_WINDOWS_AMD64)

darwin: $(BINARY_DARWIN_ARM64) $(BINARY_DARWIN_AMD64)

$(BINARY_WINDOWS_AMD64):
	GOOS=windows GOARCH=amd64 go build -o $@ .

$(BINARY_DARWIN_ARM64):
	GOOS=darwin GOARCH=arm64 go build -o $@ .

$(BINARY_DARWIN_AMD64):
	GOOS=darwin GOARCH=amd64 go build -o $@ .

clean:
	@echo "Cleaning build output..."
	@rm -rf $(OUTPUT_DIR) ../src-tauri/binaries
