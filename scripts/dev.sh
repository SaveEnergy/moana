#!/usr/bin/env bash
# Development server with a known admin account (local use only).
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
  --role=admin \
  --timezone=UTC; then
  echo "Created dev admin."
else
  go run ./cmd/moana user password \
    --email="$DEV_EMAIL" \
    --password="$DEV_PASSWORD"
  echo "Dev admin already existed; password reset to dev default."
fi

exec go run ./cmd/moana serve
