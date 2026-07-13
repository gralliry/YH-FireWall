package ctable

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/go-co-op/gocron/v2"
)

var (
	table     map[string]*Connection
	namespace map[string]*Connection
	
	mutex     sync.RWMutex
	scheduler gocron.Scheduler
)

func Start(ctx context.Context) error {
	table = make(map[string]*Connection)
	namespace = make(map[string]*Connection)
	clean()
	s, err := gocron.NewScheduler()
	if err != nil {
		return err
	}
	_, err = s.NewJob(
		gocron.DurationJob(time.Minute),
		gocron.NewTask(clean),
	)
	if err != nil {
		return err
	}
	s.Start()
	scheduler = s
	go func() {
		<-ctx.Done()
		s.Shutdown()
	}()
	return nil
}

func Close() error {
	if scheduler != nil {
		return scheduler.Shutdown()
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

func clean() {
	mutex.Lock()
	defer mutex.Unlock()
	for id, conn := range namespace {
		if conn.Expired() {
			delete(namespace, id)
			delete(table, conn.LKey())
			delete(table, conn.RKey())
		}
	}
}
