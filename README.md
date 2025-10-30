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

## Install

Make sure you have `iptables` available.

```bash
# Clone the repository
git clone https://github.com/gralliry/YH-Firewall.git

cd YH-Firewall

# Build the project
# Must be executed under root user
go build -o /usr/local/bin/yfw . && chmod +x /usr/local/bin/yfw && yfw
```

## Usage

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

## Compile

Make sure you have Go installed (`>=1.24`) available.

```bash
go install github.com/mitchellh/gox@latest
version='v1.0.0'

# go tool dist list
gox -osarch="linux/386 linux/amd64 linux/arm linux/arm64 linux/loong64 linux/mips linux/mips64 linux/mips64le linux/mipsle linux/ppc64 linux/ppc64le linux/riscv64 linux/s390x" -output="build/yfw-{{.OS}}-{{.Arch}}-$version" ./cmd
```

## License

MIT License © Gralliry

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for bug fixes or new features.

## Notes

* Requires __root__ privileges to manipulate iptables rules.

* Works best on Linux environments with iptables and libnetfilter_queue installed.