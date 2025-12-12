#!/bin/sh
set -eu

APP="backend"
OUTDIR="dist/release"
PKG="./cmd/gin-example/main.go"
VERSION="0.1.0"
mkdir -p "$OUTDIR"

# GOOS GOARCH
for t in \
  "linux amd64" \
  "linux arm64" \
  "linux arm" \
  "darwin amd64" \
  "darwin arm64" \
  "windows amd64"; do
  GOOS="$(echo "$t" | awk '{print $1}')"
  GOARCH="$(echo "$t" | awk '{print $2}')"

  EXT=""
  [ "$GOOS" = "windows" ] && EXT=".exe"
  [ "$GOOS" = "linux"   ] && EXT=".bin"
  [ "$GOOS" = "darwin"  ] && EXT=".bin"

  echo "Building Release $GOOS/$GOARCH..."
  env CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" \
    go build -trimpath -ldflags="-s -w" -o "$OUTDIR/${APP}_${GOOS}_${GOARCH}_release_${VERSION}${EXT}" "$PKG"
done

echo "Done. Outputs in ./$OUTDIR"
echo " "

OUTDIR="dist/debug"
PKG="./cmd/gin-example/main.go"

mkdir -p "$OUTDIR"

# GOOS GOARCH
for t in \
  "linux amd64" \
  "linux arm64" \
  "linux arm" \
  "darwin amd64" \
  "darwin arm64" \
  "windows amd64"; do
  GOOS="$(echo "$t" | awk '{print $1}')"
  GOARCH="$(echo "$t" | awk '{print $2}')"

  EXT=""
  [ "$GOOS" = "windows" ] && EXT=".exe"
  [ "$GOOS" = "linux"   ] && EXT=".bin"
  [ "$GOOS" = "darwin"  ] && EXT=".bin"

  echo "Building Debug $GOOS/$GOARCH..."
  env CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" \
    go build -trimpath -o "$OUTDIR/${APP}_${GOOS}_${GOARCH}_debug_${VERSION}${EXT}" "$PKG"
done

echo "Done. Outputs in ./$OUTDIR"
