#!/usr/bin/env bash
set -euo pipefail
# shellcheck source=../_shared/scripts/deploy-lib.sh
source "$(dirname "$0")/../_shared/scripts/deploy-lib.sh"
deploy_init "$0"
deploy_require_cli

tmpdir=$(mktemp -d)
trap 'rm -rf "$tmpdir"' EXIT

export ACTION_FQN=acme.corp.sdkdemo.StampDemoOnEnter

deploy_render_mdl mdl/02-lifecycle.mdl "$tmpdir/02-lifecycle.mdl" "ACTION_FQN=$ACTION_FQN"
deploy_apply_mdl "$tmpdir/02-lifecycle.mdl"
deploy_apply_mdl mdl/03-bind-object-lifecycle.mdl

deploy_render_mdl mdl/01-object.mdl "$tmpdir/01-object.mdl"
vpk="$tmpdir/action.vpk"
bash "$SHARED_SCRIPTS/package-vpk.sh" . "$vpk" \
  --component "10:Object:sdk_demo__c:$tmpdir/01-object.mdl"
deploy_vpk "$vpk"

echo "Done. Example FQNs: acme.corp.sdkdemo.SetTitleShared, ClearTitle, StampDemoOnEnter, StampDemoName"
