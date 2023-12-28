ROOT_DIR := ./
SOURCE_FILES := $(wildcard *.go)

BINARY_NAME := ServiceId

.PHONY: build
build:
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME} $(ROOT_DIR)$(SOURCE_FILES)

