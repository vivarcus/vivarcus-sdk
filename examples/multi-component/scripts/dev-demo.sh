#!/usr/bin/env bash
exec "$(cd "$(dirname "$0")/../../_shared/scripts" && pwd)/dev-demo.sh" "$(cd "$(dirname "$0")/.." && pwd)" "$@"
