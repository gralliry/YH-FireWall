package queue

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"

	"github.com/florianl/go-nfqueue"
	"github.com/mdlayher/netlink"
)

var (
	nfq *nfqueue.Nfqueue
)

// sudo iptables -L -n -v

const cmdSet = `
sudo iptables -C INPUT   -j NFQUEUE --queue-num %[1]d -m comment --comment "yfw" 2>/dev/null || sudo iptables -I INPUT   -j NFQUEUE --queue-num %[1]d -m comment --comment "yfw"
sudo iptables -C OUTPUT  -j NFQUEUE --queue-num %[1]d -m comment --comment "yfw" 2>/dev/null || sudo iptables -I OUTPUT  -j NFQUEUE --queue-num %[1]d -m comment --comment "yfw"
sudo iptables -C FORWARD -j NFQUEUE --queue-num %[1]d -m comment --comment "yfw" 2>/dev/null || sudo iptables -I FORWARD -j NFQUEUE --queue-num %[1]d -m comment --comment "yfw"
`

const cmdUnset = `
sudo iptables -L INPUT   --line-numbers | grep "NFQUEUE.*yfw" | awk '{print $1}' | xargs -r sudo iptables -D INPUT
sudo iptables -L OUTPUT  --line-numbers | grep "NFQUEUE.*yfw" | awk '{print $1}' | xargs -r sudo iptables -D OUTPUT
sudo iptables -L FORWARD --line-numbers | grep "NFQUEUE.*yfw" | awk '{print $1}' | xargs -r sudo iptables -D FORWARD
`

type Queue struct {
	// 上下文管理
	ctx    context.Context
	cancel context.CancelFunc
	// 配置管理
	config  Config
	handler Handler
	// 队列
	nfq *nfqueue.Nfqueue
}

func New(config *Config, handler Handler) *Queue {
	ctx, cancel := context.WithCancel(context.Background())
	return &Queue{
		ctx:     ctx,
		cancel:  cancel,
		config:  *config,
		handler: handler,
	}
}

func (q *Queue) Start() error {
	// 打开队列
	nfq, err := nfqueue.Open(&nfqueue.Config{
		NfQueue:      q.config.No,
		MaxPacketLen: 2048,
		MaxQueueLen:  2048,
		Copymode:     nfqueue.NfQnlCopyPacket,
		WriteTimeout: 150 * time.Millisecond,
	})
	if err != nil {
		return err
	}
	// Avoid receiving ENOBUFS errors.
	if err := nfq.SetOption(netlink.NoENOBUFS, true); err != nil {
		return err
	}
	// 注册处理函数
	if err := nfq.RegisterWithErrorFunc(q.ctx, handleFunc, errorFunc); err != nil {
		return err
	}
	// 设置包导向
	cmd := exec.Command("bash", "-c", fmt.Sprintf(cmdSet, q.config.No))
	if err := cmd.Run(); err != nil {
		return err
	}
	// 写入注册器
	q.nfq = nfq
	return nil
}

func Close() error {
	var errs []error
	// 使用 bash 执行多行命令
	cmd := exec.Command("bash", "-c", cmdUnset)
	if err := cmd.Run(); err != nil {
		errs = append(errs, err)
	}
	if err := nfq.Close(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}
