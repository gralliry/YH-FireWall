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
	ctx    context.Context
	cancel context.CancelFunc
	config Config
	nfq    *nfqueue.Nfqueue
)

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

func Start(cfg Config) error {
	ctx, cancel = context.WithCancel(context.Background())
	config = cfg

	var err error
	nfq, err = nfqueue.Open(&nfqueue.Config{
		NfQueue:      config.No,
		MaxPacketLen: 2048,
		MaxQueueLen:  2048,
		Copymode:     nfqueue.NfQnlCopyPacket,
		WriteTimeout: 150 * time.Millisecond,
	})
	if err != nil {
		return err
	}
	if err := nfq.SetOption(netlink.NoENOBUFS, true); err != nil {
		return err
	}
	if err := nfq.RegisterWithErrorFunc(ctx, handleFunc, errorFunc); err != nil {
		return err
	}
	cmd := exec.Command("bash", "-c", fmt.Sprintf(cmdSet, config.No))
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func Close() error {
	var errs []error
	cmd := exec.Command("bash", "-c", cmdUnset)
	if err := cmd.Run(); err != nil {
		errs = append(errs, err)
	}
	if nfq != nil {
		if err := nfq.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if cancel != nil {
		cancel()
	}
	return errors.Join(errs...)
}
