<h1 align="center">YH-FireWall</h1>

<p align="center">
  <img alt="Go Version" src="https://img.shields.io/github/go-mod/go-version/gralliry/YH-FireWall?style=flat-square">
  <img alt="License" src="https://img.shields.io/github/license/gralliry/YH-FireWall?style=flat-square">
</p>

<p align="center">
  A lightweight stateful firewall written in Go, using <b>iptables</b> + <b>NFQUEUE</b> for packet filtering.
</p>

---

## Overview

YH-FireWall is a Linux firewall daemon that intercepts packets at the netfilter level and applies allow/deny rules. It provides both a command-line client and a web management interface.

It is **not** a full router firewall — it does not handle NAT, port forwarding, or DMZ. It filters packets based on IP, port, protocol, and network interface.

Two programs work together:

- **yfwd** — the firewall daemon, runs as root, does the actual packet filtering
- **yfw** — CLI client, sends commands to yfwd over a Unix socket (`/tmp/yfw.sock`)

---

## Features

- **iptables + NFQUEUE** — intercepts traffic at the kernel netfilter level
- **Stateful connection tracking** — TCP/UDP connections tracked with automatic expiration
- **Rule engine** — priority-based matching: rules evaluated in order, first match wins
- **Allow/deny policy** — each rule has an `accept` field (`true` = allow, `false` = deny)
- **Default policy** — when no rule matches, the configurable `default_accept` setting applies
- **Dynamic control** — add, remove, or modify rules at runtime without restart
- **Partial updates** — only send the fields you want to change
- **Web management panel** — Vue 3 interface available at `http://<server-ip>:8080`
- **Swagger API** — OpenAPI documentation at `http://<server-ip>:8080/docs`

---

## Install

Download the latest release tarball from the [releases page](https://github.com/gralliry/YH-FireWall/releases), then:

```bash
tar xzf yfw-linux-*.tar.gz
cd yfw-linux-*
sudo ./install.sh
```

### Requirements

- Linux with `iptables` and `libnetfilter_queue` kernel module
- Root privileges

---

## Build

*For developers.*

```bash
./build.sh
```

Output: `build/yfwd` (daemon) and `build/yfw` (CLI client).

---

## Quick Start

```bash
# 1. Add a rule
./yfw rule append '{"accept":true,"protocols":"tcp","dstPorts":"80,443"}'

# 2. View all rules
./yfw rule list

# 3. Open the web interface
# http://<your-server-ip>:8080
```

---

## Usage

### CLI commands

```bash
# List all rules, or show a specific rule by ID
yfw rule list <id>
yfw rule list abc123

# Add a new rule in JSON format
yfw rule append <json>
yfw rule append '{"accept":true,"dstPorts":"80,443","protocols":"tcp"}'

# Modify an existing rule (partial update, only send fields to change)
yfw rule change <id> <json>
yfw rule change abc123 '{"comment":"updated","priority":100}'

# Set a single field by key and value
yfw rule set <id> <key> <value>
yfw rule set abc123 accept false

# Remove a rule
yfw rule remove <id>
yfw rule remove abc123

# View daemon configuration
yfw config

# List network interfaces (useful for inDevs/outDevs fields)
yfw interfaces

# List supported protocols (useful for protocols field)
yfw protocols

# Show version
yfw version
```

### Service Management

After running `install.sh`, the daemon is managed by systemd:

```bash
systemctl status yfwd       # status
systemctl start yfwd        # start
systemctl stop yfwd         # stop
systemctl restart yfwd      # restart
journalctl -u yfwd -f       # follow logs
```

### Rule JSON fields

| Field | Type | Description |
|-------|------|-------------|
| `group` | string | Rule group label for organization |
| `comment` | string | Human-readable description |
| `accept` | bool | `true` = allow, `false` = deny |
| `priority` | int | Higher number = higher precedence, checked first |
| `enable` | bool | If `false`, the rule is skipped during matching |
| `srcNets` | string | Source IP or CIDR, comma-separated (e.g. `10.0.0.0/8,192.168.1.0/24`) |
| `dstNets` | string | Destination IP or CIDR, comma-separated |
| `srcPorts` | string | Source port, single port or range, comma-separated (e.g. `1024-65535`) |
| `dstPorts` | string | Destination port, single port or range, comma-separated (e.g. `80,443,3000-4000`) |
| `inDevs` | string | Ingress network interface name(s), comma-separated (e.g. `eth0,eth1`) |
| `outDevs` | string | Egress network interface name(s), comma-separated |
| `protocols` | string | IP protocol name(s), comma-separated (e.g. `tcp,udp,icmp`) |

All fields are optional in `append` and `change` — only send what you need. Separators can be comma, space, semicolon, or newline.

### Rule matching

1. Rules are sorted by: enabled first, then by priority (higher number first)
2. The first matching rule wins — if a rule matches, later rules are not evaluated
3. If no rule matches, the **default policy** (configurable in `config.toml`) is applied — `default_accept = true` allows by default

---

## Configuration

Use `yfw config` to view the current configuration.

```toml
version = "1.0.0"          # Config file version (reserved)

[queue]
num = 0                    # NFQUEUE queue number
name = "yfw"               # iptables rule comment

[web]
enable = true              # Enable the web management interface
address = ":8080"          # Listen address (all interfaces, port 8080)
auth_username = ""         # Basic auth username (empty = no auth)
auth_password = ""         # Basic auth password
enable_cors = true         # Enable CORS for API access
static_dir = ""            # Frontend static file directory (empty = embedded files)

[cmd]
socket_path = "/tmp/yfw.sock"  # Unix socket path for CLI communication

[rule]
path = "/etc/yfw/rule.json"    # Rule persistence file
default_accept = true          # Default policy when no rule matches
```

## License

[MIT](LICENSE) © Gralliry
