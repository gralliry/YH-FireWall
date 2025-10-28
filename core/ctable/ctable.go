package ctable

import (
	"YH-FireWall/core/connection"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

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

func GetAll() []connection.Config {
	mutex.RLock()
	defer mutex.RUnlock()
	// push by process
	pushByProcess()
	// Step 1: 提取所有连接（values） // Distinct 跳过重复的连接
	connMap := make(map[string]*connection.Connection)
	for _, v := range table {
		connMap[v.Id()] = v
	}
	// Distinct 跳过重复的连接
	configList := make([]connection.Config, 0)
	for _, conn := range connMap {
		if conn.Closed() {
			continue
		}
		configList = append(configList, *conn.Unparse())
	}
	return configList
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
	if err := conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}
	return nil
}

// 自动清理过期连接
func clean(ctx context.Context) {
	for {
		// 先执行清理逻辑（立即执行一次）
		mutex.Lock()
		count := 0
		for id, conn := range namespcae {
			if conn.Expired() {
				delete(namespcae, id)
				delete(table, conn.LKey())
				delete(table, conn.RKey())
				count += 1
			}
		}
		mutex.Unlock()
		// 日志输出
		log.Printf("clean %d expired connections", count)
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Minute):
		}
	}
}
