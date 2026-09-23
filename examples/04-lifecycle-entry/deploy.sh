#!/usr/bin/env bash
set -euo pipefail
# shellcheck source=../_shared/scripts/deploy-lib.sh
source "$(dirname "$0")/../_shared/scripts/deploy-lib.sh"
deploy_init "$0"
deploy_require_cli

tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT

export ACTION_FQN=acme.corp.sdkdemo.StampOnEnter

deploy_apply_mdl mdl/01-object.mdl
deploy_sdk_put actions/stamp_on_enter.go

deploy_render_mdl mdl/02-lifecycle.mdl "$tmpdir/02-lifecycle.mdl" "ACTION_FQN=$ACTION_FQN"
deploy_render_mdl mdl/04-workflow-action.mdl "$tmpdir/04-workflow-action.mdl" "ACTION_FQN=$ACTION_FQN"
deploy_render_mdl mdl/05-workflow-cancel.mdl "$tmpdir/05-workflow-cancel.mdl" "ACTION_FQN=$ACTION_FQN"

deploy_apply_mdl "$tmpdir/02-lifecycle.mdl"
deploy_apply_mdl mdl/03-bind-object-lifecycle.mdl
deploy_apply_mdl "$tmpdir/04-workflow-action.mdl"
deploy_apply_mdl "$tmpdir/05-workflow-cancel.mdl"

echo "Done. FQN: $ACTION_FQN"
