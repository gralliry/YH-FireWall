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

const cmdSet = `
iptables -C INPUT   -j NFQUEUE --queue-num %[1]d --queue-bypass -m comment --comment "%[2]s" 2>/dev/null || iptables -I INPUT   -j NFQUEUE --queue-num %[1]d --queue-bypass -m comment --comment "%[2]s"
iptables -C OUTPUT  -j NFQUEUE --queue-num %[1]d --queue-bypass -m comment --comment "%[2]s" 2>/dev/null || iptables -I OUTPUT  -j NFQUEUE --queue-num %[1]d --queue-bypass -m comment --comment "%[2]s"
iptables -C FORWARD -j NFQUEUE --queue-num %[1]d --queue-bypass -m comment --comment "%[2]s" 2>/dev/null || iptables -I FORWARD -j NFQUEUE --queue-num %[1]d --queue-bypass -m comment --comment "%[2]s"
`

const cmdUnset = `
iptables -D INPUT   -j NFQUEUE --queue-num %[1]d --queue-bypass -m comment --comment "%[2]s" 2>/dev/null || true
iptables -D OUTPUT  -j NFQUEUE --queue-num %[1]d --queue-bypass -m comment --comment "%[2]s" 2>/dev/null || true
iptables -D FORWARD -j NFQUEUE --queue-num %[1]d --queue-bypass -m comment --comment "%[2]s" 2>/dev/null || true
`

type NFQ struct {
	ctx    context.Context
	cancel context.CancelFunc
	config Config
	nfq    *nfqueue.Nfqueue
}

type HandleFunc func(*nfqueue.Attribute) (bool, error)

func New(cfg Config, handler HandleFunc) (*NFQ, error) {
	var err error
	nfq, err := nfqueue.Open(&nfqueue.Config{
		NfQueue:      cfg.Num,
		MaxPacketLen: 2048,
		MaxQueueLen:  2048,
		Copymode:     nfqueue.NfQnlCopyPacket,
		WriteTimeout: 150 * time.Millisecond,
	})
	if err != nil {
		return nil, err
	}
	if err := nfq.SetOption(netlink.NoENOBUFS, true); err != nil {
		nfq.Close()
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	if err := nfq.RegisterWithErrorFunc(ctx, hookFunc(nfq, handler), errorFunc()); err != nil {
		nfq.Close()
		cancel()
		return nil, err
	}
	cmd := exec.Command("bash", "-c", fmt.Sprintf(cmdSet, cfg.Num, cfg.Name))
	if err := cmd.Run(); err != nil {
		nfq.Close()
		cancel()
		return nil, err
	}
	return &NFQ{
		ctx:    ctx,
		cancel: cancel,
		config: cfg,
		nfq:    nfq,
	}, nil
}

func (q *NFQ) Close() error {
	var errs []error
	if err := q.nfq.Close(); err != nil {
		errs = append(errs, err)
	}
	q.cancel()
	cmd := exec.Command("bash", "-c", fmt.Sprintf(cmdUnset, q.config.Num, q.config.Name))
	if err := cmd.Run(); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func hookFunc(nfq *nfqueue.Nfqueue, handler HandleFunc) nfqueue.HookFunc {
	return func(a nfqueue.Attribute) int {
		if ok, err := handler(&a); err != nil {
			return -1
		} else if ok {
			nfq.SetVerdict(*a.PacketID, nfqueue.NfAccept)
			return 0
		} else {
			nfq.SetVerdict(*a.PacketID, nfqueue.NfDrop)
			return 0
		}
	}
}

func errorFunc() nfqueue.ErrorFunc {
	return func(err error) int {
		return -1
	}
}
