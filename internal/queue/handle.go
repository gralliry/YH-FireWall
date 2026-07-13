package queue

import (
	"YH-FireWall/internal/ctable"
	"YH-FireWall/internal/flow"
	"YH-FireWall/internal/rtable"
	"log"

	"github.com/florianl/go-nfqueue"
)

type Handler interface {
	Match(flow *flow.Flow) bool
	Update(flow *flow.Flow) bool
}

func handleFunc(a nfqueue.Attribute) int {
	if a.Payload == nil {
		nfq.SetVerdict(*a.PacketID, nfqueue.NfDrop)
		return 0
	}
	flow, ok := flow.New(*a.Payload, a.InDev, a.OutDev)
	if !ok {
		nfq.SetVerdict(*a.PacketID, nfqueue.NfDrop)
		return 0
	}
	if !rtable.Match(flow) {
		nfq.SetVerdict(*a.PacketID, nfqueue.NfDrop)
		return 0
	}
	ctable.Push(flow)
	log.Print(flow.String())
	nfq.SetVerdict(*a.PacketID, nfqueue.NfAccept)
	return 0
}

func errorFunc(err error) int {
	return -1
}
