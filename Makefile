# Vivarcus SDK examples — build helpers (optional; see examples/README.md).

SHELL := /bin/bash
ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/examples
VIVARCUS_SDK ?= vivarcus-sdk
GOTOOLCHAIN ?= go1.22.12

.PHONY: help build-all build-multi-component build-01-hello-action build-02-update-field build-03-confirm-dialog build-04-lifecycle-entry build-05-stamp-trigger build-06-query-field

help:
	@echo "SDK examples (cd examples/ or use targets below)"
	@echo "  make build-multi-component"
	@echo "  make build-all"

build-all: build-multi-component build-01-hello-action build-02-update-field build-03-confirm-dialog build-04-lifecycle-entry build-05-stamp-trigger build-06-query-field

build-multi-component:
	GOTOOLCHAIN=$(GOTOOLCHAIN) $(VIVARCUS_SDK) build $(ROOT)/multi-component -o $(ROOT)/multi-component/action.wasm

build-01-hello-action:
	GOTOOLCHAIN=$(GOTOOLCHAIN) $(VIVARCUS_SDK) build $(ROOT)/01-hello-action --skip-compile

build-02-update-field:
	GOTOOLCHAIN=$(GOTOOLCHAIN) $(VIVARCUS_SDK) build $(ROOT)/02-update-field -o $(ROOT)/02-update-field/action.wasm

build-03-confirm-dialog:
	GOTOOLCHAIN=$(GOTOOLCHAIN) $(VIVARCUS_SDK) build $(ROOT)/03-confirm-dialog -o $(ROOT)/03-confirm-dialog/action.wasm

build-04-lifecycle-entry:
	GOTOOLCHAIN=$(GOTOOLCHAIN) $(VIVARCUS_SDK) build $(ROOT)/04-lifecycle-entry -o $(ROOT)/04-lifecycle-entry/action.wasm

build-05-stamp-trigger:
	GOTOOLCHAIN=$(GOTOOLCHAIN) $(VIVARCUS_SDK) build $(ROOT)/05-stamp-trigger -o $(ROOT)/05-stamp-trigger/action.wasm

build-06-query-field:
	GOTOOLCHAIN=$(GOTOOLCHAIN) $(VIVARCUS_SDK) build $(ROOT)/06-query-field -o $(ROOT)/06-query-field/action.wasm
