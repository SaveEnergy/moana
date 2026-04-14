#!/usr/bin/env bash
# Development: Vite watch (assets) + Air (rebuild Go on change, run serve).
# Requires: Go, bun. Optional: `go install github.com/air-verse/air@latest` (else we use `go run`).
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

export MOANA_DB_PATH="${MOANA_DB_PATH:-data/moana.db}"
export MOANA_ENV="${MOANA_ENV:-development}"
export MOANA_LISTEN="${MOANA_LISTEN:-:8080}"

DEV_EMAIL="${MOANA_DEV_EMAIL:-admin@moana.local}"
DEV_PASSWORD="${MOANA_DEV_PASSWORD:-changeme}"

echo "Moana dev — DB: $MOANA_DB_PATH  listen: $MOANA_LISTEN"
echo "Credentials: $DEV_EMAIL / $DEV_PASSWORD (local development only)"

if go run ./cmd/moana user add \
  --email="$DEV_EMAIL" \
  --password="$DEV_PASSWORD" \
  --role=admin; then
  echo "Created dev admin."
else
  go run ./cmd/moana user password \
    --email="$DEV_EMAIL" \
    --password="$DEV_PASSWORD"
  echo "Dev admin already existed; password reset to dev default."
fi

echo ""
echo "Starting hot reload:"
echo "  · Vite  — watch → internal/assets/static (CSS/JS)"
echo "  · Air   — rebuild Go when templates/sources or built assets change"

LISTEN_PORT="${MOANA_LISTEN#:}"
if [[ "$LISTEN_PORT" =~ ^[0-9]+$ ]] && [[ "${LISTEN_PORT}" == "8080" ]]; then
  echo "  · Open  http://127.0.0.1:8090  (Air proxy + browser refresh on restart)"
else
  echo "  (Air proxy in .air.toml assumes app on :8080; set MOANA_LISTEN=:8080 or edit [proxy] app_port.)"
fi
echo ""

# One production build so the first Air compile sees up-to-date embeds; then watch.
bun run build
bun run dev:frontend &
VITE_PID=$!

cleanup() {
  if kill -0 "$VITE_PID" 2>/dev/null; then
    kill "$VITE_PID" 2>/dev/null || true
    wait "$VITE_PID" 2>/dev/null || true
  fi
}
trap cleanup EXIT INT TERM

# No exec — keep this shell alive so EXIT trap stops Vite.
if command -v air >/dev/null 2>&1; then
  air
else
  go run github.com/air-verse/air@v1.65.1
fi
