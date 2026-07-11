#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"
OUT="$ROOT/bin"
PKG="./cmd/vibe-shield"
LDFLAGS='-s -w'

mkdir -p "$OUT"

build() {
  local goos=$1 goarch=$2 name=$3
  echo "==> $goos/$goarch -> bin/$name"
  CGO_ENABLED=0 GOOS=$goos GOARCH=$goarch \
    go build -ldflags="$LDFLAGS" -o "$OUT/$name" "$PKG"
}

build darwin  amd64 vibe-shield-darwin-amd64
build darwin  arm64 vibe-shield-darwin-arm64
build linux   amd64 vibe-shield-linux-amd64
build windows amd64 vibe-shield-windows-amd64.exe

echo "Done. Binaries in bin/"
