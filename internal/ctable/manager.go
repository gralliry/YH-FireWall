package ctable

import (
	"YH-FireWall/internal/model/conn"
	"YH-FireWall/internal/model/flow"
	"YH-FireWall/internal/pkg/funcs"
	"YH-FireWall/internal/pkg/multikeymap"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/gopacket/layers"
	nnet "github.com/shirou/gopsutil/v4/net"
)

type Manager struct {
	ctx    context.Context
	cancel context.CancelFunc
	mutex  sync.RWMutex

	table   *multikeymap.Map[string, *conn.Conn]
}

func New() *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	t := &Manager{
		ctx:    ctx,
		cancel: cancel,

		table:   multikeymap.New[string, *conn.Conn](),
	}
	go t.clean()
	return t
}

func (m *Manager) Close() error {
	m.cancel()
	return nil
}

func (m *Manager) Match(f *flow.Flow) (accept bool, exist bool) {
	// 校验flow是否是连接包
	if !f.IsConnection() {
		return false, false
	}
	// 加锁
	m.mutex.Lock()
	defer m.mutex.Unlock()
	// 检查是否在table里面
	conn, exists := m.table.Get(f.Key())
	if !exists {
		return false, false
	}
	if !conn.Alive() {
		return false, true
	}
	conn.Active()
	return true, true
}

func (m *Manager) Remove(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	c, exists := m.table.Get(id)
	if !exists {
		return errors.New("connection not found")
	}
	if c.Closed() {
		return errors.New("connection already closed")
	}
	m.table.Del(id)
	c.Close()
	conn.Release(c)
	return nil
}

func (m *Manager) List() ([]*conn.Info, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	// Step 1: 提取所有连接（values）
	connList := m.table.Values()
	// Distinct 跳过重复的连接
	connList = funcs.Filter(connList, func(c *conn.Conn) bool {
		return !c.Expired() && !c.Closed()
	})
	// 映射到Info
	infoList := make([]*conn.Info, 0)
	// 查找所有连接
	bconnList, err := nnet.Connections("inet")
	if err != nil {
		return nil, err
	}
	// 获取网卡ip映射 // 写入map
	connMap := multikeymap.New[string, *conn.Conn]()
	for _, conn := range connList {
		connMap.Set(conn, conn.LKey(), conn.RKey())
	}
	for _, bc := range bconnList {
		// 映射 socket type 到协议
		var proto layers.IPProtocol
		switch bc.Type {
		case 1: // SOCK_STREAM
			proto = layers.IPProtocolTCP
		case 2: // SOCK_DGRAM
			proto = layers.IPProtocolUDP
		default:
			continue
		}
		// 拼出和 LKey() 一致格式的 key
		key := fmt.Sprintf("%s-%s-%s", proto, bc.Laddr.String(), bc.Raddr.String())
		c, exist := connMap.Get(key)
		if !exist {
			continue
		}
		ci := c.Info(bc.Pid)
		// 加入
		infoList = append(infoList, ci)
	}
	return infoList, nil
}

func (m *Manager) Push(f *flow.Flow) {
	if !f.IsConnection() {
		return
	}
	// 添加连接键
	key := f.Key()
	// 加锁
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if c, exists := m.table.Get(key); exists {
		// 检测连接状态
		if c.Closed() {
			return
		}
		// 如果没有过期，更新
		if !c.Expired() {
			c.Active()
			return
		}
		// 移除连接
		m.table.Del(c.ID())
		conn.Release(c)
	}
	// 添加连接 // 获取方向
	c, ok := conn.New(f)
	if !ok {
		return
	}
	// 写入表
	m.table.Set(c, c.LKey(), c.RKey(), c.ID())
}

func (m *Manager) clean() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			// 加锁
			m.mutex.Lock()
			// 提取过期的
			disconns := m.table.Extract(func(c *conn.Conn) bool {
				return c.Expired()
			})
			// 解锁
			m.mutex.Unlock()
			// 回收过期的
			for _, c := range disconns {
				conn.Release(c)
			}
		case <-m.ctx.Done():
			return
		}
	}
}
