BINARY_NAME=go-am-discord-rpc
OUTPUT_DIR=./builds

BINARY_WINDOWS_AMD64=$(OUTPUT_DIR)/$(BINARY_NAME)-x86_64-pc-windows-msvc.exe
BINARY_DARWIN_ARM64=$(OUTPUT_DIR)/$(BINARY_NAME)-aarch64-apple-darwin
BINARY_DARWIN_AMD64=$(OUTPUT_DIR)/$(BINARY_NAME)-x86_64-apple-darwin

# Detect host OS via Go
HOST_OS := $(shell go env GOOS)

all: build

build:
	@mkdir -p ../src-tauri/binaries

ifeq ($(HOST_OS),windows)
	@echo "Windows detected — building Windows binary first"
	@$(MAKE) windows
	@echo "Copying directory and exe on Windows only..."
	@cmd /c copy /Y "$(subst /,\,$(BINARY_WINDOWS_AMD64))" "..\src-tauri\binaries\"
	@cmd /c del /Q "$(subst /,\,$(BINARY_WINDOWS_AMD64))"
else ifeq ($(HOST_OS),darwin)
	@echo "macOS detected — building all targets"
	@$(MAKE) darwin
	cp "$(BINARY_DARWIN_ARM64)" "../src-tauri/binaries/"
	cp "$(BINARY_DARWIN_AMD64)" "../src-tauri/binaries/"
	rm "$(BINARY_DARWIN_ARM64)"
	rm "$(BINARY_DARWIN_AMD64)"
else
	@echo "Unsupported OS for Darwin cross-compilation — building Windows binary only"
	@$(MAKE) windows-gnu
endif

windows: $(BINARY_WINDOWS_AMD64)

darwin: $(BINARY_DARWIN_ARM64) $(BINARY_DARWIN_AMD64)


$(BINARY_WINDOWS_AMD64):
	@if not exist $(subst /,\,$(OUTPUT_DIR)) mkdir $(subst /,\,$(OUTPUT_DIR))
	go build -o $(subst /,\,$@) .

$(BINARY_DARWIN_ARM64):
	@mkdir -p $(OUTPUT_DIR)
	GOOS=darwin GOARCH=arm64 go build -o $@ .

$(BINARY_DARWIN_AMD64):
	@mkdir -p $(OUTPUT_DIR)
	GOOS=darwin GOARCH=amd64 go build -o $@ .

clean:
ifeq ($(HOST_OS),windows)
	@del /q $(subst /,\,$(OUTPUT_DIR))\*
else
	@rm -f $(OUTPUT_DIR)/*
endif
