# Shared helpers for sdk/examples/*/deploy.sh (source, do not execute).
# Requires: bash, python3, vivarcus CLI + Vault session (see deploy.env.example).
set -euo pipefail

deploy_init() {
  local self="${1:?deploy_init: path to deploy.sh}"
  EXAMPLE_ROOT="$(cd "$(dirname "$self")" && pwd)"
  EXAMPLES_ROOT="$(cd "$EXAMPLE_ROOT/.." && pwd)"
  SHARED_SCRIPTS="$(cd "$EXAMPLES_ROOT/_shared/scripts" && pwd)"
  VIVARCUS="${VIVARCUS_BIN:-vivarcus}"
  cd "$EXAMPLE_ROOT"
  deploy_load_env
}

deploy_load_env() {
  if [ -f "$EXAMPLES_ROOT/deploy.env" ]; then
    set -a
    # shellcheck disable=SC1091
    source "$EXAMPLES_ROOT/deploy.env"
    set +a
  fi
  if [ -f "$EXAMPLE_ROOT/deploy.env" ]; then
    set -a
    # shellcheck disable=SC1091
    source "$EXAMPLE_ROOT/deploy.env"
    set +a
  fi
}

deploy_auth_configured() {
  if [ -n "${VIVARCUS_TOKEN:-}" ] && [ -n "${VIVARCUS_ENDPOINT:-}" ] && [ -n "${VIVARCUS_VAULT:-}" ]; then
    return 0
  fi
  "$VIVARCUS" auth status >/dev/null 2>&1
}

deploy_print_not_logged_in() {
  cat >&2 <<EOF
需要先连接 Vault（deploy.sh 不会保存或传递密码）。

  方式 A — 本机首次登录（会提示账号/密码或 SSO，写入 vivarcus 本地配置）：
    $VIVARCUS auth login --endpoint https://<你的租户>.vivarcus.com
    $VIVARCUS config set default_vault <vault_id>

  方式 B — 使用 session（CI / Agent）：复制 deploy.env.example 为 examples/deploy.env，填写：
    VIVARCUS_ENDPOINT、VIVARCUS_VAULT、VIVARCUS_TOKEN（token 是登录后的会话，不是密码）

  说明：docs/01-prerequisites.md
EOF
}

deploy_show_target() {
  if [ -n "${VIVARCUS_ENDPOINT:-}" ] && [ -n "${VIVARCUS_VAULT:-}" ]; then
    echo "target: $VIVARCUS_ENDPOINT (vault $VIVARCUS_VAULT)"
    return
  fi
  local out
  out=$("$VIVARCUS" auth status --json 2>/dev/null) || return
  python3 -c '
import json, sys
d = json.load(sys.stdin)
ep = d.get("endpoint") or "?"
vault = d.get("default_vault") or d.get("vault_id") or "?"
user = d.get("username") or d.get("user") or ""
line = f"target: {ep} (vault {vault})"
if user:
    line += f" user {user}"
print(line)
' <<<"$out"
}

deploy_require_cli() {
  if ! command -v "$VIVARCUS" >/dev/null 2>&1; then
    echo "need vivarcus on PATH (or set VIVARCUS_BIN)" >&2
    exit 1
  fi
  if ! deploy_auth_configured; then
    deploy_print_not_logged_in
    exit 1
  fi
  deploy_show_target
}

deploy_cli_json() {
  local out
  if ! out=$("$VIVARCUS" "$@" --json 2>&1); then
    echo "$out" >&2
    return 1
  fi
  if ! python3 -c '
import json, sys
d = json.load(sys.stdin)
s = d.get("responseStatus") or ""
if s and s != "SUCCESS":
    print(json.dumps(d), file=sys.stderr)
    sys.exit(1)
' <<<"$out"; then
    echo "$out" >&2
    return 1
  fi
  printf '%s' "$out"
}

deploy_apply_mdl() {
  local f="${1:?mdl file}"
  deploy_cli_json component apply-mdl --confirm -f "$f" >/dev/null
  echo "apply-mdl: $f"
}

