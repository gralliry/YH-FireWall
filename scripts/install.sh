#!/bin/bash
set -e

REPO="gralliry/YH-FireWall"
VERSION="${1:-latest}"
CONFIG_PATH="${2:-/etc/yfw/config.yaml}"
TMP_DIR="$(mktemp -d)"

cleanup() { rm -rf "$TMP_DIR"; }
trap cleanup EXIT

# ---------- root check ----------
if [ "$EUID" -ne 0 ]; then
    echo "Error: This script must be run as root."
    echo "Usage: curl -fsSL https://raw.githubusercontent.com/$REPO/master/scripts/install.sh | sudo bash"
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
FILE="yfw-linux-$ARCH"

# ---------- download ----------
echo "Downloading $FILE ..."
curl -fsSL -o "$TMP_DIR/$FILE" "$BASE_URL/$FILE"
curl -fsSL -o "$TMP_DIR/$FILE.sha256" "$BASE_URL/$FILE.sha256"
EXPECTED="$(cut -d' ' -f1 "$TMP_DIR/$FILE.sha256")"
ACTUAL="$(sha256sum "$TMP_DIR/$FILE" | cut -d' ' -f1)"
if [ "$EXPECTED" != "$ACTUAL" ]; then
    echo "Error: Checksum mismatch."
    echo "Expected: $EXPECTED"
    echo "Got:      $ACTUAL"
    exit 1
fi
echo "Checksum OK."

# ---------- install ----------
install -m 755 "$TMP_DIR/$FILE" /usr/local/bin/yfw

# ---------- systemd ----------
SERVICE_FILE="/etc/systemd/system/yfw.service"
echo "Creating systemd service: $SERVICE_FILE"
tee "$SERVICE_FILE" > /dev/null <<EOF
[Unit]
Description=YH Firewall Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/yfw core -c $CONFIG_PATH
Restart=always
User=root
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl enable yfw
systemctl start yfw

echo ""
echo "============================================"
echo "  YH-FireWall $VERSION installed successfully."
echo "  arch: $ARCH"
echo "  service: systemctl status yfw"
echo "  config:  $CONFIG_PATH"
echo "  client:  yfw help"
echo "============================================"
