#!/usr/bin/env bash
# Unified entry for examples demos.
#
# Usage:
#   ./_shared/scripts/dev-demo.sh <example-dir> [--setup-only|--verify-only]
#   make -C examples demo-02-update-field
set -euo pipefail

SHARED_DIR=$(cd "$(dirname "$0")" && pwd)
EXAMPLES_ROOT=$(cd "$SHARED_DIR/../.." && pwd)

usage() {
  sed -n '2,8p' "$0" | cut -c3-
  echo "examples: multi-component 02-update-field 03-confirm-dialog 04-lifecycle-entry"
}

if [ $# -lt 1 ]; then
  usage >&2
  exit 2
fi

EXAMPLE_DIR=$(cd "$1" && pwd)
shift

SETUP=1
VERIFY=1
while [ $# -gt 0 ]; do
  case "$1" in
    --setup-only) VERIFY=0 ;;
    --verify-only) SETUP=0 ;;
    -h|--help)
      usage
      exit 0
      ;;
    *) echo "unknown arg: $1" >&2; exit 2 ;;
  esac
  shift
done

case "$EXAMPLE_DIR" in
  "$EXAMPLES_ROOT"/*) ;;
  *)
    echo "example dir must live under $EXAMPLES_ROOT" >&2
    exit 2
    ;;
esac

ENV_FILE="$EXAMPLE_DIR/demo.env"
if [ -f "$ENV_FILE" ]; then
  # shellcheck source=/dev/null
  source "$ENV_FILE"
fi

# shellcheck source=/dev/null
source "$SHARED_DIR/dev-demo-common.sh"

VIVARCUS="${VIVARCUS_BIN:-vivarcus}"
VIVARCUS_SDK="${VIVARCUS_SDK_BIN:-vivarcus-sdk}"
export GOTOOLCHAIN="${GOTOOLCHAIN:-go1.26.2}"
export VIVARCUS_BIN="$VIVARCUS"

MDL_DIR="$EXAMPLE_DIR/${DEMO_MDL_DIR:-mdl}"
DEMO_DIST=$(demo_dist_init)
WORKDIR="${DEMO_WORKDIR:-$DEMO_DIST}"
mkdir -p "$WORKDIR"

DEMO_SCENARIO="${DEMO_SCENARIO:-}"
if [ -z "$DEMO_SCENARIO" ]; then
  case "$(basename "$EXAMPLE_DIR")" in
    02-update-field|03-confirm-dialog) DEMO_SCENARIO=user-action ;;
    04-lifecycle-entry) DEMO_SCENARIO=lifecycle ;;
    multi-component) DEMO_SCENARIO=multi-component ;;
    *)
      echo "unknown example $(basename "$EXAMPLE_DIR")" >&2
      exit 1
      ;;
  esac
fi
export DEMO_SCENARIO
demo_apply_example_defaults


VERIFY_SCRIPT="$SHARED_DIR/verify-dev-demo.py"

setup_demo() {
  case "$DEMO_SCENARIO" in
    user-action) demo_scenario_user_action ;;
    lifecycle) demo_scenario_lifecycle ;;
    multi-component) demo_scenario_multi_component ;;
    *)
      echo "unknown DEMO_SCENARIO=$DEMO_SCENARIO" >&2
      exit 1
      ;;
  esac
}

verify_demo() {
  demo_require_auth
  case "$DEMO_SCENARIO" in
    user-action)
      demo_fill_from_describe
      demo_export_user_action_env
      ;;
    lifecycle)
      demo_fill_from_describe
      demo_export_lifecycle_verify_env
      ;;
    multi-component)
      demo_fill_multi_component_from_describe
      demo_export_multi_component_verify_env
      ;;
  esac
  python3 "$VERIFY_SCRIPT"
}

if [ "$SETUP" = 1 ]; then setup_demo; fi
if [ "$VERIFY" = 1 ]; then verify_demo; fi
