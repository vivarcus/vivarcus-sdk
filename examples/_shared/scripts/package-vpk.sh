#!/usr/bin/env bash
# Package an Inbound VPK: optional components/<step>/*.mdl + gosdk/ Go source.
#
# Platform deploy order: all component steps (by step number), then gosdk last.
# Import validates gosdk/ Go source; deploy compiles on the platform (ADR-21).
#
# Usage:
#   package-vpk.sh [options] <module-dir> [output.vpk]
#
# Options (repeatable):
#   --component <step>:<Type>:<name>:<mdl-file>
#   --name <vaultpackage name>
#   --summary <vaultpackage summary>
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
# shellcheck source=/dev/null
. "${SCRIPT_DIR}/lib-sdk-module-version.sh"

module_go_toolchain_env() {
  local mod="$1"
  local ver
  ver="$(sed -n 's/^go //p' "$mod" | head -1)"
  if [ -n "$ver" ]; then
    echo "GOTOOLCHAIN=go${ver}"
  else
    echo "GOTOOLCHAIN=local"
  fi
}

# Rewrite require to the Go module tag (dual-tag: v26.3.3-N from platform v26R3.3-N),
# run go mod tidy with a local SDK replace, then strip replace for customer VPK (ADR-21).
prepare_gosdk_mod_for_vpk() {
  local gosdk_dir="$1"
  local repo_root="$2"
  local module_root="${3:-$repo_root}"
  local mod="${gosdk_dir}/go.mod"
  if [ ! -f "$mod" ]; then
    return 0
  fi

  local sdk_version tag_root="${VIVARCUS_REPO_ROOT:-$repo_root}"
  if ! sdk_version="$(resolve_sdk_module_version "$tag_root")"; then
    echo "prepare_gosdk_mod_for_vpk: could not resolve github.com/vivarcus/vivarcus-sdk version (set VIVARCUS_PLATFORM_TAG or VIVARCUS_SDK_MODULE_REF)" >&2
    return 1
  fi
  local gosdk_go
  gosdk_go="$(sed -n 's/^go //p' "${repo_root}/sdk/go.mod" | head -1)"
  [ -n "$gosdk_go" ] || gosdk_go="1.26.2"

  python3 - "$mod" "$sdk_version" "$gosdk_go" <<'PY'
import re, sys
path, sdk_version, go_ver = sys.argv[1], sys.argv[2], sys.argv[3]
sdk_modules = (
    "github.com/vivarcus/vivarcus-sdk",
)
lines = open(path, encoding="utf-8").read().splitlines()
out = []
in_require = False
seen_sdk = False
for line in lines:
    stripped = line.strip()
    if stripped.startswith("go "):
        out.append(f"go {go_ver}")
        continue
    if stripped == "require (":
        in_require = True
        out.append(line)
        continue
    if in_require and stripped == ")":
        in_require = False
        if not seen_sdk:
            out.append(f"\t{sdk_modules[0]} {sdk_version}")
        out.append(line)
        continue
    if stripped.startswith("require "):
        fields = stripped.split()
        if len(fields) >= 3 and fields[1] in sdk_modules:
            out.append(f"require {fields[1]} {sdk_version}")
            seen_sdk = True
            continue
    if in_require:
        fields = stripped.split()
        if len(fields) >= 2 and fields[0] in sdk_modules:
            out.append(f"\t{fields[0]} {sdk_version}")
            seen_sdk = True
            continue
    out.append(line)
if not seen_sdk:
    if out and out[-1].strip() != "":
        out.append("")
    out.append(f"require {sdk_modules[0]} {sdk_version}")
with open(path, "w", encoding="utf-8") as f:
    f.write("\n".join(out).rstrip() + "\n")
PY

  {
    echo ""
    echo "replace github.com/vivarcus/vivarcus-sdk => ${module_root}"
  } >> "$mod"

  (
    cd "$gosdk_dir"
    # shellcheck disable=SC1090
    export $(module_go_toolchain_env "$mod")
    go mod tidy
  )

  python3 - "$mod" <<'PY'
import re, sys
path = sys.argv[1]
lines = open(path, encoding="utf-8").read().splitlines()
out = []
skip = False
for line in lines:
    if skip:
        if line.strip() == ")":
            skip = False
        continue
    if re.match(r"^replace\s*\(", line):
        skip = True
        continue
    if re.match(r"^replace\s+", line):
        continue
    out.append(line)
with open(path, "w", encoding="utf-8") as f:
    f.write("\n".join(out).rstrip() + "\n")
PY

  populate_gosdk_sum_from_proxy "$gosdk_dir" "$repo_root"

  if [ ! -f "${gosdk_dir}/go.sum" ]; then
    echo "prepare_gosdk_mod_for_vpk: warning: no go.sum (SDK resolved via platform injectGuestReplaces)" >&2
  fi
}

