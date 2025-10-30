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