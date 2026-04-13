.PHONY: build-frontend test test-e2e build dev

# Vite bundles CSS/JS into internal/handlers/static (required for go:embed).
build-frontend:
	bun install && bun run build

# Binary output: bin/moana (gitignored).
build:
	@mkdir -p bin
	go build -o bin/moana ./cmd/moana

test: build-frontend
	go test -race ./...

test-e2e: build-frontend
	bun run test:e2e

# Local dev: Vite build, ensure admin@moana.local / changeme, then go run serve.
dev: build-frontend
	bash scripts/dev.sh
