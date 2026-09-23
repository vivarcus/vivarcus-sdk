#!/usr/bin/env bash
set -euo pipefail
# shellcheck source=../_shared/scripts/deploy-lib.sh
source "$(dirname "$0")/../_shared/scripts/deploy-lib.sh"
deploy_init "$0"
deploy_require_cli
deploy_apply_mdl mdl/01-webapigroup.mdl
deploy_sdk_put webapis/hello_api.go
echo "Done. FQN: acme.corp.sdkdemo.HelloAPI"
