#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

cd "$ROOT_DIR/api"

echo "[order-smoke] running minimal order smoke tests"
go test ./internal/logic/order -run 'TestOrderIntegration_CreateVerifyScanConsistency|TestOrderIntegration_ConcurrentDuplicateCreateGuard' -count=1 -v
