# Frontend (Vite + Bun) → internal/assets/static
FROM oven/bun:1 AS assets
WORKDIR /app
COPY package.json bun.lock ./
COPY frontend ./frontend
COPY internal/assets/static ./internal/assets/static
RUN bun install --frozen-lockfile && bun run build

# Go binary
FROM golang:1.26-alpine AS build
RUN apk add --no-cache git
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=assets /app/internal/assets/static ./internal/assets/static
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/moana ./cmd/moana

# Run
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /data
ENV MOANA_DB_PATH=/data/moana.db
ENV MOANA_LISTEN=:8080
ENV MOANA_ENV=production
EXPOSE 8080
VOLUME ["/data"]
COPY --from=build /out/moana /bin/moana
CMD ["/bin/moana"]
