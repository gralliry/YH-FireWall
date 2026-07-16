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
	go t.handle()

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
	// push by process
	pushByProcess()
	// Step 1: 提取所有连接（values）
	connList := m.table.Values()
	// Distinct 跳过重复的连接
	connList = funcs.Filter(connList, func(c *conn.Conn) bool {
		return !c.Expired()
	})
	// 映射到Info
	infoList := funcs.Transform(connList, func(c *conn.Conn) *conn.Info {
		return c.Info()
	})
	return infoList
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
			m.mutex.Lock()
			m.table.Filter(func(c *conn.Conn) bool {
				return !c.Expired()
			})
			m.mutex.Unlock()
		case <-m.ctx.Done():
			return
		}
	}
}

func (m *Manager) handle() {
	for {
		select {
		case flow := <-m.channel:
			// 添加连接键
			key := flow.Key()
			m.mutex.Lock()

			if c, exists := m.table.Get(key); exists {
				// 检测连接状态
				isClosed, isExpired := c.Closed(), c.Expired()
				if isClosed {
					m.mutex.Unlock()
					continue
				}
				// 如果没有过期，更新
				if !isExpired {
					c.Active()
					m.mutex.Unlock()
					continue
				}
				// 移除连接
				m.table.Del(c.ID())
			} else {
				// 添加连接 // 获取方向
				c = conn.New(flow)
				// 写入表
				m.table.Set(c, c.LKey(), c.RKey(), c.ID())
			}
			m.mutex.Unlock()
		case <-m.ctx.Done():
			// 退出
			return
		}
	}
}
