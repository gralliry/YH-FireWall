#!/bin/bash
set -e

VERSION='1.0.0'
BUILD_DIR="build"

# 获取当前脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
INSTALL_SCRIPT="$SCRIPT_DIR/install.sh"

# 支持的架构列表
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

# 清理旧的构建目录
rm -rf "$BUILD_DIR"
mkdir -p "$BUILD_DIR"

# 循环每个架构
for OSARCH in "${ARCH_LIST[@]}"; do
    IFS="/" read -r GOOS GOARCH <<< "$OSARCH"
    echo "Building for $GOOS/$GOARCH..."

    # 创建输出目录
    OUTPUT_DIR="$BUILD_DIR/$GOOS-$GOARCH"
    mkdir -p "$OUTPUT_DIR"

    # 构建 core（纯 Go，禁用 cgo）
    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -o "$OUTPUT_DIR/yfw-core" ./cmd/core

    # 构建 client（纯 Go，禁用 cgo）
    CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -o "$OUTPUT_DIR/yfw-client" ./cmd/client

    # 复制 install.sh
    cp "$INSTALL_SCRIPT" "$OUTPUT_DIR/"

    # 打包为 tar.gz
    TAR_FILE="$BUILD_DIR/yfw-$GOOS-$GOARCH-$VERSION.tar.gz"
    tar -czvf "$TAR_FILE" -C "$OUTPUT_DIR" .

    # 生成 sha256 校验值
    sha256sum "$TAR_FILE" > "$TAR_FILE.sha256"
    echo "Generated hash: $TAR_FILE.sha256"

    # 删除临时目录
    rm -rf "$OUTPUT_DIR"
done

echo "All builds completed."
