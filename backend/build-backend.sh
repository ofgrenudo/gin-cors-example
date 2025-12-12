#!/bin/sh
set -eu

APP="backend" # output name
OUTDIR="dist" # output folder
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

  echo "Building $GOOS/$GOARCH..."
  env CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" \
    go build -trimpath -ldflags="-s -w" -o "$OUTDIR/${APP}_${GOOS}_${GOARCH}${EXT}" "$PKG"
done

echo "Done. Outputs in ./$OUTDIR"
