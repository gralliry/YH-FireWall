#!/bin/bash
set -e

INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/yfw"

# ---------- root check ----------
if [ "$EUID" -ne 0 ]; then
    echo "Error: This script must be run as root."
    exit 1
fi

# ---------- resolve script directory ----------
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

# ---------- install binaries ----------
echo "Installing yfwd..."
install -m 755 "$SCRIPT_DIR/yfwd" "$INSTALL_DIR/yfwd"

echo "Installing yfw..."
install -m 755 "$SCRIPT_DIR/yfw" "$INSTALL_DIR/yfw"

# ---------- install config ----------
mkdir -p "$CONFIG_DIR"

# ---------- systemd service ----------
SERVICE_FILE="/etc/systemd/system/yfwd.service"
echo "Creating systemd service..."
cat > "$SERVICE_FILE" <<EOF
[Unit]
Description=YH FireWall Service
After=network.target

[Service]
Type=simple
ExecStart=$INSTALL_DIR/yfwd -c $CONFIG_DIR/config.toml
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
echo "  YH-FireWall installed successfully."
echo ""
echo "  config:  $CONFIG_DIR/config.toml"
echo "  service: systemctl status yfwd"
echo "  client:  yfw help"
echo "============================================"
