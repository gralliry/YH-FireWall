#!/bin/bash

# 判断是否 root 用户
if [ "$EUID" -ne 0 ]; then
    echo "Error: This script must be run as root. Please run with sudo."
    exit 1
fi

# 安装程序
install -m 755 yfw-core /usr/local/bin/yfw-core
install -m 755 yfw-client /usr/local/bin/yfw

# 创建 systemd service 文件
SERVICE_FILE="/etc/systemd/system/yfw.service"

echo "Creating systemd service file at $SERVICE_FILE..."

sudo tee "$SERVICE_FILE" > /dev/null <<EOF
[Unit]
Description=YH Firewall Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/yfw-core
Restart=always
User=root
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF

# 重新加载 systemd 配置
echo "Reloading systemd daemon..."
sudo systemctl daemon-reload

# 启动服务
echo "Starting $SERVICE_NAME service..."
sudo systemctl start yfw

# 设置开机自启
echo "Enabling $SERVICE_NAME to start on boot..."
sudo systemctl enable yfw

# 显示状态
echo "Service status:"
sudo systemctl status yfw --no-pager

echo "Installation complete!"