populate_gosdk_sum_from_proxy() {
  local gosdk_dir="$1"
  local repo_root="${2:-}"
  local ref helper resolved
  ref="$(normalize_sdk_module_list_ref "$(resolve_sdk_module_ref "${VIVARCUS_REPO_ROOT:-$repo_root}")")"

  resolved="$(GOPROXY="${GOPROXY:-direct}" go list -m "github.com/vivarcus/vivarcus-sdk@${ref}" 2>/dev/null | awk '{print $2}')" || return 0
  [ -n "$resolved" ] || return 0

  helper="$(mktemp -d)"
  local gosdk_go="1.26.2"
  if [ -n "$repo_root" ] && [ -f "${repo_root}/sdk/go.mod" ]; then
    gosdk_go="$(sed -n 's/^go //p' "${repo_root}/sdk/go.mod" | head -1)"
  fi
  cat > "${helper}/go.mod" <<EOF
module example.com/vivarcus-sdk-sum

go ${gosdk_go}

require github.com/vivarcus/vivarcus-sdk ${resolved}
EOF
  cat > "${helper}/main.go" <<'EOF'
package main

import _ "github.com/vivarcus/vivarcus-sdk/action"

func main() {}
EOF
  if ( cd "$helper" && GOPROXY="${GOPROXY:-direct}" go mod tidy >/dev/null 2>&1 ) && [ -f "${helper}/go.sum" ]; then
    cp "${helper}/go.sum" "${gosdk_dir}/go.sum"
  fi
  rm -rf "$helper"
}

PKG_NAME="${PKG_NAME:-Demo SDK Action}"
PKG_SUMMARY="${PKG_SUMMARY:-Customer Record Action (gosdk + optional MDL)}"
COMPONENTS=()

while [ $# -gt 0 ] && [[ "$1" == --* ]]; do
  case "$1" in
    --component)
      COMPONENTS+=("${2:?--component requires step:Type:name:file}")
      shift 2
      ;;
    --name)
      PKG_NAME="$2"
      shift 2
      ;;
    --summary)
      PKG_SUMMARY="$2"
      shift 2
      ;;
    -h|--help)
      sed -n '1,16p' "$0"
      exit 0
      ;;
    *)
      echo "unknown option: $1" >&2
      exit 2
      ;;
  esac
done

MODULE_DIR="${1:?module dir required}"
OUT="${2:-}"

if [ ! -d "$MODULE_DIR" ]; then
  echo "module dir not found: $MODULE_DIR" >&2
  exit 1
fi
if [ ! -f "$MODULE_DIR/go.mod" ]; then
  echo "go.mod required in module dir: $MODULE_DIR" >&2
  exit 1
fi

DIST="${TMPDIR:-/tmp}/vivarcus-vpk-$$"
rm -rf "$DIST"
mkdir -p "$DIST/gosdk"

for spec in "${COMPONENTS[@]}"; do
  IFS=: read -r step ctype cname mdlfile <<< "$spec"
  if [ -z "$step" ] || [ -z "$ctype" ] || [ -z "$cname" ] || [ -z "$mdlfile" ]; then
    echo "invalid --component (want step:Type:name:file): $spec" >&2
    exit 1
  fi
  if [ ! -f "$mdlfile" ]; then
    echo "mdl not found: $mdlfile" >&2
    exit 1
  fi
  folder=$(printf "components/%05d" "$step")
  base="${ctype}.${cname}"
  mkdir -p "$DIST/$folder"
  cp "$mdlfile" "$DIST/$folder/${base}.mdl"
  md5=$(md5sum "$mdlfile" | awk '{print $1}')
  printf '%s %s\n' "$md5" "$base" > "$DIST/$folder/${base}.md5"
