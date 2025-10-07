#!/bin/bash

ACTION="$1"

QUEUE_NUM="${2:-0}"  # 默认队列号 0

case "$ACTION" in
  set)
    echo "[*] 添加 NFQUEUE 到 INPUT/OUTPUT/FORWARD"
    sudo iptables -I INPUT   -j NFQUEUE --queue-num "$QUEUE_NUM" -m comment --comment "yfw"
    sudo iptables -I OUTPUT  -j NFQUEUE --queue-num "$QUEUE_NUM" -m comment --comment "yfw"
    sudo iptables -I FORWARD -j NFQUEUE --queue-num "$QUEUE_NUM" -m comment --comment "yfw"
    ;;

  clear)
    echo "[*] 清空所有表规则"
    for table in raw mangle nat filter; do
      sudo iptables -t $table -F
      sudo iptables -t $table -X
    done
    sudo iptables -P INPUT ACCEPT
    sudo iptables -P OUTPUT ACCEPT
    sudo iptables -P FORWARD ACCEPT
    ;;

  remove)
    echo "[*] 删除所有带 yfw 注释的规则"
    sudo iptables -D INPUT   -j NFQUEUE --queue-num "$QUEUE_NUM" -m comment --comment "yfw"
    sudo iptables -D OUTPUT  -j NFQUEUE --queue-num "$QUEUE_NUM" -m comment --comment "yfw"
    sudo iptables -D FORWARD -j NFQUEUE --queue-num "$QUEUE_NUM" -m comment --comment "yfw"
    ;;

  list)
    echo "[*] 显示所有规则"
    sudo iptables -L -n -v
    ;;

  *)
    echo "Usage: $0 {set|clear|remove|list} [queue-num]"
    exit 1
    ;;
esac
