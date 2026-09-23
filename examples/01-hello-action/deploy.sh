#!/usr/bin/env bash
# Deploy this example to the logged-in Vault (see README.md).
set -euo pipefail
# shellcheck source=../_shared/scripts/deploy-lib.sh
source "$(dirname "$0")/../_shared/scripts/deploy-lib.sh"
deploy_init "$0"
deploy_require_cli
deploy_sdk_put actions/noop_action.go
echo "Done. FQN: acme.corp.sdkdemo.NoopAction"