done

while IFS= read -r -d '' file; do
  rel="${file#"$MODULE_DIR"/}"
  case "$rel" in
    action.wasm|*.wasm) continue ;;
    dist/*) continue ;;
  esac
  case "$(basename "$rel")" in
    zz_generated_reactor.go) continue ;;
  esac
  mkdir -p "$DIST/gosdk/$(dirname "$rel")"
  cp "$file" "$DIST/gosdk/$rel"
done < <(find "$MODULE_DIR" -type f \( -name '*.go' -o -name 'go.mod' -o -name 'go.sum' \) ! -path '*/dist/*' -print0)

# VPK gosdk/ must not ship replace directives (ADR-21 validate).
if [ -f "$DIST/gosdk/go.mod" ]; then
  python3 - "$DIST/gosdk/go.mod" <<'PY'
import re, sys
path = sys.argv[1]
lines = open(path, encoding="utf-8").read().splitlines()
out = []
skip = False
for line in lines:
    if skip:
        if line.strip() == ")":
            skip = False
        continue
    if re.match(r"^replace\s*\(", line):
        skip = True
        continue
    if re.match(r"^replace\s+", line):
        continue
    out.append(line)
with open(path, "w", encoding="utf-8") as f:
    f.write("\n".join(out).rstrip() + "\n")
PY
fi

# Pin SDK require to Go module tag (dual-tag) and emit go.sum (customer VPK path).
resolve_sdk_module_root() {
  local dir="$1"
  while [ "$dir" != "/" ]; do
    if [ -f "$dir/go.mod" ] && grep -q '^module github.com/vivarcus/vivarcus-sdk' "$dir/go.mod" 2>/dev/null; then
      echo "$dir"
      return 0
    fi
    dir="$(dirname "$dir")"
  done
  if [ -n "${VIVARCUS_SDK_EXPORT:-}" ]; then
    echo "$VIVARCUS_SDK_EXPORT"
    return 0
  fi
  echo "resolve_sdk_module_root: github.com/vivarcus/vivarcus-sdk module root not found" >&2
  return 1
}
REPO_ROOT="$(git -C "$SCRIPT_DIR" rev-parse --show-toplevel 2>/dev/null || true)"
if [ -z "$REPO_ROOT" ]; then
  REPO_ROOT="${VIVARCUS_SDK_EXPORT:-}"
fi
SDK_MODULE_ROOT="$(resolve_sdk_module_root "$SCRIPT_DIR")"
if [ -z "$REPO_ROOT" ]; then
  REPO_ROOT="$SDK_MODULE_ROOT"
fi
prepare_gosdk_mod_for_vpk "$DIST/gosdk" "$REPO_ROOT" "$SDK_MODULE_ROOT"

if ! find "$DIST/gosdk" -name '*.go' -print -quit | grep -q .; then
  echo "no .go files copied from $MODULE_DIR" >&2
  exit 1
fi

cat > "$DIST/vaultpackage.xml" <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<VaultPackage xmlns="https://www.veevavault.com/schema/vaultpackage">
  <name>${PKG_NAME}</name>
  <summary>${PKG_SUMMARY}</summary>
  <packagetype>migration__v</packagetype>
</VaultPackage>
EOF

PKG_PATH="${OUT:-$DIST/demo-action.vpk}"
ZIP_BASE="$(basename "$PKG_PATH")"
( cd "$DIST" && {
  find vaultpackage.xml gosdk -type f -print
  if [ -d components ]; then
    find components -type f -print
  fi
} | zip -qr "$ZIP_BASE" -@ )
if [ "$PKG_PATH" != "$DIST/$ZIP_BASE" ]; then
  mv "$DIST/$ZIP_BASE" "$PKG_PATH"
fi

echo "created $PKG_PATH" >&2
if [ "${#COMPONENTS[@]}" -gt 0 ]; then
  echo "  components: ${#COMPONENTS[@]} mdl file(s)" >&2
fi
rm -rf "$DIST"
