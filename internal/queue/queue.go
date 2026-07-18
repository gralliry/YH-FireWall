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

type Handler interface {
	HandleFlow(*nfqueue.Attribute) (bool, bool)
}

func New(cfg Config, handler Handler) (*NFQ, error) {
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
	if err := iptables(fmt.Sprintf(cmdSet, cfg.Num, cfg.Name)); err != nil {
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
	q.cancel()
	if err := q.nfq.Close(); err != nil {
		errs = append(errs, err)
	}
	if err := iptables(fmt.Sprintf(cmdUnset, q.config.Num, q.config.Name)); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

func iptables(cmd string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, "bash", "-c", cmd).Run()
}

func hookFunc(nfq *nfqueue.Nfqueue, handler Handler) nfqueue.HookFunc {
	return func(a nfqueue.Attribute) int {
		if a.PacketID == nil {
			return -1
		}
		if accept, ok := handler.HandleFlow(&a); !ok {
			nfq.SetVerdict(*a.PacketID, nfqueue.NfDrop)
			return -1
		} else if accept {
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
