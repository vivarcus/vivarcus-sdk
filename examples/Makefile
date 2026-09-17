# SDK examples — build & deploy demos.
# From repo root: make -C sdk/examples help
# Or: make sdk-update-field-demo  (delegates here)

SHELL := /bin/bash
ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
SHARED := $(ROOT)/_shared/scripts
DEMO_SH := $(SHARED)/dev-demo.sh
VIVARCUS_SDK ?= vivarcus-sdk
GOTOOLCHAIN ?= go1.26.2

DEMOS := 02-update-field 03-confirm-dialog 04-lifecycle-entry 99-full-stack
BUILD_EXAMPLES := 01-hello-action 02-update-field 03-confirm-dialog 04-lifecycle-entry 99-full-stack

.PHONY: help build build-all demo demo-all setup verify
.PHONY: build-01-hello-action build-02-update-field build-03-confirm-dialog build-04-lifecycle-entry build-99-full-stack
.PHONY: demo-02-update-field demo-03-confirm-dialog demo-04-lifecycle-entry demo-99-full-stack
.PHONY: setup-02-update-field setup-03-confirm-dialog setup-04-lifecycle-entry setup-99-full-stack
.PHONY: verify-02-update-field verify-03-confirm-dialog verify-04-lifecycle-entry verify-99-full-stack

help:
	@echo "SDK examples (cd $(ROOT))"
	@echo ""
	@echo "  make build-all                vivarcus-sdk build for 01–04 + 99"
	@echo "  make build-01-hello-action    build one example (--skip-compile)"
	@echo "  make demo-02-update-field     setup + verify (needs vivarcus auth)"
	@echo "  make setup-04-lifecycle-entry deploy only"
	@echo "  make verify-99-full-stack     verify only"
	@echo "  make demo-all                 run all $(words $(DEMOS)) demos"
	@echo ""
	@echo "Env: VIVARCUS_SDK_BIN VIVARCUS_BIN GOTOOLCHAIN"

build-all: build-01-hello-action build-02-update-field build-03-confirm-dialog build-04-lifecycle-entry build-99-full-stack

build-01-hello-action:
	GOTOOLCHAIN=$(GOTOOLCHAIN) $(VIVARCUS_SDK) build $(ROOT)/01-hello-action --skip-compile

build-02-update-field:
	GOTOOLCHAIN=$(GOTOOLCHAIN) $(VIVARCUS_SDK) build $(ROOT)/02-update-field -o $(ROOT)/02-update-field/action.wasm

build-03-confirm-dialog:
	GOTOOLCHAIN=$(GOTOOLCHAIN) $(VIVARCUS_SDK) build $(ROOT)/03-confirm-dialog -o $(ROOT)/03-confirm-dialog/action.wasm

build-04-lifecycle-entry:
	GOTOOLCHAIN=$(GOTOOLCHAIN) $(VIVARCUS_SDK) build $(ROOT)/04-lifecycle-entry -o $(ROOT)/04-lifecycle-entry/action.wasm

build-99-full-stack:
	GOTOOLCHAIN=$(GOTOOLCHAIN) $(VIVARCUS_SDK) build $(ROOT)/99-full-stack -o $(ROOT)/99-full-stack/action.wasm

demo-all: demo-02-update-field demo-03-confirm-dialog demo-04-lifecycle-entry demo-99-full-stack

demo-02-update-field:
	@bash $(DEMO_SH) $(ROOT)/02-update-field

demo-03-confirm-dialog:
	@bash $(DEMO_SH) $(ROOT)/03-confirm-dialog

demo-04-lifecycle-entry:
	@bash $(DEMO_SH) $(ROOT)/04-lifecycle-entry

demo-99-full-stack:
	@bash $(DEMO_SH) $(ROOT)/99-full-stack

setup-02-update-field:
	@bash $(DEMO_SH) $(ROOT)/02-update-field --setup-only

setup-03-confirm-dialog:
	@bash $(DEMO_SH) $(ROOT)/03-confirm-dialog --setup-only

setup-04-lifecycle-entry:
	@bash $(DEMO_SH) $(ROOT)/04-lifecycle-entry --setup-only

setup-99-full-stack:
	@bash $(DEMO_SH) $(ROOT)/99-full-stack --setup-only

verify-02-update-field:
	@bash $(DEMO_SH) $(ROOT)/02-update-field --verify-only

verify-03-confirm-dialog:
	@bash $(DEMO_SH) $(ROOT)/03-confirm-dialog --verify-only

verify-04-lifecycle-entry:
	@bash $(DEMO_SH) $(ROOT)/04-lifecycle-entry --verify-only

verify-99-full-stack:
	@bash $(DEMO_SH) $(ROOT)/99-full-stack --verify-only
