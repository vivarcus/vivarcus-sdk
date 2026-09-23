#!/usr/bin/env bash
# Package an Inbound VPK: optional components/<step>/*.mdl + gosdk/ Go source.
#
# Platform deploy order: all component steps (by step number), then gosdk last.
# gosdk/ contains only .go; module path is written to vaultpackage.xml (ADR-21).
# Local go.mod is required in the module dir for Go tooling and to copy <module>.
#
# Usage:
#   package-vpk.sh [options] <module-dir> [output.vpk]
#
# Options (repeatable):
#   --component <step>:<Type>:<name>:<mdl-file>
#   --name <vaultpackage name>
#   --summary <vaultpackage summary>
#   --replace   gosdk deployment_option=replace_all (drop other files on the vault tree)
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

PKG_NAME="${PKG_NAME:-Demo SDK Action}"
PKG_SUMMARY="${PKG_SUMMARY:-Customer Record Action (gosdk + optional MDL)}"
DEPLOY_OPTION="incremental"
COMPONENTS=()
POSITIONAL=()

while [ $# -gt 0 ]; do
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
    --replace)
      DEPLOY_OPTION="replace_all"
      shift
      ;;
    -h|--help)
      sed -n '1,20p' "$0"
      exit 0
      ;;
    --)
      shift
      POSITIONAL+=("$@")
      break
      ;;
    --*)
      echo "unknown option: $1" >&2
      exit 2
      ;;
    *)
      POSITIONAL+=("$1")
      shift
      ;;
  esac
done

MODULE_DIR="${POSITIONAL[0]:?module dir required}"
OUT="${POSITIONAL[1]:-}"

if [ ! -d "$MODULE_DIR" ]; then
  echo "module dir not found: $MODULE_DIR" >&2
  exit 1
fi
if [ ! -f "$MODULE_DIR/go.mod" ]; then
  echo "go.mod required in module dir (local DX; not packed into VPK): $MODULE_DIR" >&2
  exit 1
fi

MODULE_PATH="$(awk '/^module / {print $2; exit}' "$MODULE_DIR/go.mod")"
if [ -z "$MODULE_PATH" ]; then
  echo "go.mod has no module path: $MODULE_DIR/go.mod" >&2
  exit 1
fi
MODULE_PATH_XML="$(python3 -c 'import xml.sax.saxutils,sys; print(xml.sax.saxutils.escape(sys.argv[1]))' "$MODULE_PATH")"

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

is_local_main_harness() {
  python3 - "$1" <<'PY'
import re, sys
src = open(sys.argv[1], encoding="utf-8").read()
src = re.sub(r"/\*.*?\*/", "", src, flags=re.S)
src = re.sub(r"//.*?$", "", src, flags=re.M)
src = re.sub(r"`(?:\\.|[^`])*`", "``", src)
src = re.sub(r'"(?:\\.|[^"\\])*"', '""', src)
if not re.search(r"(?m)^\s*package\s+main\s*$", src):
    sys.exit(1)
if re.search(r"(?m)^\s*(type|const|var)\s+", src):
    sys.exit(1)
funcs = re.findall(r"(?m)^\s*func\s+(?:\([^)]*\)\s*)?(\w+)\s*\(", src)
if funcs not in ([], ["main"]):
    sys.exit(1)
for m in re.finditer(r'(?m)^\s*import\s+(?:(\w+)\s+)?"', src):
    if m.group(1) != "_":
        sys.exit(1)
block = re.search(r"import\s*\((.*?)\)", src, re.S)
if block:
    for line in block.group(1).splitlines():
        line = line.strip()
        if not line:
            continue
        if line.startswith("_"):
            continue
        sys.exit(1)
sys.exit(0)
PY
}

while IFS= read -r -d '' file; do
  rel="${file#"$MODULE_DIR"/}"
  case "$rel" in
    action.wasm|*.wasm) continue ;;
    dist/*) continue ;;
  esac
  base="$(basename "$rel")"
  case "$base" in
    zz_generated_reactor.go|*_test.go) continue ;;
  esac
  # A root main.go that only blank-imports packages is local DX (same as go.mod).
  # Entry files named after their type are still packed.
  # Rule matches sourcetree.IsLocalMainHarness.
  if [ "$rel" = "main.go" ] && is_local_main_harness "$file"; then
    continue
  fi
  mkdir -p "$DIST/gosdk/$(dirname "$rel")"
  cp "$file" "$DIST/gosdk/$rel"
done < <(find "$MODULE_DIR" -type f -name '*.go' ! -path '*/dist/*' -print0)

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
  <gosdk>
    <deployment_option>${DEPLOY_OPTION}</deployment_option>
    <module>${MODULE_PATH_XML}</module>
  </gosdk>
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
echo "  gosdk module: $MODULE_PATH" >&2
if [ "${#COMPONENTS[@]}" -gt 0 ]; then
  echo "  components: ${#COMPONENTS[@]} mdl file(s)" >&2
fi
rm -rf "$DIST"
