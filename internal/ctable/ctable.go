package ctable

import (
	"YH-FireWall/internal/model/flow"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

var (
	ctx    context.Context
	cancel context.CancelFunc
	mutex  sync.RWMutex

	table     map[string]*Connection
	namespace map[string]*Connection
	channel   chan *flow.Flow

	scheduler *cron.Cron
)

func Start() error {
	mutex.Lock()
	defer mutex.Unlock()

	table = make(map[string]*Connection)
	namespace = make(map[string]*Connection)

	channel = make(chan *flow.Flow, 1024)

	ctx, cancel = context.WithCancel(context.Background())

	go clean(ctx)
	go handle(ctx)

	return nil
}

func Close() error {
	if scheduler != nil {
		<-scheduler.Stop().Done()
	}
	return nil
}

func Infos() []Info {
	mutex.Lock()
	defer mutex.Unlock()
	// push by process
	pushByProcess()
	// Step 1: 提取所有连接（values） // Distinct 跳过重复的连接
	connMap := make(map[string]*Connection)
	for _, v := range table {
		connMap[v.Id()] = v
	}
	// Distinct 跳过重复的连接
	configList := make([]Info, 0)
	for _, conn := range connMap {
		if conn.Expired() {
			continue
		}
		configList = append(configList, *conn.Unparse())
	}
	return configList
}

func Remove(id string) error {
	mutex.Lock()
	defer mutex.Unlock()
	conn, exists := namespace[id]
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

func clean(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			mutex.Lock()
			defer mutex.Unlock()
			for id, conn := range namespace {
				if !conn.Expired() {
					continue
				}
				delete(namespace, id)
				delete(table, conn.LKey())
				delete(table, conn.RKey())
			}
		case <-ctx.Done():
			return
		}
	}
}

func handle(ctx context.Context) {
	for {
		select {
		case flow := <-channel:
			// 添加连接键
			lkey := flow.LKey()
			rkey := flow.RKey()
			//
			conn, exists := table[lkey]
			if exists {
				// 检测连接状态
				isClosed, isExpired := conn.Closed(), conn.Expired()
				if isClosed {
					continue
				}

				if !isExpired {
					// 如果没有过期，更新
					conn.Active()
					continue
				}

				delete(table, lkey)
				delete(table, rkey)
				delete(namespace, conn.Id())
			}
			// 添加连接 // 获取方向
			conn = New(flow)
			// 写入表
			table[lkey] = conn
			table[rkey] = conn
			namespace[conn.Id()] = conn
		case <-ctx.Done():
			// 退出
			return
		}
	}
}
