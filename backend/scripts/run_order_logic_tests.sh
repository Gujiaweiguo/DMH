#!/usr/bin/env bash

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
API_DIR="$ROOT_DIR/api"

start_ts="$(date +%s)"

echo "=========================================="
echo "Order Logic Test Suite"
echo "=========================================="

cd "$API_DIR"

echo "[1/2] unit + integration (order package)"
go test ./internal/logic/order -count=1 -v

echo "[2/2] smoke regression"
"$ROOT_DIR/scripts/order_smoke.sh"

end_ts="$(date +%s)"
duration="$((end_ts - start_ts))"

echo "=========================================="
echo "Order logic tests passed in ${duration}s"
echo "=========================================="