deploy_sdk_put() {
  local f="${1:?go file}"
  local out
  out=$(deploy_cli_json sdk put -f "$f")
  deploy_wait_code_compile "$out"
  echo "sdk put: $f"
}

# Block until a queued sdk put compile is SUCCESS. A body with no job_status, or
# job_status SUCCESS, is already final (CLI polled, or an older synchronous PUT).
deploy_wait_code_compile() {
  local payload="${1:?sdk put json}"
  VIVARCUS_COMPILE_PAYLOAD="$payload" python3 - "$VIVARCUS" <<'PY'
import json, os, subprocess, sys, time, urllib.error, urllib.request

vivarcus = sys.argv[1]
try:
    body = json.loads(os.environ["VIVARCUS_COMPILE_PAYLOAD"])
except json.JSONDecodeError as exc:
    print(f"sdk put: invalid json: {exc}", file=sys.stderr)
    sys.exit(1)

def cfg(key):
    try:
        out = subprocess.check_output(
            [vivarcus, "config", "get", key], text=True, stderr=subprocess.DEVNULL
        )
    except (subprocess.CalledProcessError, FileNotFoundError):
        return ""
    return out.strip()

status = (body.get("job_status") or "").strip()
if status in ("", "SUCCESS"):
    sys.exit(0)
url = (body.get("url") or "").strip()
if not url:
    print(f"sdk put: compile job url missing (job_status={status})", file=sys.stderr)
    sys.exit(1)
endpoint = (os.environ.get("VIVARCUS_ENDPOINT") or cfg("endpoint")).rstrip("/")
token = (os.environ.get("VIVARCUS_TOKEN") or cfg("token")).strip()
vault = (os.environ.get("VIVARCUS_VAULT") or cfg("default_vault")).strip()
if not endpoint or not token or not vault:
    print("sdk put: need endpoint, token, and vault to wait for compile", file=sys.stderr)
    sys.exit(1)
deadline = time.time() + 30 * 60
while status not in ("", "SUCCESS"):
    if status in ("FAILURE", "CANCELLED"):
        print(f"sdk compile {status}: {json.dumps(body)}", file=sys.stderr)
        sys.exit(1)
    if time.time() > deadline:
        print(f"sdk compile timed out ({status})", file=sys.stderr)
        sys.exit(1)
    time.sleep(2)
    req = urllib.request.Request(
        endpoint + url,
        headers={"Authorization": "Bearer " + token, "X-Vault-Id": vault},
    )
    try:
        with urllib.request.urlopen(req, timeout=60) as resp:
            raw = resp.read().decode()
    except urllib.error.HTTPError as exc:
        print(exc.read().decode(), file=sys.stderr)
        sys.exit(1)
    body = json.loads(raw)
    if body.get("responseStatus") == "FAILURE":
        print(raw, file=sys.stderr)
        sys.exit(1)
    status = (body.get("job_status") or "").strip()
    nxt = (body.get("url") or "").strip()
    if nxt:
        url = nxt
PY
}

# Usage: deploy_render_mdl <src> <dst> [VAR=val ...]
deploy_render_mdl() {
  local src="${1:?}" dst="${2:?}"
  shift 2
  mkdir -p "$(dirname "$dst")"
  if [ "$#" -gt 0 ]; then
    env "$@" python3 "$SHARED_SCRIPTS/render-mdl.py" "$src" "$dst"
  else
    python3 "$SHARED_SCRIPTS/render-mdl.py" "$src" "$dst"
  fi
  echo "render-mdl: $src -> $dst"
}

deploy_vpk() {
  local vpk="${1:?path to .vpk}"
  local out pkg_id
  out=$(deploy_cli_json package import -f "$vpk")
  pkg_id=$(python3 -c 'import json,sys; d=json.loads(sys.argv[1]); print(d.get("package_id") or d.get("id") or "")' "$out")
  if [ -z "$pkg_id" ]; then
    echo "package import: missing package_id: $out" >&2
    exit 1
  fi
  echo "package import: $pkg_id"
  deploy_cli_json package validate "$pkg_id" >/dev/null
  echo "package validate: ok"
  deploy_cli_json package deploy "$pkg_id" --confirm >/dev/null
  echo "package deploy: ok"
}
