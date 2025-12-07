.PHONY: ckeck install upload

BUILD_VERSION ?= $(shell cat VERSION)
BUILD_OUTPUT ?= fractals
GO ?= GO111MODULE=on CGO_ENABLED=0 go
GOOS ?= $(shell $(GO) version | cut -d' ' -f4 | cut -d'/' -f1)
GOARCH ?= $(shell $(GO) version | cut -d' ' -f4 | cut -d'/' -f2)
DOCKER_IMAGE ?= ghcr.io/teadove/fractals:$(BUILD_VERSION)

test:
	@$(GO) test -v ./...
	pre-commit run

run:
	@$(GO) run main.go

run-image:
	@$(GO) run main.go image --xpos -0.00575 --radius 0.00125689

build:
	@$(GO) build -o $(BUILD_OUTPUT) main.go

clean:
	@echo -n ">> CLEAN"
	@$(GO) clean -i ./...
	@rm -f fractals-*-*
	@rm -rf dist/*
	@printf '%s\n' '$(OK)'


crosscompile:
	@echo -n ">> CROSSCOMPILE linux/amd64"
	@GOOS=linux GOARCH=amd64 $(GO) build -o $(BUILD_OUTPUT)-$(BUILD_VERSION)-linux-amd64
	@printf '%s\n' '$(OK)'
	@echo -n ">> CROSSCOMPILE darwin/amd64"
	@GOOS=darwin GOARCH=amd64 $(GO) build -o $(BUILD_OUTPUT)-$(BUILD_VERSION)-darwin-amd64
	@printf '%s\n' '$(OK)'
	@echo -n ">> CROSSCOMPILE windows/amd64"
	@GOOS=windows GOARCH=amd64 $(GO) build -o $(BUILD_OUTPUT)-$(BUILD_VERSION)-windows-amd64
	@printf '%s\n' '$(OK)'

	@echo -n ">> CROSSCOMPILE linux/arm64"
	@GOOS=linux GOARCH=arm64 $(GO) build -o $(BUILD_OUTPUT)-$(BUILD_VERSION)-linux-arm64
	@printf '%s\n' '$(OK)'
	@echo -n ">> CROSSCOMPILE darwin/arm64"
	@GOOS=darwin GOARCH=arm64 $(GO) build -o $(BUILD_OUTPUT)-$(BUILD_VERSION)-darwin-arm64
	@printf '%s\n' '$(OK)'
	@echo -n ">> CROSSCOMPILE windows/arm64"
	@GOOS=windows GOARCH=arm64 $(GO) build -o $(BUILD_OUTPUT)-$(BUILD_VERSION)-windows-arm64
	@printf '%s\n' '$(OK)'
