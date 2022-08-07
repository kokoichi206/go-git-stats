BINARY_NAME=ggs
BUILD_LDFLAGS = "-s -w -X main.revision=$(CURRENT_REVISION)"

.PHONY: build
build:
	go build -o ${BINARY_NAME} -ldflags=$(BUILD_LDFLAGS)
