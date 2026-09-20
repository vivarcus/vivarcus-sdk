# shellcheck shell=bash
# Shared helpers for sdk/examples demos (_shared/scripts/dev-demo.sh).
# Source after setting EXAMPLE_DIR, WORKDIR, VIVARCUS, VIVARCUS_SDK.

SHARED_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
EXAMPLES_ROOT="$(cd "$SHARED_DIR/../.." && pwd)"
RENDER_MDL="$SHARED_DIR/render-mdl.py"
PACKAGE_VPK="$SHARED_DIR/package-vpk.sh"

demo_dist_dir() {
  echo "$EXAMPLE_DIR/${DEMO_DIST_DIR:-dist}"
}

# Create dist/ (mdl/, unpacked/) and export DEMO_DIST for artifact persistence.
demo_dist_init() {
  DEMO_DIST=$(demo_dist_dir)
  rm -rf "$DEMO_DIST"
  mkdir -p "$DEMO_DIST/mdl" "$DEMO_DIST/unpacked"
  export DEMO_DIST
  echo "$DEMO_DIST"
}

demo_dist_save_mdl() {
  local src="$1"
  local name="${2:-$(basename "$src")}"
  local dest="$DEMO_DIST/mdl/$name"
  if [ "$src" != "$dest" ]; then
    cp "$src" "$dest"
  fi
}

demo_dist_publish_vpk() {
  local vpk="$1"
  local name="${2:-$(basename "$vpk")}"
  local dest="$DEMO_DIST/$name"
  if [ "$vpk" != "$dest" ]; then
    cp -f "$vpk" "$dest"
    vpk="$dest"
  fi
  local unpack="$DEMO_DIST/unpacked/${name%.vpk}"
  rm -rf "$unpack"
  mkdir -p "$unpack"
  unzip -q -o "$vpk" -d "$unpack"
  echo "artifacts: $vpk (unpacked/${name%.vpk}/)" >&2
}

demo_require_auth() {
  if ! "$VIVARCUS" auth status --json >/dev/null 2>&1; then
    cat >&2 <<EOF
not authenticated — configure vivarcus CLI first, for example:
  vivarcus auth login
  vivarcus config set endpoint http://127.0.0.1:8080
  vivarcus config set default_vault <vault-uuid>
or export VIVARCUS_TOKEN and VIVARCUS_VAULT
EOF
    exit 1
  fi
}

demo_vault_args() {
  if [ -n "${VIVARCUS_VAULT:-}" ]; then
    echo --vault "$VIVARCUS_VAULT"
  fi
}

demo_apply_mdl_file() {
  local rendered="$1"
  local label="$2"
  if ! "$VIVARCUS" component apply-mdl --confirm -f "$rendered" $(demo_vault_args) --json >/dev/null; then
    echo "MDL apply failed: $label" >&2
    exit 1
  fi
}

demo_apply_mdl_template() {
  local template="$1"
  local rendered="$WORKDIR/$(basename "$template")"
  python3 "$RENDER_MDL" "$template" "$rendered"
  if [ -n "${DEMO_DIST:-}" ]; then
    demo_dist_save_mdl "$rendered"
  fi
  demo_apply_mdl_file "$rendered" "$template"
}

demo_apply_mdl_template_optional() {
  local template="$1"
  local rendered="$WORKDIR/$(basename "$template")"
  python3 "$RENDER_MDL" "$template" "$rendered"
  if ! "$VIVARCUS" component apply-mdl --confirm -f "$rendered" $(demo_vault_args) --json >/dev/null; then
    echo "note: optional MDL skipped ($template)" >&2
  fi
}

demo_build_wasm() {
  (cd "$EXAMPLE_DIR" && "$VIVARCUS_SDK" build . -o "$WORKDIR/action.wasm")
  demo_load_describe
}

