#!/usr/bin/env sh
# Build script run inside the Docker container.
# Usage: sh build/build.sh [windows-only]
set -e

cd /src
go mod tidy

# Resolve version metadata.
# VERSION: nearest git tag (e.g. v1.0.0), or "dev" if no tags exist.
# COMMIT:  short commit hash, or "unknown" if no commits.
# DATE:    build date in YYYY-MM-DD format.
VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "dev")
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE=$(date -u +%Y-%m-%d)

LDFLAGS="-s -w \
  -X github.com/raz0red/radiscordpresence/internal/buildinfo.Version=${VERSION} \
  -X github.com/raz0red/radiscordpresence/internal/buildinfo.Commit=${COMMIT} \
  -X github.com/raz0red/radiscordpresence/internal/buildinfo.Date=${DATE}"

if [ "${1}" = "windows-only" ]; then
    echo "Building Windows amd64 (${VERSION} ${COMMIT})..."
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="${LDFLAGS}" -o dist/radiscordpresence.exe ./cmd/radiscordpresence
    echo "Windows build complete: dist/radiscordpresence.exe"
else
    echo "Building Windows amd64 (${VERSION} ${COMMIT})..."
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="${LDFLAGS}" -o dist/radiscordpresence-windows-amd64.exe ./cmd/radiscordpresence

    echo "Building Linux amd64..."
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="${LDFLAGS}" -o dist/radiscordpresence-linux-amd64 ./cmd/radiscordpresence

    echo "Building macOS amd64..."
    GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="${LDFLAGS}" -o dist/radiscordpresence-darwin-amd64 ./cmd/radiscordpresence

    echo "Building macOS arm64..."
    GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="${LDFLAGS}" -o dist/radiscordpresence-darwin-arm64 ./cmd/radiscordpresence

    echo "All builds complete."
fi
