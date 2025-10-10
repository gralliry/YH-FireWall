package queue

import (
	"YH-FireWall/core/config"
	"YH-FireWall/core/manager"
	"YH-FireWall/core/packet"
	"context"
	"errors"
	"fmt"
	"github.com/florianl/go-nfqueue"
	"github.com/mdlayher/netlink"
	"os"
	"os/exec"
	"time"
)

var (
	nfq             *nfqueue.Nfqueue
	isDefaultAccept bool
	queueNum        uint16
)

var cmdSet = `
sudo iptables -C INPUT   -j NFQUEUE --queue-num %d -m comment --comment "yfw" 2>/dev/null || sudo iptables -I INPUT   -j NFQUEUE --queue-num %d -m comment --comment "yfw"
sudo iptables -C OUTPUT  -j NFQUEUE --queue-num %d -m comment --comment "yfw" 2>/dev/null || sudo iptables -I OUTPUT  -j NFQUEUE --queue-num %d -m comment --comment "yfw"
sudo iptables -C FORWARD -j NFQUEUE --queue-num %d -m comment --comment "yfw" 2>/dev/null || sudo iptables -I FORWARD -j NFQUEUE --queue-num %d -m comment --comment "yfw"
`

var cmdUnset = `
sudo iptables -L INPUT   --line-numbers | grep "NFQUEUE.*yfw" | awk '{print $1}' | xargs -r sudo iptables -D INPUT
sudo iptables -L OUTPUT  --line-numbers | grep "NFQUEUE.*yfw" | awk '{print $1}' | xargs -r sudo iptables -D OUTPUT
sudo iptables -L FORWARD --line-numbers | grep "NFQUEUE.*yfw" | awk '{print $1}' | xargs -r sudo iptables -D FORWARD
`

func Start(ctx context.Context, cfg config.Queue) (err error) {
	// 队列编号
	queueNum = cfg.Num
	// 默认规则
	isDefaultAccept = cfg.Accept
	// 打开队列
	nfq, err = nfqueue.Open(&nfqueue.Config{
		NfQueue:      queueNum,
		MaxPacketLen: 2048,
		MaxQueueLen:  2048,
		Copymode:     nfqueue.NfQnlCopyPacket,
		WriteTimeout: 150 * time.Millisecond,
	})
	// Avoid receiving ENOBUFS errors.
	if err = nfq.SetOption(netlink.NoENOBUFS, true); err != nil {
		return err
	}
	// 注册处理函数
	if err = nfq.RegisterWithErrorFunc(ctx, handler, func(err error) int { return -1 }); err != nil {
		return err
	}
	// 设置包导向
	cmd := exec.Command("bash", "-c", fmt.Sprintf(cmdSet, queueNum, queueNum, queueNum, queueNum, queueNum, queueNum))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Run(); err != nil {
		return err
	}
	return nil
}

func Close() error {
	var errs []error
	// 使用 bash 执行多行命令
	cmd := exec.Command("bash", "-c", cmdUnset)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		errs = append(errs, err)
	}
	if err := nfq.Close(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func handler(a nfqueue.Attribute) int {
	if a.PacketID == nil {
		return 0
	}
	// 解析报文
	p, err := packet.Parse(&a)
	if err != nil {
		_ = nfq.SetVerdict(*a.PacketID, nfqueue.NfDrop)
		return 0
	}
	// 匹配规则组
	match, accept := manager.Match(p)
	// 匹配规则
	if match {
		if accept {
			_ = nfq.SetVerdict(*a.PacketID, nfqueue.NfAccept)
		} else {
			_ = nfq.SetVerdict(*a.PacketID, nfqueue.NfDrop)
		}
	}
	// 默认规则
	if isDefaultAccept {
		_ = nfq.SetVerdict(*a.PacketID, nfqueue.NfAccept)
	} else {
		_ = nfq.SetVerdict(*a.PacketID, nfqueue.NfDrop)
	}
	return 0
}
