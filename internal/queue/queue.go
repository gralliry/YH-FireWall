package queue

import (
	"YH-FireWall/internal/group"
	"YH-FireWall/internal/packet"
	"context"
	"github.com/florianl/go-nfqueue"
	"log"
	"time"
)

type Queue struct {
	no uint16

	nfq *nfqueue.Nfqueue

	ctx context.Context

	group *group.Group

	attributes chan *nfqueue.Attribute // 缓存队列
}

func New(ctx context.Context, no uint16) (*Queue, error) {
	q := &Queue{
		no:  no,
		ctx: ctx,

		attributes: make(chan *nfqueue.Attribute, 1000),
	}
	nf, err := nfqueue.Open(&nfqueue.Config{
		NfQueue:      no,
		MaxPacketLen: 0xFFFF,
		MaxQueueLen:  0xFFF,
		Copymode:     nfqueue.NfQnlCopyPacket,
		WriteTimeout: 15 * time.Millisecond,
	})
	if err != nil {
		return nil, err
	}
	if err = nf.RegisterWithErrorFunc(
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
	q.nfq = nf
	return q, nil
}

func (q *Queue) Close() error {
	return q.nfq.Close()
}

func (q *Queue) handle() {
	for a := range q.attributes {
		p := packet.Parse(a)
		log.Println(p)
		if p == nil {
			_ = q.nfq.SetVerdict(p.ID(), nfqueue.NfDrop)
			continue
		}
		// 匹配规则组
		match, accept := q.group.Match(p)
		// 匹配规则
		if match && accept {
			_ = q.nfq.SetVerdict(p.ID(), nfqueue.NfAccept)
		} else {
			_ = q.nfq.SetVerdict(p.ID(), nfqueue.NfDrop)
		}
	}
}