# Example names live in main.go Meta() and mdl/*.mdl. demo.env is optional.
demo_apply_example_defaults() {
  TITLE_FIELD="${TITLE_FIELD:-title__c}"
  case "$(basename "$EXAMPLE_DIR")" in
    02-update-field)
      EXPECTED_TITLE="${EXPECTED_TITLE:-from-sdk}"
      ;;
    03-confirm-dialog)
      EXPECTED_TITLE="${EXPECTED_TITLE:-confirmed-by-sdk}"
      DEMO_SETUP_NOTE="${DEMO_SETUP_NOTE:-UI: All Actions → Confirm Update (dialog + banner)}"
      DEMO_VERIFY_NOTE="${DEMO_VERIFY_NOTE:-note: OnPreExecute dialog / OnPostExecute banner require UI manual test}"
      ;;
    04-lifecycle-entry)
      OBJECT="${OBJECT:-demo_request__c}"
      LIFECYCLE="${LIFECYCLE:-demo_request_lc__c}"
      WF_ACTION="${WF_ACTION:-stamp_on_enter_wf__c}"
      WF_CANCEL="${WF_CANCEL:-stamp_on_enter_wf_cancel__c}"
      SUBMIT_ACTION="${SUBMIT_ACTION:-submit__c}"
      IN_REVIEW_STATE="${IN_REVIEW_STATE:-in_review__c}"
      EXPECTED_TITLE="${EXPECTED_TITLE:-stamped-on-enter}"
      ;;
    multi-component)
      SET_TITLE_ACTION="${SET_TITLE_ACTION:-set_title__c}"
      CLEAR_TITLE_ACTION="${CLEAR_TITLE_ACTION:-clear_title__c}"
      SET_TITLE_EXPECTED="${SET_TITLE_EXPECTED:-from-sdk}"
      CLEAR_TITLE_EXPECTED="${CLEAR_TITLE_EXPECTED-}"
      EXPECT_NAME_SUFFIX="${EXPECT_NAME_SUFFIX:--trig}"
      LIFECYCLE="${LIFECYCLE:-sdk_demo_lc__c}"
      SUBMIT_ACTION="${SUBMIT_ACTION:-submit__c}"
      ENTRY_TITLE_EXPECTED="${ENTRY_TITLE_EXPECTED:-stamped-on-enter}"
      STAMP_ON_ENTER_MATCH="${STAMP_ON_ENTER_MATCH:-StampOnEnter}"
      ;;
  esac
}

# object / object_action / Recordaction name come from wasm __sdk_describe.
demo_load_describe() {
  local wasm="$WORKDIR/action.wasm"
  if [ ! -f "$wasm" ]; then
    wasm="$EXAMPLE_DIR/action.wasm"
  fi
  if [ ! -f "$wasm" ]; then
    echo "wasm not found (build first): $WORKDIR/action.wasm" >&2
    exit 1
  fi
  DEMO_DESCRIBE_JSON=$("$VIVARCUS_SDK" describe "$wasm")
  export DEMO_DESCRIBE_JSON
}

demo_describe_field() {
  python3 -c '
import json, sys
d = json.load(sys.stdin)
field = sys.argv[1]
match = sys.argv[2] if len(sys.argv) > 2 else ""
acts = d.get("actions") or []
if not acts:
    raise SystemExit("wasm describe returned no actions")
chosen = None
if not match:
    chosen = acts[0]
else:
    for a in acts:
        name = a.get("component_name") or ""
        obj_act = a.get("object_action") or ""
        if obj_act == match or name == match or name.endswith("." + match):
            chosen = a
            break
    if chosen is None:
        raise SystemExit("no action matching %r in wasm describe" % match)
print(chosen.get(field) or "")
' "$1" "${2:-}"
}

demo_fill_from_describe() {
  if [ -z "${DEMO_DESCRIBE_JSON:-}" ]; then
    demo_load_describe
  fi
  if [ -z "${ACTION_FQN:-}" ]; then
    ACTION_FQN=$(printf '%s' "$DEMO_DESCRIBE_JSON" | demo_describe_field component_name)
  fi
  if [ -z "${OBJECT:-}" ]; then
    OBJECT=$(printf '%s' "$DEMO_DESCRIBE_JSON" | demo_describe_field object)
  fi
  if [ -z "${OBJECT_ACTION:-}" ]; then
    OBJECT_ACTION=$(printf '%s' "$DEMO_DESCRIBE_JSON" | demo_describe_field object_action)
  fi
  if [ -z "$ACTION_FQN" ]; then
    echo "could not read Recordaction FQN from wasm describe" >&2
    exit 1
  fi
  export ACTION_FQN OBJECT OBJECT_ACTION
}

