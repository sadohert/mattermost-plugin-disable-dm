# Makefile for Mattermost Disable DM Plugin
# Based on mattermost-plugin-starter-template

GO ?= go
PLUGIN_ID = disable-dm
PLUGIN_VERSION = 2.0.0
BUNDLE_NAME = $(PLUGIN_ID)-$(PLUGIN_VERSION).tar.gz
GO_BUILD_FLAGS ?=
MM_DEBUG ?=

ifneq ($(MM_DEBUG),)
	GO_BUILD_GCFLAGS = -gcflags "all=-N -l"
else
	GO_BUILD_GCFLAGS =
endif

## Default target
.PHONY: all
all: check-style dist

## Download Go dependencies
.PHONY: deps
deps:
	$(GO) mod download
	$(GO) mod tidy

## Runs govet and gofmt
.PHONY: check-style
check-style: govet gofmt

.PHONY: govet
govet:
	@echo Running govet
	$(GO) vet ./...
	@echo Govet success

.PHONY: gofmt
gofmt:
	@echo Running gofmt
	@if [ -n "$$(gofmt -l .)" ]; then \
		echo "gofmt failed:"; \
		gofmt -d .; \
		exit 1; \
	fi
	@echo Gofmt success

## Runs any tests
.PHONY: test
test:
	$(GO) test -v -race ./...

## Builds the server executables
.PHONY: server
server:
ifneq ($(MM_DEBUG),)
	$(info DEBUG mode is on; to disable, unset MM_DEBUG)
endif
	mkdir -p server/dist
	cd server && env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) $(GO_BUILD_GCFLAGS) -trimpath -o dist/plugin-linux-amd64
	cd server && env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build $(GO_BUILD_FLAGS) $(GO_BUILD_GCFLAGS) -trimpath -o dist/plugin-linux-arm64
	cd server && env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) $(GO_BUILD_GCFLAGS) -trimpath -o dist/plugin-darwin-amd64
	cd server && env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build $(GO_BUILD_FLAGS) $(GO_BUILD_GCFLAGS) -trimpath -o dist/plugin-darwin-arm64
	cd server && env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) $(GO_BUILD_GCFLAGS) -trimpath -o dist/plugin-windows-amd64.exe

## Bundles the plugin
.PHONY: bundle
bundle:
	rm -rf dist/
	mkdir -p dist/$(PLUGIN_ID)
	cp plugin.json dist/$(PLUGIN_ID)/
	mkdir -p dist/$(PLUGIN_ID)/server
	cp -r server/dist dist/$(PLUGIN_ID)/server/
	cd dist && tar -czf $(BUNDLE_NAME) $(PLUGIN_ID)
	@echo plugin built at: dist/$(BUNDLE_NAME)

## Build the plugin bundle
.PHONY: dist
dist: deps check-style server bundle

## Cleans build artifacts
.PHONY: clean
clean:
	rm -rf dist
	rm -rf server/dist
	$(GO) clean -cache

## Deploy plugin to local Mattermost server (requires MM_SERVICESETTINGS_SITEURL and MM_ADMIN_TOKEN)
.PHONY: deploy
deploy: dist
	@if [ -z "$(MM_SERVICESETTINGS_SITEURL)" ]; then \
		echo "Error: MM_SERVICESETTINGS_SITEURL is not set"; \
		exit 1; \
	fi
	@if [ -z "$(MM_ADMIN_TOKEN)" ]; then \
		echo "Error: MM_ADMIN_TOKEN is not set"; \
		exit 1; \
	fi
	curl -F "plugin=@dist/$(BUNDLE_NAME)" -H "Authorization: Bearer $(MM_ADMIN_TOKEN)" $(MM_SERVICESETTINGS_SITEURL)/api/v4/plugins

## Show help
.PHONY: help
help:
	@echo "Mattermost Disable DM Plugin - available make targets:"
	@echo ""
	@echo "  all          - Run check-style and build distribution (default)"
	@echo "  dist         - Build plugin bundle for all platforms"
	@echo "  server       - Build server binaries only"
	@echo "  bundle       - Create plugin bundle from built artifacts"
	@echo "  deps         - Download and tidy Go dependencies"
	@echo "  check-style  - Run govet and gofmt"
	@echo "  test         - Run tests"
	@echo "  deploy       - Deploy to local Mattermost (requires env vars)"
	@echo "  clean        - Remove build artifacts"
	@echo "  help         - Show this help message"
