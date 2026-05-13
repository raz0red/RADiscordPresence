#!/usr/bin/env sh
# Build script run inside the Docker container.
# Usage: sh build/build.sh [windows-only|mac-only|linux-only]
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
  -X github.com/raz0red/radpresence/internal/buildinfo.Version=${VERSION} \
  -X github.com/raz0red/radpresence/internal/buildinfo.Commit=${COMMIT} \
  -X github.com/raz0red/radpresence/internal/buildinfo.Date=${DATE}"

mkdir -p dist

if [ "${1}" = "windows-only" ]; then
    echo "Building Windows amd64 (${VERSION} ${COMMIT})..."
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="${LDFLAGS}" -o dist/radpresence.exe ./cmd/radiscordpresence
    echo "Build complete: dist/radpresence.exe"
elif [ "${1}" = "mac-only" ]; then
    echo "Building macOS amd64 (${VERSION} ${COMMIT})..."
    GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="${LDFLAGS}" -o dist/radpresence-darwin-amd64 ./cmd/radiscordpresence
    echo "Building macOS arm64..."
    GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="${LDFLAGS}" -o dist/radpresence-darwin-arm64 ./cmd/radiscordpresence
    echo "Build complete: dist/radpresence-darwin-amd64, dist/radpresence-darwin-arm64"
elif [ "${1}" = "linux-only" ]; then
    echo "Building Linux amd64 (${VERSION} ${COMMIT})..."
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="${LDFLAGS}" -o dist/radpresence-linux-amd64 ./cmd/radiscordpresence
    echo "Build complete: dist/radpresence-linux-amd64"
else
    echo "Building Windows amd64 (${VERSION} ${COMMIT})..."
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="${LDFLAGS}" -o dist/radpresence-windows-amd64.exe ./cmd/radiscordpresence

    echo "Building Linux amd64..."
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="${LDFLAGS}" -o dist/radpresence-linux-amd64 ./cmd/radiscordpresence

    echo "Building macOS amd64..."
    GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="${LDFLAGS}" -o dist/radpresence-darwin-amd64 ./cmd/radiscordpresence

    echo "Building macOS arm64..."
    GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -ldflags="${LDFLAGS}" -o dist/radpresence-darwin-arm64 ./cmd/radiscordpresence

    echo "All builds complete."
fi