demo_fill_multi_component_from_describe() {
  if [ -z "${DEMO_DESCRIBE_JSON:-}" ]; then
    demo_load_describe
  fi
  SET_TITLE_ACTION="${SET_TITLE_ACTION:-set_title__c}"
  CLEAR_TITLE_ACTION="${CLEAR_TITLE_ACTION:-clear_title__c}"
  if [ -z "${SET_TITLE_FQN:-}" ]; then
    SET_TITLE_FQN=$(printf '%s' "$DEMO_DESCRIBE_JSON" | demo_describe_field component_name "$SET_TITLE_ACTION")
  fi
  if [ -z "${CLEAR_TITLE_FQN:-}" ]; then
    CLEAR_TITLE_FQN=$(printf '%s' "$DEMO_DESCRIBE_JSON" | demo_describe_field component_name "$CLEAR_TITLE_ACTION")
  fi
  if [ -z "${OBJECT:-}" ]; then
    OBJECT=$(printf '%s' "$DEMO_DESCRIBE_JSON" | demo_describe_field object "$SET_TITLE_ACTION")
  fi
  STAMP_ON_ENTER_MATCH="${STAMP_ON_ENTER_MATCH:-StampOnEnter}"
  if [ -z "${ENTRY_ACTION_FQN:-}" ]; then
    ENTRY_ACTION_FQN=$(printf '%s' "$DEMO_DESCRIBE_JSON" | demo_describe_field component_name "$STAMP_ON_ENTER_MATCH")
  fi
  if [ -z "$ENTRY_ACTION_FQN" ]; then
    echo "could not read entry Recordaction FQN (match=${STAMP_ON_ENTER_MATCH}) from wasm describe" >&2
    exit 1
  fi
  export SET_TITLE_ACTION CLEAR_TITLE_ACTION SET_TITLE_FQN CLEAR_TITLE_FQN OBJECT ENTRY_ACTION_FQN
}

demo_deploy_vpk() {
  demo_deploy_vpk_at "$EXAMPLE_DIR" "$WORKDIR/action.vpk"
}

demo_deploy_vpk_at() {
  local module_dir="$1"
  local vpk="$2"
  shift 2
  "$PACKAGE_VPK" "$@" "$module_dir" "$vpk" >/dev/null
  if [ -n "${DEMO_DIST:-}" ]; then
    demo_dist_publish_vpk "$vpk"
  fi
  local pkg
  pkg=$("$VIVARCUS" package import -f "$vpk" $(demo_vault_args) --json \
    | python3 -c "import json,sys; d=json.load(sys.stdin); print(d.get('package_id') or d.get('id') or '')")
  "$VIVARCUS" package validate "$pkg" $(demo_vault_args) --json >/dev/null
  "$VIVARCUS" package deploy "$pkg" --confirm $(demo_vault_args) --json >/dev/null
}

demo_render_mdl() {
  local template="$1"
  local rendered="$2"
  python3 "$RENDER_MDL" "$template" "$rendered"
  if [ -n "${DEMO_DIST:-}" ]; then
    demo_dist_save_mdl "$rendered"
  fi
}

demo_object_component_arg() {
  local rendered="$1"
  echo "--component" "10:Object:${OBJECT}:${rendered}"
}

demo_export_user_action_env() {
  export DEMO_OBJECT="${OBJECT:?OBJECT required (wasm Meta.Object)}"
  export DEMO_OBJECT_ACTION="${OBJECT_ACTION:?OBJECT_ACTION required (wasm Meta.ObjectAction)}"
  export DEMO_TITLE_FIELD="${TITLE_FIELD:-title__c}"
  export DEMO_EXPECTED_TITLE="${EXPECTED_TITLE:?EXPECTED_TITLE required}"
}

demo_export_lifecycle_verify_env() {
  export DEMO_OBJECT="${OBJECT:?OBJECT required}"
  export DEMO_LIFECYCLE="${LIFECYCLE:?LIFECYCLE required}"
  export DEMO_WF_ACTION="${WF_ACTION:?WF_ACTION required}"
  export DEMO_WF_CANCEL="${WF_CANCEL:?WF_CANCEL required}"
  export DEMO_TITLE_FIELD="${TITLE_FIELD:-title__c}"
  export DEMO_EXPECTED_TITLE="${EXPECTED_TITLE:-stamped-on-enter}"
  export DEMO_SUBMIT_ACTION="${SUBMIT_ACTION:-submit__c}"
  export DEMO_IN_REVIEW_STATE="${IN_REVIEW_STATE:-in_review__c}"
}

