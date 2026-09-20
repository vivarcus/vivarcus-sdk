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
  esac
  case "$(basename "$rel")" in
    zz_generated_reactor.go) continue ;;
  esac
  mkdir -p "$DIST/gosdk/$(dirname "$rel")"
  cp "$file" "$DIST/gosdk/$rel"
done < <(find "$MODULE_DIR" -type f \( -name '*.go' -o -name 'go.mod' -o -name 'go.sum' \) -print0)

# VPK gosdk/ must not ship monorepo replace directives (ADR-21 validate).
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
