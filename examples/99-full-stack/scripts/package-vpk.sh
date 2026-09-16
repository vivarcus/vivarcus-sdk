#!/usr/bin/env bash
# Package a gosdk/ VPK from wasm + sdk_manifest.json.
#
# Usage:
#   ./scripts/package-vpk.sh <action.wasm> <sdk_manifest.json> [output.vpk]
#
# Example:
#   vivarcus-sdk build ../02-update-field -o action.wasm
#   ./scripts/package-vpk.sh ./action.wasm ./action.sdk_manifest.json
set -euo pipefail

WASM="${1:?wasm path required}"
MANIFEST="${2:?sdk_manifest.json path required}"
OUT="${3:-}"

if [ ! -f "$WASM" ]; then
  echo "wasm not found: $WASM" >&2
  exit 1
fi
if [ ! -f "$MANIFEST" ]; then
  echo "manifest not found: $MANIFEST" >&2
  exit 1
fi

ROOT=$(cd "$(dirname "$0")/.." && pwd)
DIST="$ROOT/dist"
rm -rf "$DIST"
mkdir -p "$DIST/gosdk"

cp "$WASM" "$DIST/gosdk/action.wasm"
cp "$MANIFEST" "$DIST/gosdk/sdk_manifest.json"

cat > "$DIST/vaultpackage.xml" <<'EOF'
<?xml version="1.0" encoding="UTF-8"?>
<VaultPackage xmlns="https://www.veevavault.com/schema/vaultpackage">
  <name>Demo SDK Action</name>
  <summary>Customer Record Action (gosdk)</summary>
  <packagetype>migration__v</packagetype>
</VaultPackage>
EOF

PKG_NAME="${OUT:-$DIST/demo-action.vpk}"
(cd "$DIST" && zip -qr "$(basename "$PKG_NAME")" vaultpackage.xml gosdk/)
if [ "$PKG_NAME" != "$DIST/$(basename "$PKG_NAME")" ]; then
  mv "$DIST/$(basename "$PKG_NAME")" "$PKG_NAME"
fi

echo "created $PKG_NAME"
sha256sum "$DIST/gosdk/action.wasm"
