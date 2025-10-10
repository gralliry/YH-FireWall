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

## Build

```bash
sudo go build -o yfw ./cmd 

sudo cp yfw /usr/local/bin/

sudo chmod +x /usr/local/bin/yfw
```

```bash
go install github.com/mitchellh/gox@latest

version='v1.0.0'

go tool dist list

gox -osarch="linux/386 linux/amd64 linux/arm linux/arm64 linux/loong64 linux/mips linux/mips64 linux/mips64le linux/mipsle linux/ppc64 linux/ppc64le linux/riscv64 linux/s390x" -output="build/yfw-{{.OS}}-{{.Arch}}-$version" ./cmd
```

## Install

Make sure you have Go installed (`>=1.20`) and `iptables` available.

```bash
# Clone the repository
git clone https://github.com/gralliry/YH-Firewall.git

cd YH-Firewall

# Build the project
go build -o yfw ./cmd
```

## Usage

Run the firewall service:

```bash
# Start the core service
yfw start

# Stop the service
yfw stop

# Check status
yfw status
```

## Web Interface

```bash
# Start web interface: yfw web start <address> <username> <password>
yfw web start 0.0.0.0:8080 admin admin123

# Stop web interface
yfw web stop

# Check web interface status
yfw web status
```

## Rule Management

```bash
# List all rules
yfw rule list

# Get a specific rule
yfw rule list <rule_id>

# Add a new rule
yfw rule add '{"src_net":"0.0.0.0/0","tar_port":"80"}'

# Remove a rule
yfw rule remove <rule_id>

# Change a rule
yfw rule change <rule_id> '{"src_net":"192.168.1.0/24"}'

# Enable/disable a rule
yfw rule enable <rule_id>
yfw rule disable <rule_id>
```

## Group Management

```bash
# Enable/disable a group
yfw group enable <group_name>
yfw group disable <group_name>
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