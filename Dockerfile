# syntax=docker/dockerfile:1
#
# BuildKit: uses cache mounts (repeat builds are much faster on CI/VPS). Requires
# BuildKit (default in recent Docker: `docker build` / `docker compose build`).
# First-time `apk` / image pulls and `bun run build` can take several minutes on
# a small host or slow link to Docker Hub / apk mirrors — that is normal, not a hang.

# Frontend (Vite + Bun) → internal/assets/static
FROM oven/bun:1.3.13 AS assets
WORKDIR /app
ENV NODE_ENV=production
COPY package.json bun.lock ./
COPY frontend ./frontend
COPY internal/assets/static ./internal/assets/static
# Bun download cache: speeds rebuilds; safe to change path with Bun major upgrades
RUN --mount=type=cache,target=/root/.bun/install/cache \
  bun install --frozen-lockfile && bun run build \
  && test -f internal/assets/static/favicon-mo.svg \
  && test -f internal/assets/static/css/app.css \
  && test -f internal/assets/static/js/app.js

# Go binary
FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
  go mod download
COPY . .
COPY --from=assets /app/internal/assets/static ./internal/assets/static
RUN test -f internal/assets/static/favicon-mo.svg \
  && test -f internal/assets/static/css/app.css \
  && test -f internal/assets/static/js/app.js
# module cache: speeds rebuild; public modules need no `git` in the image
RUN --mount=type=cache,target=/go/pkg/mod \
  CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/moana ./cmd/moana

# Run — use a current Alpine (3.23+) for fresh ca-cert / tz / base image security fixes; the Go binary is static.
FROM alpine:3.23
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /data
ENV MOANA_DB_PATH=/data/moana.db
ENV MOANA_LISTEN=:8080
ENV MOANA_ENV=production
EXPOSE 8080
VOLUME ["/data"]
COPY --from=build /out/moana /bin/moana
CMD ["/bin/moana"]
