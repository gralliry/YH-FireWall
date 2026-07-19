<h1 align="center">YH-FireWall</h1>

<p align="center">
  <img alt="Go Version" src="https://img.shields.io/github/go-mod/go-version/gralliry/YH-FireWall?style=flat-square">
  <img alt="License" src="https://img.shields.io/github/license/gralliry/YH-FireWall?style=flat-square">
  <img alt="Build" src="https://img.shields.io/github/actions/workflow/status/gralliry/YH-FireWall/build.yml?style=flat-square">
</p>

<p align="center">
  A lightweight firewall written in Go, using <b>iptables</b> + <b>NFQUEUE</b> for packet filtering.
</p>

---

## Features

- **Pure Go** — single binary, no runtime dependencies
- **iptables + NFQUEUE** — intercepts packets at the netfilter level
- **Rule engine** — priority-based rule matching with IP, port, protocol, and interface filtering
- **Connection tracking** — stateful tracking of TCP/UDP connections
- **Dynamic control** — enable, disable, add, or remove rules at runtime without restart
- **Command-line client** — manage everything from the terminal
- **Web panel** — Vue-based web management interface
- **Swagger API** — OpenAPI documentation at `/docs`

---

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/gralliry/YH-FireWall/master/install.sh | sudo bash
```

Or pin a specific version, or change the config file path:

```bash
# Latest version, custom config path
curl -fsSL https://raw.githubusercontent.com/gralliry/YH-FireWall/master/install.sh | sudo bash -s -- -c /custom/config.toml

# Specific version, default config path
curl -fsSL https://raw.githubusercontent.com/gralliry/YH-FireWall/master/install.sh | sudo bash -s -- -v v1.0.0

# Both
curl -fsSL https://raw.githubusercontent.com/gralliry/YH-FireWall/master/install.sh | sudo bash -s -- -v v1.0.0 -c /custom/config.toml
```

The installer auto-detects your architecture, downloads the latest release, verifies the checksum, installs the binary, and sets up a systemd service.

Two separate binaries:
- `yfwd` — the firewall daemon (runs as a service)
- `yfw` — CLI client to manage the daemon

### Requirements

- Linux with `iptables` and `libnetfilter_queue`
- Root privileges

---

## Build

```bash
# Install swagger CLI (first time only)
go install github.com/swaggo/swag/cmd/swag@latest

# Build everything
./build.sh
```

`build.sh` runs: swagger generation → frontend build (npm) → Go backend compilation. Output goes to `build/`.

---

## Usage

```bash
# Start the daemon
sudo yfwd -c /etc/yfw/config.toml

# List all rules
yfw rule list

# Get a specific rule
yfw rule list <rule_id>

# Add a rule
yfw rule add '{"srcPrefixs":"192.168.1.0/24","dstPortRanges":"80","accept":true}'

# Modify a rule (partial update)
yfw rule change <rule_id> '{"comment":"updated"}'

# Enable / disable a rule
yfw rule enable <rule_id>
yfw rule disable <rule_id>

# Remove a rule
yfw rule remove <rule_id>

# View config
yfw config

# Version
yfw version
```

### Rule fields

Rules use codec strings for IP/port/protocol/device fields:

| Field | Type | Example |
|-------|------|---------|
| `group` | `string` | `"web"` |
| `comment` | `string` | `"Allow HTTP"` |
| `accept` | `bool` | `true` / `false` |
| `priority` | `int` | `100` |
| `enable` | `bool` | `true` / `false` |
| `srcPrefixs` | `string` | `"192.168.1.0/24,10.0.0.0/8"` |
| `dstPrefixs` | `string` | `"0.0.0.0/0"` |
| `srcPortRanges` | `string` | `"1024-65535"` |
| `dstPortRanges` | `string` | `"80,443"` |
| `inDevs` | `string` | `"eth0,eth1"` |
| `outDevs` | `string` | `"eth0"` |
| `protocols` | `string` | `"tcp,udp"` |

All codec fields accept comma/space/semicolon/newline as separators.

### Web API

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/ping` | Health check |
| `GET` | `/api/rule` | List rules |
| `POST` | `/api/rule` | Add rule |
| `PUT` | `/api/rule/{id}` | Update rule |
| `DELETE` | `/api/rule/{id}` | Delete rule |
| `GET` | `/api/config` | Get config |
| `POST` | `/api/config` | Update config |
| `GET` | `/api/connection` | List active connections |
| `DELETE` | `/api/connection/{id}` | Close a connection |
| `GET` | `/api/interface` | List network interfaces |
| `GET` | `/docs` | Swagger UI |

---

## Configuration

The config file lives at `/etc/yfw/config.toml`:

```toml
version = "1.0.0"

[queue]
no = 0
name = "yfw"

[web]
enable = true
address = ":8080"
auth_username = "admin"
auth_password = "admin"
enable_cors = true

[cmd]
socket_path = "/tmp/yfw.sock"

[rule]
path = "/etc/yfw/rule.json"
default_accept = true
```

---

## Service Management

```bash
# Managed by systemd (after install)
systemctl status yfwd     # check status
systemctl start yfwd      # start
systemctl stop yfwd       # stop
systemctl restart yfwd    # restart
journalctl -u yfwd -f     # follow logs
```

---

## License

[MIT](LICENSE) © Gralliry
