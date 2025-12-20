<h1 style="text-align:center;">YH-Firewall</h1>

**YH-Firewall** is a super-lightweight firewall written in Go, leveraging `iptables` and `NFQUEUE` to provide flexible
packet filtering and management.

---

## Features

- Lightweight and high-performance, written in pure Go
- Works with `iptables` to intercept packets via NFQUEUE (Make sure you have `iptables` available)
- Support for custom rules, groups, and priority-based filtering
- Easy to enable/disable rules and groups dynamically
- Web-based control panel (optional)

---

## Install

```bash
tar -xzvf yfw-XXX-XXX-XXX.tar.gz -C /tmp/yfw
cd /tmp/yfw
bash install.sh
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

## Config

In```/etc/yfw/config.yaml```, for example:

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

## License

MIT License © Gralliry

## Contributing

Contributions are welcome! Please open an issue or submit a pull request for bug fixes or new features.

## Notes

* Requires __root__ privileges to manipulate iptables rules.

* Works best on Linux environments with iptables and libnetfilter_queue installed.

## Author's words
It is just a homework about computer network application in my college, so it is just a toy project. Don't be too serious.
And ... DO NOT USE IT IN PRODUCTION ENVIRONMENT!!!
I need to submit my homework to my teacher, (s)he will check it, so I must use chinese comments.
If someone thinks it is a good project or star it, i will try to translate it into English and keep it updated.
