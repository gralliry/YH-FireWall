package ctable

import (
	"YH-FireWall/core/connection"
	"context"
	"fmt"
	"github.com/google/gopacket/layers"
	"log"
	"net"
	"sync"
)

var (
	conn  map[string]connection.Connection
	mutex sync.RWMutex

	params chan *connection.Parameter
	ctx    context.Context
	cancel context.CancelFunc

	isWork = false
)

func Start(parent context.Context) error {
	ctx, cancel = context.WithCancel(parent)
	conn = make(map[string]connection.Connection)
	params = make(chan *connection.Parameter, 10000)
	go handle()
	return nil
}

func Close() error {
	cancel()
	close(params)
	return nil
}

func handle() {
	for {
		select {
		case <-ctx.Done():
			return
		case param := <-params:
			// 检查是否已存在连接，避免重复写入
			lkey := param.Key()
			if c, exists := conn[lkey]; exists {
				c.Update()
				continue
			}
			c, err := connection.New(param)
			if err != nil {
				continue
			}
			rkey := param.ReverseKey()
			conn[lkey] = c
			conn[rkey] = c
		}
	}
}

func Push(
	proto layers.IPProtocol,
	srcIP net.IP,
	srcPort uint16,
	dstIP net.IP,
	dstPort uint16,

	inDev *uint32,
	outDev *uint32,
) {
	//
	if !isWork {
		// 打印 报文
		pringLog(proto, srcIP, srcPort, dstIP, dstPort, inDev, outDev)
		return
	}
	// 对协议判断
	// 参数可以不保证 inDev 、 outDev 不同时为 nil
	var direction connection.Direction
	switch {
	case inDev != nil && outDev != nil:
		direction = connection.Forward
	case inDev != nil:
		direction = connection.Inbound
	case outDev != nil:
		direction = connection.Outbound
	default:
		return
	}
	select {
	case <-ctx.Done():
	case params <- &connection.Parameter{
		Proto:     proto,
		SrcIP:     srcIP,
		SrcPort:   srcPort,
		DstIP:     dstIP,
		DstPort:   dstPort,
		Direction: direction,
	}:
	default:
	}
}

func pringLog(proto layers.IPProtocol, srcIP net.IP, srcPort uint16, dstIP net.IP, dstPort uint16, inDev *uint32, outDev *uint32) {
	var direction string
	switch {
	case inDev == nil && outDev == nil:
		direction = "   ->   "
	case inDev == nil:
		direction = fmt.Sprintf("%3d->   ", *outDev)
	case outDev == nil:
		direction = fmt.Sprintf("   ->%-3d", *inDev)
	default:
		direction = fmt.Sprintf("%3d->%-3d", *inDev, *outDev)
	}
	log.Printf("[%5s] %15s:%5d -> %15s:%5d (%s)", proto, srcIP, srcPort, dstIP, dstPort, direction)
}
