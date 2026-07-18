#!/bin/bash
set -e

REPO="gralliry/YH-FireWall"
VERSION="latest"
CONFIG_PATH="/etc/yfw/config.toml"

usage() {
    echo "Usage: curl -fsSL https://raw.githubusercontent.com/$REPO/main/install.sh | sudo bash"
    echo "       curl -fsSL ... | sudo bash -s -- -v v1.0.0 -c /path/to/config.toml"
    echo ""
    echo "Options:"
    echo "  -v  Version to install (default: latest)"
    echo "  -c  Config file path  (default: /etc/yfw/config.toml)"
    exit 1
}

while getopts "v:c:h" opt; do
    case "$opt" in
        v) VERSION="$OPTARG" ;;
        c) CONFIG_PATH="$OPTARG" ;;
        h) usage ;;
        *) usage ;;
    esac
done

TMP_DIR="$(mktemp -d)"
cleanup() { rm -rf "$TMP_DIR"; }
trap cleanup EXIT

# ---------- root check ----------
if [ "$EUID" -ne 0 ]; then
    echo "Error: This script must be run as root."
    exit 1
fi

# ---------- arch detection ----------
detect_arch() {
    local machine
    machine="$(uname -m)"
    case "$machine" in
        x86_64)    echo "amd64" ;;
        i386|i686) echo "386" ;;
        armv6l|armv7l) echo "arm" ;;
        aarch64)   echo "arm64" ;;
        loongarch64) echo "loong64" ;;
        mips)      echo "mips" ;;
        mips64)    echo "mips64" ;;
        mips64el)  echo "mips64le" ;;
        mipsel)    echo "mipsle" ;;
        ppc64)     echo "ppc64" ;;
        ppc64le)   echo "ppc64le" ;;
        riscv64)   echo "riscv64" ;;
        s390x)     echo "s390x" ;;
        *)         echo "" ;;
    esac
}

ARCH="$(detect_arch)"
if [ -z "$ARCH" ]; then
    echo "Error: Unsupported architecture: $(uname -m)"
    exit 1
fi
echo "Detected architecture: $ARCH"

# ---------- resolve version ----------
if [ "$VERSION" = "latest" ]; then
    echo "Resolving latest version..."
    VERSION="$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
        | grep '"tag_name":' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')"
    if [ -z "$VERSION" ]; then
        echo "Error: Failed to resolve latest version."
        exit 1
    fi
fi
echo "Installing version: $VERSION"

BASE_URL="https://github.com/$REPO/releases/download/$VERSION"

# ---------- download ----------
download_verify() {
    local name="$1" out="$TMP_DIR/$name"
    echo "Downloading $name ..."
    curl -fsSL -o "$out" "$BASE_URL/$name"
    curl -fsSL -o "$out.sha256" "$BASE_URL/$name.sha256"
    expected="$(cut -d' ' -f1 "$out.sha256")"
    actual="$(sha256sum "$out" | cut -d' ' -f1)"
    if [ "$expected" != "$actual" ]; then
        echo "Error: Checksum mismatch for $name."
        echo "Expected: $expected"
        echo "Got:      $actual"
        exit 1
    fi
    echo "Checksum OK."
}

download_verify "yfwd-linux-$ARCH"
download_verify "yfw-linux-$ARCH"

# ---------- install ----------
install -m 755 "$TMP_DIR/yfwd-linux-$ARCH" /usr/local/bin/yfwd
install -m 755 "$TMP_DIR/yfw-linux-$ARCH" /usr/local/bin/yfw

# ---------- systemd ----------
SERVICE_FILE="/etc/systemd/system/yfwd.service"
echo "Creating systemd service: $SERVICE_FILE"
tee "$SERVICE_FILE" > /dev/null <<EOF
[Unit]
Description=YH Firewall Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/yfwd -c $CONFIG_PATH
Restart=always
User=root
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable yfwd
systemctl start yfwd

echo ""
echo "============================================"
echo "  YH-FireWall $VERSION installed successfully."
echo "  arch: $ARCH"
echo "  config:  $CONFIG_PATH"
echo "  service: systemctl status yfwd"
echo "  client:  yfw help"
echo "============================================"
