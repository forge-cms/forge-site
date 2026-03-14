# syntax=docker/dockerfile:1

# ── build stage ──────────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS build

WORKDIR /src

# Download dependencies first (layer-cached until go.mod/go.sum change).
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build.
# CGO_ENABLED=0: modernc.org/sqlite is pure-Go — no C toolchain needed.
# -trimpath: strip local file paths from the binary.
# -ldflags: embed VERSION and strip debug symbols for a smaller image.
COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 go build \
      -trimpath \
      -ldflags "-s -w -X main.Version=${VERSION}" \
      -o /forge-site \
      .

# ── runtime stage ─────────────────────────────────────────────────────────────
FROM alpine:latest

# ca-certificates: needed for outbound TLS (e.g. ACME, external APIs).
RUN apk add --no-cache ca-certificates

# Run as a non-root user.
RUN adduser -D -u 1000 app
USER app

WORKDIR /app

COPY --from=build /forge-site /app/forge-site

# Data directory is provided by a Docker volume at runtime.
VOLUME ["/app/data"]

EXPOSE 8080

ENV PORT=8080

ENTRYPOINT ["/app/forge-site"]
