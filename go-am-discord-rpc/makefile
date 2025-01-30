BINARY_NAME=go-am-discord-rpc
OUTPUT_DIR=./builds

# build out go applications in both darwin-x86-x64 and Apple Silicon
# naming conventions match to tauri specs
BINARY_DARWIN_ARM64=${OUTPUT_DIR}/${BINARY_NAME}-aarch64-apple-darwin
BINARY_DARWIN_AMD64=${OUTPUT_DIR}/${BINARY_NAME}-x86_64-apple-darwin

all: build

build: ${BINARY_DARWIN_ARM64} ${BINARY_DARWIN_AMD64}
	cp ${OUTPUT_DIR}/* ../src-tauri/binaries
	rm ${BINARY_DARWIN_ARM64}
	rm ${BINARY_DARWIN_AMD64}

${BINARY_DARWIN_ARM64}:
	GOOS=darwin GOARCH=arm64 go build -o $@ .

${BINARY_DARWIN_AMD64}:
	GOOS=darwin GOARCH=amd64 go build -o $@ .
