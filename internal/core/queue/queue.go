package queue

import (
	"YH-FireWall/internal/core/manager"
	"YH-FireWall/internal/core/packet"
	"context"
	"github.com/florianl/go-nfqueue"
	"github.com/mdlayher/netlink"
	"time"
)

var (
	nfq    *nfqueue.Nfqueue
	ctx    context.Context
	cancel context.CancelFunc

	defaultAccept = true
)

func Start(parent context.Context, num uint16, accept bool) (err error) {
	ctx, cancel = context.WithCancel(parent)
	//
	defaultAccept = accept
	// 打开队列
	nfq, err = nfqueue.Open(&nfqueue.Config{
		NfQueue:      num,
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
	if err = nfq.RegisterWithErrorFunc(ctx, handleAttribute, handleError); err != nil {
		return err
	}
	return nil
}

func Close() error {
	cancel()
	return nfq.Close()
}

func handleError(_ error) int {
	return -1
}

func handleAttribute(a nfqueue.Attribute) int {
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
			_ = nfq.SetVerdict(p.Id(), nfqueue.NfAccept)
		} else {
			_ = nfq.SetVerdict(p.Id(), nfqueue.NfDrop)
		}
	} else if defaultAccept {
		_ = nfq.SetVerdict(p.Id(), nfqueue.NfAccept)
	} else {
		_ = nfq.SetVerdict(p.Id(), nfqueue.NfDrop)
	}
	return 0
}

//func handleAttribute(a nfqueue.Attribute) int {
//	select {
//	case attributes <- &a:
//	default:
//		_ = nfq.SetVerdict(*a.PacketID, nfqueue.NfDrop)
//	}
//	return 0
//}
//
//func handlePackets() {
//	for {
//		select {
//		case <-ctx.Done():
//			return
//		case a := <-attributes:
//			p, err := packet.Parse(a)
//			if err != nil {
//				_ = nfq.SetVerdict(*a.PacketID, nfqueue.NfDrop)
//				break
//			}
//			// 打印包日志
//			log.Println(p)
//			// 匹配规则组
//			match, accept := manager.Match(p)
//			// 匹配规则
//			if match && accept {
//				_ = nfq.SetVerdict(p.Id(), nfqueue.NfAccept)
//			} else {
//				_ = nfq.SetVerdict(p.Id(), nfqueue.NfDrop)
//			}
//		}
//	}
//}
