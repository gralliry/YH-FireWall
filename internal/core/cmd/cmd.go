package cmd

import (
	"fmt"
	"os"
	"os/exec"
)

var cmdSet = `
sudo iptables -I INPUT   -j NFQUEUE --queue-num %d -m comment --comment "yfw"
sudo iptables -I OUTPUT  -j NFQUEUE --queue-num %d -m comment --comment "yfw"
sudo iptables -I FORWARD -j NFQUEUE --queue-num %d -m comment --comment "yfw"
`

func Set(qnum uint16) error {
	// 使用 bash 执行多行命令
	cmd := exec.Command("bash", "-c", fmt.Sprintf(cmdSet, qnum, qnum, qnum))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}
	return nil
}

var cmdUnset = `
sudo iptables -C INPUT   -j NFQUEUE --queue-num %d -m comment --comment "yfw" 2>/dev/null && sudo iptables -D INPUT   -j NFQUEUE --queue-num %d -m comment --comment "yfw" || true
sudo iptables -C OUTPUT  -j NFQUEUE --queue-num %d -m comment --comment "yfw" 2>/dev/null && sudo iptables -D OUTPUT  -j NFQUEUE --queue-num %d -m comment --comment "yfw" || true
sudo iptables -C FORWARD -j NFQUEUE --queue-num %d -m comment --comment "yfw" 2>/dev/null && sudo iptables -D FORWARD -j NFQUEUE --queue-num %d -m comment --comment "yfw" || true
`

func Unset(qnum uint16) error {
	// 使用 bash 执行多行命令
	cmd := exec.Command("bash", "-c", fmt.Sprintf(cmdUnset, qnum, qnum, qnum, qnum, qnum, qnum))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run command: %v", err)
	}
	return nil
}
