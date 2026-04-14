.PHONY: build-frontend test test-e2e build dev

# Vite bundles CSS/JS into internal/assets/static (required for go:embed).
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

# Hot reload: Vite watch (CSS/JS → internal/assets/static) + Air (Go rebuild & serve).
# Use http://127.0.0.1:8090 for Air’s proxy + browser live reload (app listens on MOANA_LISTEN, default :8080).
dev:
	bun install
	bash scripts/dev.sh
