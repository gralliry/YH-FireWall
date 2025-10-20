package ctable

import (
	"YH-FireWall/core/connection"
	"YH-FireWall/core/pkg/funcall"
	"errors"
	"fmt"
	"github.com/google/gopacket/layers"
	"net"
	"sync"
)

var (
	table map[string]*connection.Connection
	mutex sync.RWMutex
)

func Init() error {
	table = make(map[string]*connection.Connection)
	return nil
}

func Close() error {
	return nil
}

func makeKey(proto layers.IPProtocol, srcIP net.IP, srcPort uint16, dstIP net.IP, dstPort uint16) string {
	return fmt.Sprintf("%s-%s-%d-%s-%d", proto, srcIP, srcPort, dstIP, dstPort)
}

func Push(
	family uint8, proto layers.IPProtocol,
	srcIP net.IP, srcPort uint16, dstIP net.IP, dstPort uint16,
	inDev, outDev *uint32,
) bool {
	mutex.Lock()
	defer mutex.Unlock()
	// 添加连接键
	lkey := makeKey(proto, srcIP, srcPort, dstIP, dstPort)
	rkey := makeKey(proto, dstIP, dstPort, srcIP, srcPort)
	conn, exists := table[lkey]
	if exists {
		// 检测连接状态
		isClosed, isExpired := conn.Status()
		switch {
		case isClosed && isExpired:
			// 如果被关闭 且 已过期
			delete(table, lkey)
			delete(table, rkey)
			return true
		case isClosed && !isExpired:
			// 表示该连接任然未过期
			return false
		default:
			// 如果没有过期，更新
			conn.Update()
			return true
		}
	}
	// 添加连接 // 获取方向
	switch {
	case inDev != nil && outDev != nil:
		// 转发模式
		conn = connection.New(family, proto, srcIP, srcPort, connection.Forward, dstIP, dstPort)
	case inDev != nil:
		// 入口模式 // 源是外部连接
		conn = connection.New(family, proto, dstIP, dstPort, connection.Inbound, srcIP, srcPort)
	case outDev != nil:
		// 出口模式 // 源是内部连接
		conn = connection.New(family, proto, srcIP, srcPort, connection.Outbound, dstIP, dstPort)
	default:
		// 未知数据，直接停止
		return false
	}
	// 写入表
	table[lkey] = conn
	table[rkey] = conn
	//
	return true
}

func GetAll() []connection.Config {
	mutex.RLock()
	defer mutex.RUnlock()
	return funcall.Convert(funcall.Set(funcall.Map2List(table,
		func(_ string, conn *connection.Connection) *connection.Config {
			return conn.Unparse()
		}),
		func(conn *connection.Config) string {
			return conn.Id
		}),
		func(conn *connection.Config) connection.Config {
			return *conn
		})
}

func Disable(id string) error {
	mutex.RLock()
	defer mutex.RUnlock()
	conn, exists := table[id]
	if !exists {
		return errors.New("connection not found")
	}
	if conn.Closed() {
		return errors.New("connection already closed")
	}
	conn.Close()
	return nil
}

// todo 自动清理过期连接