demo_export_multi_component_verify_env() {
  export DEMO_OBJECT="${OBJECT:?OBJECT required}"
  export DEMO_TITLE_FIELD="${TITLE_FIELD:-title__c}"
  export DEMO_SET_TITLE_ACTION="${SET_TITLE_ACTION:?SET_TITLE_ACTION required}"
  export DEMO_SET_TITLE_EXPECTED="${SET_TITLE_EXPECTED:?SET_TITLE_EXPECTED required}"
  export DEMO_CLEAR_TITLE_ACTION="${CLEAR_TITLE_ACTION:?CLEAR_TITLE_ACTION required}"
  export DEMO_CLEAR_TITLE_EXPECTED="${CLEAR_TITLE_EXPECTED-}"
  export DEMO_EXPECT_NAME_SUFFIX="${EXPECT_NAME_SUFFIX-}"
  export DEMO_ENTRY_TITLE_EXPECTED="${ENTRY_TITLE_EXPECTED-}"
  export DEMO_SUBMIT_ACTION="${SUBMIT_ACTION:-submit__c}"
}

# 02 / 03 — single user Objectaction on a demo object.
demo_scenario_user_action() {
  demo_require_auth
  demo_build_wasm
  demo_fill_from_describe
  demo_export_user_action_env

  local object_mdl="$DEMO_DIST/mdl/01-object.mdl"
  demo_render_mdl "$MDL_DIR/01-object.mdl" "$object_mdl"

  echo "== deploy VPK (object MDL in components/00010 + gosdk source) =="
  demo_deploy_vpk_at \
    "$EXAMPLE_DIR" "$DEMO_DIST/action.vpk" \
    $(demo_object_component_arg "$object_mdl")

  echo "setup complete (object=${OBJECT} action=${OBJECT_ACTION} fqn=${ACTION_FQN})"
  if [ -n "${DEMO_SETUP_NOTE:-}" ]; then
    echo "$DEMO_SETUP_NOTE"
  fi
  echo "see artifacts under $DEMO_DIST/"
}

# 04 — lifecycle event / entry / workflow / cancel.
demo_scenario_lifecycle() {
  demo_require_auth
  echo "== build wasm =="
  demo_build_wasm
  demo_fill_from_describe
  demo_export_lifecycle_verify_env

  echo "== apply MDL (lifecycle; object MDL ships in VPK components/00010) =="
  demo_apply_mdl_template "$MDL_DIR/02-lifecycle.mdl"
  demo_apply_mdl_template_optional "$MDL_DIR/03-bind-object-lifecycle.mdl"

  local object_mdl="$DEMO_DIST/mdl/01-object.mdl"
  demo_render_mdl "$MDL_DIR/01-object.mdl" "$object_mdl"

  echo "== deploy Recordaction VPK (object + gosdk source) =="
  demo_deploy_vpk_at \
    "$EXAMPLE_DIR" "$DEMO_DIST/action.vpk" \
    $(demo_object_component_arg "$object_mdl")

  demo_apply_mdl_template "$MDL_DIR/04-workflow-action.mdl"
  demo_apply_mdl_template "$MDL_DIR/05-workflow-cancel.mdl"

  echo "setup complete (object=${OBJECT} lifecycle=${LIFECYCLE} action=${ACTION_FQN})"
  echo "see artifacts under $DEMO_DIST/"
}

# multi-component — user actions + lifecycle entry + trigger, one module / one VPK.
demo_scenario_multi_component() {
  demo_require_auth

  echo "== build =="
  demo_build_wasm
  demo_fill_multi_component_from_describe

  echo "== apply MDL (lifecycle entry_action; object MDL ships in VPK components/00010) =="
  ACTION_FQN="${ENTRY_ACTION_FQN:?ENTRY_ACTION_FQN required}"
  export ACTION_FQN
  demo_apply_mdl_template "$MDL_DIR/02-lifecycle.mdl"
  demo_apply_mdl_template_optional "$MDL_DIR/03-bind-object-lifecycle.mdl"
  unset ACTION_FQN

  local object_mdl="$DEMO_DIST/mdl/01-object.mdl"
  demo_render_mdl "$MDL_DIR/01-object.mdl" "$object_mdl"

  echo "== deploy combo VPK (object + gosdk source with actions + trigger) =="
  demo_deploy_vpk_at \
    "$EXAMPLE_DIR" "$DEMO_DIST/action.vpk" \
    $(demo_object_component_arg "$object_mdl")

  echo "setup complete (object=${OBJECT} actions=${SET_TITLE_ACTION},${CLEAR_TITLE_ACTION} entry=${ENTRY_ACTION_FQN})"
  echo "see artifacts under $DEMO_DIST/"
}
