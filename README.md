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
- **Rule engine** — custom rules with groups and priority-based filtering
- **Dynamic control** — enable, disable, add, or remove rules at runtime
- **Command-line client** — manage everything from the terminal
- **Web panel** — optional web-based management interface

---

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/gralliry/YH-FireWall/master/scripts/install.sh | sudo bash
```

Or pin a specific version:

```bash
curl -fsSL https://raw.githubusercontent.com/gralliry/YH-FireWall/master/scripts/install.sh | sudo bash -s v1.0.0
```

The installer auto-detects your architecture, downloads the latest release, verifies the checksum, installs binaries to `/usr/local/bin`, and sets up a systemd service.

### Requirements

- Linux with `iptables` and `libnetfilter_queue`
- Root privileges

---

## Usage

```bash
# List all rules
yfw rule list

# Get a specific rule
yfw rule list <rule_id>

# Add a rule
yfw rule add '{"srcNet":"0.0.0.0/0","tarPort":"80"}'

# Modify a rule
yfw rule change <rule_id> '{"srcNet":"192.168.1.0/24"}'

# Enable / disable a rule
yfw rule enable <rule_id>
yfw rule disable <rule_id>

# Remove a rule
yfw rule remove <rule_id>
```

---

## Configuration

The config file lives at `/etc/yfw/config.yaml`:

```yaml
queue_no: 0

web:
  enable: true
  address: 0.0.0.0:8080
  auth_username: admin
  auth_password: admin
  enable_cors: true

cmd:
  enable: true
  socket_path: /tmp/yfw.sock

rule_table:
  path: /etc/yfw/rule.json
```

---

## Service Management

```bash
systemctl status yfw     # check status
systemctl start yfw      # start
systemctl stop yfw       # stop
systemctl restart yfw    # restart
journalctl -u yfw -f     # follow logs
```

---

## License

[MIT](LICENSE) © Gralliry

---

> [!WARNING]  
> This is a college course project and **not production-ready**. Use at your own risk.