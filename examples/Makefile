# SDK examples — build & deploy demos.
# From repo root: make -C sdk/examples help
# Or: make sdk-multi-component-demo  (delegates here)

SHELL := /bin/bash
ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
SHARED := $(ROOT)/_shared/scripts
DEMO_SH := $(SHARED)/dev-demo.sh
VIVARCUS_SDK ?= vivarcus-sdk
GOTOOLCHAIN ?= go1.26.2

DEMOS := multi-component 02-update-field 03-confirm-dialog 04-lifecycle-entry
BUILD_EXAMPLES := multi-component 01-hello-action 02-update-field 03-confirm-dialog 04-lifecycle-entry 05-stamp-trigger

.PHONY: help build build-all demo demo-all setup verify
.PHONY: build-multi-component build-01-hello-action build-02-update-field build-03-confirm-dialog build-04-lifecycle-entry build-05-stamp-trigger
.PHONY: demo-multi-component demo-02-update-field demo-03-confirm-dialog demo-04-lifecycle-entry
.PHONY: setup-multi-component setup-02-update-field setup-03-confirm-dialog setup-04-lifecycle-entry
.PHONY: verify-multi-component verify-02-update-field verify-03-confirm-dialog verify-04-lifecycle-entry

help:
	@echo "SDK examples (cd $(ROOT))"
	@echo ""
	@echo "  make demo-multi-component     多 Action/Trigger 端到端（推荐入口）"
	@echo "  make build-multi-component    vivarcus-sdk build multi-component"
	@echo "  make verify-multi-component   verify only"
	@echo "  make build-all                vivarcus-sdk build for all examples"
	@echo "  make demo-02-update-field     单按钮 demo（需 vivarcus auth）"
	@echo "  make demo-all                 run all $(words $(DEMOS)) demos"
	@echo ""
	@echo "Env: VIVARCUS_SDK_BIN VIVARCUS_BIN GOTOOLCHAIN"

build-all: build-multi-component build-01-hello-action build-02-update-field build-03-confirm-dialog build-04-lifecycle-entry build-05-stamp-trigger

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

demo-all: demo-multi-component demo-02-update-field demo-03-confirm-dialog demo-04-lifecycle-entry

demo-multi-component:
	@bash $(DEMO_SH) $(ROOT)/multi-component

demo-02-update-field:
	@bash $(DEMO_SH) $(ROOT)/02-update-field

demo-03-confirm-dialog:
	@bash $(DEMO_SH) $(ROOT)/03-confirm-dialog

demo-04-lifecycle-entry:
	@bash $(DEMO_SH) $(ROOT)/04-lifecycle-entry

setup-multi-component:
	@bash $(DEMO_SH) $(ROOT)/multi-component --setup-only

setup-02-update-field:
	@bash $(DEMO_SH) $(ROOT)/02-update-field --setup-only

setup-03-confirm-dialog:
	@bash $(DEMO_SH) $(ROOT)/03-confirm-dialog --setup-only

setup-04-lifecycle-entry:
	@bash $(DEMO_SH) $(ROOT)/04-lifecycle-entry --setup-only

verify-multi-component:
	@bash $(DEMO_SH) $(ROOT)/multi-component --verify-only

verify-02-update-field:
	@bash $(DEMO_SH) $(ROOT)/02-update-field --verify-only

verify-03-confirm-dialog:
	@bash $(DEMO_SH) $(ROOT)/03-confirm-dialog --verify-only

verify-04-lifecycle-entry:
	@bash $(DEMO_SH) $(ROOT)/04-lifecycle-entry --verify-only
