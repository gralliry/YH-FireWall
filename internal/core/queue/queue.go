package queue

import (
	"YH-FireWall/internal/core/group"
	"YH-FireWall/internal/core/packet"
	"context"
	"github.com/florianl/go-nfqueue"
	"log"
	"time"
)

type Queue struct {
	nfq *nfqueue.Nfqueue

	ctx    context.Context
	cancel context.CancelFunc

	group *group.Group

	attributes chan *nfqueue.Attribute // 缓存队列
}

func New(g *group.Group) (*Queue, error) {
	ctx, cancel := context.WithCancel(context.Background())
	q := &Queue{
		ctx:        ctx,
		cancel:     cancel,
		group:      g,
		attributes: make(chan *nfqueue.Attribute, 2048),
	}
	// 打开队列
	if nf, err := nfqueue.Open(&nfqueue.Config{
		NfQueue:      g.Qnum,
		MaxPacketLen: 2048,
		MaxQueueLen:  2048,
		Copymode:     nfqueue.NfQnlCopyPacket,
		WriteTimeout: 15 * time.Millisecond,
	}); err != nil {
		return nil, err
	} else {
		q.nfq = nf
	}
	// 注册处理函数
	if err := q.nfq.RegisterWithErrorFunc(
		ctx,
		func(a nfqueue.Attribute) int {
			select {
			case q.attributes <- &a:
				return nfqueue.NfStolen
			default:
				return nfqueue.NfDrop
			}
		},
		func(e error) int {
			return nfqueue.NfDrop
		},
	); err != nil {
		return nil, err
	}
	return q, nil
}

func (q *Queue) Close() error {
	q.cancel()
	return q.nfq.Close()
}

func (q *Queue) handle() {
	for {
		select {
		case <-q.ctx.Done():
			return
		case a := <-q.attributes:
			if a == nil {
				continue
			}
			p, err := packet.Parse(a)
			if err != nil {
				if p != nil {
					_ = q.nfq.SetVerdict(p.Id, nfqueue.NfDrop)
				} else {
					continue
				}
			}
			// todo 打印包日志
			log.Println(p)
			// 匹配规则组
			match, accept := q.group.Match(p)
			// 匹配规则
			if match && accept {
				_ = q.nfq.SetVerdict(p.Id, nfqueue.NfAccept)
			} else {
				_ = q.nfq.SetVerdict(p.Id, nfqueue.NfDrop)
			}
		}
	}
}
