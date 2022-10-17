PACKAGE = github.com/mrmoneyc/slack-exporter
BIN  = slack-exporter
DATE    ?= $(shell date +%FT%T%z)
VERSION ?= $(shell git describe --tags --always --dirty --match=v* 2> /dev/null || \
			cat $(CURDIR)/.version 2> /dev/null || echo nightly)
BIN_DIR = $(CURDIR)/bin
DIST_DIR ?= $(CURDIR)/dist
SCRIPT_DIR = $(CURDIR)/scripts
WINDOWS_RUN = run.bat

GO = go

# Colors
RED	:= $(shell tput -Txterm setaf 1)
GREEN	:= $(shell tput -Txterm setaf 2)
YELLOW	:= $(shell tput -Txterm setaf 3)
BLUE	:= $(shell tput -Txterm setaf 4)
WHITE	:= $(shell tput -Txterm setaf 7)
RESET	:= $(shell tput -Txterm sgr0)

TARGET_MAX_CHAR_NUM := 20

V ?= 0
Q := $(if $(filter 1,$V),,@)
M := $(shell printf "$(YELLOW)>>$(RESET)")

all: help

## Build binary for current OS and ARCH
build: $(BIN_DIR) ; $(info $(M) Building $(BIN)...)
	$Q $(GO) build \
		-tags release \
		-ldflags '-X $(PACKAGE)/pkg/cmd.Version=$(VERSION) -X $(PACKAGE)/pkg/cmd.BuildDate=$(DATE)' \
		-o $(BIN_DIR)/$(BIN) cmd/slack-exporter/main.go

## Build cross-compile Linux binary
build-linux: build-linux-amd64 build-linux-arm64

## Build cross-compile macOS binary
build-darwin: build-darwin-amd64 build-darwin-arm64

## Build cross-compile Windows binary
build-windows: build-windows-amd64

## Build cross-compile binary
build-cross: build-linux build-darwin build-windows

dist-linux: \
	$(DIST_DIR)/$(VERSION)/$(BIN)-$(VERSION)-linux-arm64.tar.gz \
	$(DIST_DIR)/$(VERSION)/$(BIN)-$(VERSION)-linux-amd64.tar.gz

dist-darwin: \
	$(DIST_DIR)/$(VERSION)/$(BIN)-$(VERSION)-darwin-arm64.tar.gz \
	$(DIST_DIR)/$(VERSION)/$(BIN)-$(VERSION)-darwin-amd64.tar.gz

dist-windows: \
	$(DIST_DIR)/$(VERSION)/$(BIN)-$(VERSION)-windows-amd64.zip

## Generate dist archive file
dist: dist-linux dist-darwin dist-windows

# Tools

$(BIN_DIR):
	@mkdir -p $@

$(DIST_DIR)/$(VERSION):
	@mkdir -p $(DIST_DIR)/$(VERSION)

$(DIST_DIR)/$(VERSION)/$(BIN)-$(VERSION)-%.tar.gz: build-% $(DIST_DIR)/$(VERSION)
		tar -C $(BIN_DIR) -cvzf $@ $(BIN)-$*

$(DIST_DIR)/$(VERSION)/$(BIN)-$(VERSION)-windows-%.zip: build-windows-% $(DIST_DIR)/$(VERSION)
		cp $(SCRIPT_DIR)/$(WINDOWS_RUN) $(BIN_DIR)/$(WINDOWS_RUN)
		cd $(BIN_DIR) && zip -r $@ $(BIN).exe $(WINDOWS_RUN)

build-%: $(BIN_DIR) ; $(info $(M) Building $(BIN)-$*...)
	$Q os=$$(echo $* | cut -d "-" -f 1); \
	arch=$$(echo $* | cut -d "-" -f 2); \
	if [ "$$os" != "$$arch" ]; then \
		output=$(BIN_DIR)/$(BIN)-$*; \
		[[ "$$os" == "windows" ]] && output=$(BIN_DIR)/$(BIN).exe; \
		GOOS=$$os GOARCH=$$arch $(GO) build \
			-tags release \
			-ldflags '-X $(PACKAGE)/pkg/cmd.Version=$(VERSION) -X $(PACKAGE)/pkg/cmd.BuildDate=$(DATE)' \
			-o $$output cmd/slack-exporter/main.go; \
	fi

# Misc

## Cleanup
clean: ; $(info $(M) Cleaning...)
	rm -rf $(BIN_DIR) $(DIST_DIR)
	rm -rf test/tests.* test/coverage.*

## Show version
version:
	@echo $(VERSION)

## Show help
help:
	@echo 'Usage:'
	@echo '  $(YELLOW)make$(RESET) $(GREEN)<target>$(RESET)'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\_0-9]+:/ { \
		helpMessage = match(lastLine, /^## (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 3, RLENGTH); \
			printf "  $(YELLOW)%-$(TARGET_MAX_CHAR_NUM)s$(RESET) $(GREEN)%s$(RESET)\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

require-%:
	@ if [ "$($(*))" = "" ]; then \
		echo "Environment variable $* not set"; \
		exit 1; \
	fi

.PHONY: all build build-cross build-linux build-darwin build-windows dist clean version help .FORCE
