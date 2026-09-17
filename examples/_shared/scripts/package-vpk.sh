#!/usr/bin/env bash
# Package an Inbound VPK: optional components/<step>/*.mdl + gosdk/*.wasm.
#
# Platform deploy order: all component steps (by step number), then gosdk last.
# Import scans gosdk/*.wasm and projects Recordactions from __sdk_describe.
#
# Usage:
#   package-vpk.sh [options] <action.wasm> [output.vpk]
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

WASM="${1:?wasm path required}"
OUT="${2:-}"

if [ ! -f "$WASM" ]; then
  echo "wasm not found: $WASM" >&2
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

cp "$WASM" "$DIST/gosdk/action.wasm"

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
( cd "$DIST" && zip -qr "$ZIP_BASE" vaultpackage.xml gosdk/ )
if [ -d "$DIST/components" ]; then
  ( cd "$DIST" && zip -qr "$ZIP_BASE" components/ )
fi
if [ "$PKG_PATH" != "$DIST/$ZIP_BASE" ]; then
  mv "$DIST/$ZIP_BASE" "$PKG_PATH"
fi

echo "created $PKG_PATH" >&2
if [ "${#COMPONENTS[@]}" -gt 0 ]; then
  echo "  components: ${#COMPONENTS[@]} mdl file(s)" >&2
fi
sha256sum "$DIST/gosdk/action.wasm" >&2
rm -rf "$DIST"
