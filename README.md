<h1 style="text-align:center;">YH-Firewall</h1>

**YH-Firewall** is a super-lightweight firewall written in Go, leveraging `iptables` and `NFQUEUE` to provide flexible
packet filtering and management.

---

## Features

- Lightweight and high-performance, written in pure Go
- Works with `iptables` to intercept packets via NFQUEUE
- Support for custom rules, groups, and priority-based filtering
- Easy to enable/disable rules and groups dynamically
- Web-based control panel (optional)

---

## Build or Install

```bash
# go install github.com/mitchellh/gox@latest
version='v1.0.0'
# go tool dist list
gox -osarch="linux/386 linux/amd64 linux/arm linux/arm64 linux/loong64 linux/mips linux/mips64 linux/mips64le linux/mipsle linux/ppc64 linux/ppc64le linux/riscv64 linux/s390x" -output="build/yfw-{{.OS}}-{{.Arch}}-$version" ./cmd
```

Make sure you have Go installed (`>=1.20`) and `iptables` available.

And `conntrack` installed.
And `tcpkill` installed. sudo yum install dsniff -y

```bash
apt install dsniff


# Clone the repository
git clone https://github.com/gralliry/YH-Firewall.git

cd YH-Firewall

# Build the project
# Must be executed under root user
go build -o /usr/local/bin/yfw . && chmod +x /usr/local/bin/yfw && yfw
```

## Usage

### Core
Run the firewall service:

```bash
# Start the core service
yfw start

# Stop the service
yfw stop

# Check status
yfw status
```

### Web Interface

```bash
# Start web interface: yfw web start <address> <username> <password>
yfw web start 0.0.0.0:8080 admin admin123

# Stop web interface
yfw web stop

# Check web interface status
yfw web status
```

### Rule Management

```bash
# List all rules
yfw rule list
# Get a specific rule
yfw rule list <rule_id>

# Add a new rule
yfw rule add '{"srcNet":"0.0.0.0/0","tarPort":"80"}'

# Remove a rule
yfw rule remove <rule_id>

# Change a rule
yfw rule change <rule_id> '{"srcNet":"192.168.1.0/24"}'

# Enable/disable a rule
yfw rule enable <rule_id>
yfw rule disable <rule_id>
```

### Group Management

```bash
# Enable/disable a group
yfw group enable <group_name>
yfw group disable <group_name>
```

### Other

```bash
# Show all rules
sudo iptables -L -n -v
# 
QUEUE_NUM=1
# Append NFQUEUE to INPUT/OUTPUT/FORWARD
sudo iptables -I INPUT   -j NFQUEUE --queue-num "$QUEUE_NUM" -m comment --comment "yfw"
sudo iptables -I OUTPUT  -j NFQUEUE --queue-num "$QUEUE_NUM" -m comment --comment "yfw"
sudo iptables -I FORWARD -j NFQUEUE --queue-num "$QUEUE_NUM" -m comment --comment "yfw"
# Delete all rules with yfw comment
sudo iptables -D INPUT   -j NFQUEUE --queue-num "$QUEUE_NUM" -m comment --comment "yfw"
sudo iptables -D OUTPUT  -j NFQUEUE --queue-num "$QUEUE_NUM" -m comment --comment "yfw"
sudo iptables -D FORWARD -j NFQUEUE --queue-num "$QUEUE_NUM" -m comment --comment "yfw"
```

## License

MIT License © Gralliry

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for bug fixes or new features.

## Notes

* Requires __root__ privileges to manipulate iptables rules.

* Works best on Linux environments with iptables and libnetfilter_queue installed.

## Other

```bash
sudo iptables -L -v -n
```