#!/usr/bin/env bash
# Install ov-sdk from a GitHub Release tarball.
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/vivarcus/vivarcus-sdk/main/scripts/install-ov-sdk.sh | bash
#   VERSION=v0.1.0 bash install-ov-sdk.sh
set -euo pipefail

VERSION="${VERSION:-latest}"
INSTALL_DIR="${INSTALL_DIR:-$HOME/.local/bin}"
REPO="vivarcus/vivarcus-sdk"

mkdir -p "$INSTALL_DIR"

if [ "$VERSION" = "latest" ]; then
  URL=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" \
    | grep -o 'https://[^"]*ov-sdk[^"]*linux-amd64[^"]*' | head -1 || true)
else
  URL="https://github.com/${REPO}/releases/download/${VERSION}/ov-sdk-linux-amd64"
fi

if [ -z "${URL:-}" ]; then
  cat <<'EOF'
No ov-sdk release found for this repository yet.

Download ov-sdk from your Vivarcus Vault release notes, or contact Vivarcus support
for the build toolchain matching your Vault version.

See docs/01-prerequisites.md
EOF
  exit 1
fi

TMP=$(mktemp)
curl -fsSL "$URL" -o "$TMP"
chmod +x "$TMP"
mv "$TMP" "$INSTALL_DIR/ov-sdk"
echo "installed ov-sdk to $INSTALL_DIR/ov-sdk"
