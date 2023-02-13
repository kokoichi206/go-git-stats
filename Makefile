BINARY_NAME=ggs
BUILD_LDFLAGS = "-s -w -X main.revision=$(CURRENT_REVISION)"

.PHONY: build help
.DEFAULT_GOAL := help

build:	## バージョン値にリビジョンを埋め込んでビルド。
	go build -o ${BINARY_NAME} -ldflags=$(BUILD_LDFLAGS) cmd/ggs/*

help:	## https://postd.cc/auto-documented-makefile/
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
