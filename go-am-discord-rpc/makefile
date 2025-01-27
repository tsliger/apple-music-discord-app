BINARY_NAME=go-am-discord-rpc

BINARY_NAME=go-am-discord-rpc
OUTPUT_DIR=./builds

# build out go applications in both linux-x86-x64 and Apple Silicon
# naming conventions match to tauri specs
BINARY_DARWIN=${OUTPUT_DIR}/${BINARY_NAME}-aarch64-apple-darwin
BINARY_LINUX=${OUTPUT_DIR}/${BINARY_NAME}-x86_64-unknown-linux-gnu

all: build

build: ${BINARY_DARWIN} ${BINARY_LINUX}
	cp ${OUTPUT_DIR}/* ../src-tauri/binaries
	rm ${BINARY_DARWIN}
	rm ${BINARY_LINUX}

${BINARY_DARWIN}:
	GOOS=darwin GOARCH=arm64 go build -o $@ .

${BINARY_LINUX}:
	GOOS=linux GOARCH=amd64 go build -o $@ .
