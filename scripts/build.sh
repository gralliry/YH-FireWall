#!/bin/bash
set -e

BUILD_DIR="build"
mkdir -p "$BUILD_DIR"

ARCH_LIST=(
    "linux/386"
    "linux/amd64"
    "linux/arm"
    "linux/arm64"
    "linux/loong64"
    "linux/mips"
    "linux/mips64"
    "linux/mips64le"
    "linux/mipsle"
    "linux/ppc64"
    "linux/ppc64le"
    "linux/riscv64"
    "linux/s390x"
)

for OSARCH in "${ARCH_LIST[@]}"; do
    IFS="/" read -r GOOS GOARCH <<< "$OSARCH"
    echo "Building for $GOOS/$GOARCH..."

    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -o "$BUILD_DIR/yfw-core-$GOOS-$GOARCH"   ./cmd/core
    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -o "$BUILD_DIR/yfw-client-$GOOS-$GOARCH" ./cmd/client

    sha256sum "$BUILD_DIR/yfw-core-$GOOS-$GOARCH"   > "$BUILD_DIR/yfw-core-$GOOS-$GOARCH.sha256"
    sha256sum "$BUILD_DIR/yfw-client-$GOOS-$GOARCH" > "$BUILD_DIR/yfw-client-$GOOS-$GOARCH.sha256"
done

echo "All builds completed."
