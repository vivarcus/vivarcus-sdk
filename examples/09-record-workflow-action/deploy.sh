#!/usr/bin/env bash
set -euo pipefail
# shellcheck source=../_shared/scripts/deploy-lib.sh
source "$(dirname "$0")/../_shared/scripts/deploy-lib.sh"
deploy_init "$0"
deploy_require_cli
deploy_apply_mdl mdl/01-object.mdl
deploy_apply_mdl mdl/02-lifecycle.mdl
deploy_apply_mdl mdl/03-bind-object-lifecycle.mdl
deploy_sdk_put workflowactions/capture_participants.go
deploy_apply_mdl mdl/04-workflow.mdl
echo "Done. FQN: acme.corp.sdkdemo.CaptureParticipants"
