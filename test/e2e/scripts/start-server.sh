#!/usr/bin/env bash
set -euo pipefail
ROOT="$(cd "$(dirname "$0")/../../.." && pwd)"
cd "$ROOT"
mkdir -p "$ROOT/data"
rm -f "$ROOT/data/e2e-test.db"
export MOANA_DB_PATH="$ROOT/data/e2e-test.db"
export MOANA_LISTEN=":18080"
export MOANA_ENV="development"
export MOANA_SESSION_SECRET="e2e-test-session-secret-32bytes!!"
go run ./cmd/moana user add --email=e2e@moana.test --password=password123 --role=user
exec go run ./cmd/moana serve
