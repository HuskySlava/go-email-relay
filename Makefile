APP_NAME := go-email-relay
BUILD_DIR := bin

.PHONY: run build build-linux build-macos clean

run:
	go run .

build: build-linux build-macos

build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(APP_NAME)-linux-amd64 .

build-macos:
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(APP_NAME)-darwin-arm64 .

clean:
	rm -rf $(BUILD_DIR)
