package queue

import (
	"YH-FireWall/core/ctable"
	"YH-FireWall/core/rtable"
	"context"
	"errors"
	"fmt"
	"github.com/florianl/go-nfqueue"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/mdlayher/netlink"
	"net"
	"os"
	"os/exec"
	"time"
)

var (
	NfQueueNo uint16
)

var (
	nfq *nfqueue.Nfqueue
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

func Start(ctx context.Context) (err error) {
	// 打开队列
	nfq, err = nfqueue.Open(&nfqueue.Config{
		NfQueue:      NfQueueNo,
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
	cmd := exec.Command("bash", "-c", fmt.Sprintf(cmdSet, NfQueueNo, NfQueueNo, NfQueueNo, NfQueueNo, NfQueueNo, NfQueueNo))
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
	var (
		srcIP    net.IP
		srcPort  uint16
		dstIP    net.IP
		dstPort  uint16
		inDev    = a.InDev
		outDev   = a.OutDev
		protocol layers.IPProtocol
	)
	// 使用 gopacket 解析 Payload
	rawpacket := gopacket.NewPacket(*a.Payload, layers.LayerTypeIPv4, gopacket.Default)
	// 获取 IPv4 或 IPv6 地址
	if ip4 := rawpacket.Layer(layers.LayerTypeIPv4); ip4 != nil {
		ip := ip4.(*layers.IPv4)
		srcIP, dstIP, protocol = ip.SrcIP, ip.DstIP, ip.Protocol
	} else if ip6 := rawpacket.Layer(layers.LayerTypeIPv6); ip6 != nil {
		ip := ip6.(*layers.IPv6)
		srcIP, dstIP, protocol = ip.SrcIP, ip.DstIP, ip.NextHeader
	} else {
		_ = nfq.SetVerdict(*a.PacketID, nfqueue.NfDrop)
		return 0
	}
	// 匹配端口 // TCP/UDP 端口
	switch protocol {
	case layers.IPProtocolTCP:
		tcp := rawpacket.Layer(layers.LayerTypeTCP)
		if tcp == nil {
			_ = nfq.SetVerdict(*a.PacketID, nfqueue.NfDrop)
			return 0
		}
		t := tcp.(*layers.TCP)
		srcPort, dstPort = uint16(t.SrcPort), uint16(t.DstPort)
	case layers.IPProtocolUDP:
		udp := rawpacket.Layer(layers.LayerTypeUDP)
		if udp == nil {
			_ = nfq.SetVerdict(*a.PacketID, nfqueue.NfDrop)
			return 0
		}
		u := udp.(*layers.UDP)
		srcPort, dstPort = uint16(u.SrcPort), uint16(u.DstPort)
	}
	// 匹配规则
	if !rtable.Match(srcIP, srcPort, dstIP, dstPort, inDev, outDev, protocol) {
		_ = nfq.SetVerdict(*a.PacketID, nfqueue.NfDrop)
		return 0
	}
	// 写入 ctable // 这个因为没有结果上的关系，可以使用异步处理
	ctable.Push(protocol, srcIP, srcPort, dstIP, dstPort, inDev, outDev)
	// 继续处理
	_ = nfq.SetVerdict(*a.PacketID, nfqueue.NfAccept)
	return 0
}
