package ctable

import (
	"YH-FireWall/core/connection"
	"YH-FireWall/core/pkg/fc"
	"context"
	"errors"
	"github.com/google/gopacket/layers"
	"time"

	"net"
	"sync"
)

var (
	table     map[string]*connection.Connection
	namespcae map[string]*connection.Connection
	mutex     sync.RWMutex
)

func Start(ctx context.Context) error {
	table = make(map[string]*connection.Connection)
	namespcae = make(map[string]*connection.Connection)
	go clean(ctx)
	return nil
}

func Close() error {
	return nil
}

func Push(
	family uint8, proto layers.IPProtocol,
	srcIP net.IP, srcPort uint16, dstIP net.IP, dstPort uint16,
	inDev, outDev *uint32,
) bool {
	mutex.Lock()
	defer mutex.Unlock()
	// 添加连接键
	lkey := connection.MakeKey(proto, srcIP, srcPort, dstIP, dstPort)
	rkey := connection.MakeKey(proto, dstIP, dstPort, srcIP, srcPort)
	conn, exists := table[lkey]
	if exists {
		// 检测连接状态
		isClosed, isExpired := conn.Status()
		switch {
		case isClosed && isExpired:
			// 如果被关闭 且 已过期
			delete(table, lkey)
			delete(table, rkey)
			delete(namespcae, conn.Id())
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
	namespcae[conn.Id()] = conn
	//
	return true
}

func GetAll() []connection.Config {
	mutex.RLock()
	defer mutex.RUnlock()
	// Step 1: 提取所有连接（values）
	connList := fc.Map2List(table, func(k string, v *connection.Connection) *connection.Connection {
		return v
	})
	// Distinct 跳过重复的连接
	connList = fc.Distinct(connList, func(conn *connection.Connection) string {
		return conn.Id()
	})
	// Filter 跳过已关闭的连接
	connList = fc.Filter(connList, func(conn *connection.Connection) bool {
		return !conn.Closed()
	})
	// Map 转换为 Config
	return fc.List2List(connList, func(conn *connection.Connection) connection.Config {
		return *conn.Unparse()
	})
}

func Remove(id string) error {
	mutex.RLock()
	defer mutex.RUnlock()
	conn, exists := namespcae[id]
	if !exists {
		return errors.New("connection not found")
	}
	if conn.Closed() {
		return errors.New("connection already closed")
	}
	conn.Close()
	return nil
}

// 自动清理过期连接
func clean(ctx context.Context) {
	for {
		// 先执行清理逻辑（立即执行一次）
		mutex.Lock()
		for id, conn := range namespcae {
			if conn.Expired() {
				delete(namespcae, id)
				delete(table, conn.LKey())
				delete(table, conn.RKey())
			}
		}
		mutex.Unlock()
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Minute):
		}
	}
}
