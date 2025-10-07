# Y-Firewall

a super-lighning firewall based on iptables by go

```shell
sudo iptables -I INPUT -j NFQUEUE --queue-num 1
sudo iptables -I OUTPUT -j NFQUEUE --queue-num 1
sudo iptables -I FORWARD -j NFQUEUE --queue-num 1
```

```shell
sudo iptables -t raw -F
sudo iptables -t mangle -F
sudo iptables -t nat -F
sudo iptables -t filter -F

sudo iptables -t raw -X
sudo iptables -t mangle -X
sudo iptables -t nat -X
sudo iptables -t filter -X

sudo iptables -P INPUT ACCEPT
sudo iptables -P OUTPUT ACCEPT
sudo iptables -P FORWARD ACCEPT
```

```shell
# raw 表：最早阶段
sudo iptables -t raw -A PREROUTING  -j NFQUEUE --queue-num 0 -m comment --comment "yfw"
sudo iptables -t raw -A OUTPUT      -j NFQUEUE --queue-num 0 -m comment --comment "yfw"

# mangle 表：处理前后所有阶段
sudo iptables -t mangle -A PREROUTING  -j NFQUEUE --queue-num 0 -m comment --comment "yfw"
sudo iptables -t mangle -A INPUT       -j NFQUEUE --queue-num 0 -m comment --comment "yfw"
sudo iptables -t mangle -A FORWARD     -j NFQUEUE --queue-num 0 -m comment --comment "yfw"
sudo iptables -t mangle -A OUTPUT      -j NFQUEUE --queue-num 0 -m comment --comment "yfw"
sudo iptables -t mangle -A POSTROUTING -j NFQUEUE --queue-num 0 -m comment --comment "yfw"

# filter 表（常用表）
sudo iptables -t filter -A INPUT   -j NFQUEUE --queue-num 0 -m comment --comment "yfw"
sudo iptables -t filter -A OUTPUT  -j NFQUEUE --queue-num 0 -m comment --comment "yfw"
sudo iptables -t filter -A FORWARD -j NFQUEUE --queue-num 0 -m comment --comment "yfw"

# nat 表（做地址转换前后）
sudo iptables -t nat -A PREROUTING  -j NFQUEUE --queue-num 0 -m comment --comment "yfw"
sudo iptables -t nat -A OUTPUT      -j NFQUEUE --queue-num 0 -m comment --comment "yfw"
sudo iptables -t nat -A POSTROUTING -j NFQUEUE --queue-num 0 -m comment --comment "yfw"
```

```shell
# 遍历所有表和链，删除 comment 为 yfw 的规则
for table in raw mangle filter nat; do
  for chain in $(sudo iptables -t $table -S | grep '^-A' | awk '{print $2}'); do
    while sudo iptables -t $table -C $chain -m comment --comment "yfw" &>/dev/null; do
      sudo iptables -t $table -D $chain -m comment --comment "yfw"
    done
  done
done
```

```shell
sudo iptables -L -n -v
```