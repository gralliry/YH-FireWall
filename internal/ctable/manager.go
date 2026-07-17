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

	nnet "github.com/shirou/gopsutil/v4/net"
)

type Manager struct {
	ctx    context.Context
	cancel context.CancelFunc
	mutex  sync.RWMutex

	table   *multikeymap.Map[string, *conn.Conn]
	channel chan *flow.Flow
}

func New() *Manager {
	ctx, cancel := context.WithCancel(context.Background())
	t := &Manager{
		ctx:    ctx,
		cancel: cancel,

		table:   multikeymap.New[string, *conn.Conn](),
		channel: make(chan *flow.Flow, 1024),
	}
	go t.clean()
	return t
}

func (m *Manager) Close() error {
	m.cancel()
	close(m.channel)
	return nil
}

func (m *Manager) Remove(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	conn, exists := m.table.Get(id)
	if !exists {
		return errors.New("connection not found")
	}
	if conn.Closed() {
		return errors.New("connection already closed")
	}
	if err := conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}

func (m *Manager) List() []*conn.Info {
	m.mutex.Lock()
	defer m.mutex.Unlock()
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
		return infoList
	}
	// 获取网卡ip映射 // 写入map
	connMap := multikeymap.New[string, *conn.Conn]()
	for _, conn := range connList {
		connMap.Set(conn, conn.LKey(), conn.RKey())
	}
	for _, bc := range bconnList {
		// 获取连接
		c, exist := connMap.Get(bc.Laddr.String())
		if !exist {
			continue
		}
		ci := c.Info(bc.Pid)
		// 加入
		infoList = append(infoList, ci)
	}
	return infoList
}

func (m *Manager) Push(f *flow.Flow) {
	// 添加连接键
	key := f.CKey()
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
	} else {
		// 添加连接 // 获取方向
		c = conn.New(f)
		// 写入表
		m.table.Set(c, c.LKey(), c.RKey(), c.ID())
	}
}

func (m *Manager) Match(f *flow.Flow) bool {
	return false
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
