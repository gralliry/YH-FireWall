package rtable

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"YH-FireWall/internal/model/flow"
	"YH-FireWall/internal/model/rule"
	"YH-FireWall/internal/pkg/funcs"
	"YH-FireWall/internal/pkg/lfile"
	"YH-FireWall/internal/pkg/skiplist"
)

type DevMap interface {
	Name2Index(name string) (index uint32, exist bool)
	Index2Name(index uint32) (name string, exist bool)
}

type Manager struct {
	mutex  sync.RWMutex
	config Config

	rules *skiplist.SkipList[string, *rule.Rule]
	file  *lfile.LockedFile
	// 设备调用映射
	devs DevMap
	//
	ctx       context.Context
	cancel    context.CancelFunc
	flushReqs chan struct{}
}

func New(config Config, devs DevMap) (*Manager, error) {
	file, err := lfile.Open(config.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open rule file: %w", err)
	}
	buf, err := file.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read rule file: %w", err)
	}
	rds := make([]*rule.Data, 0)
	if len(buf) > 0 {
		if err := json.Unmarshal(buf, &rds); err != nil {
			return nil, fmt.Errorf("failed to decode rule file: %w", err)
		}
	}
	// 对其去重
	rds = funcs.Distinct(rds, func(ro *rule.Data) string {
		return ro.ID
	})
	// 加载配置
	rules := skiplist.New[string](func(a, b *rule.Rule) int {
		return a.Compare(b)
	})
	for _, rd := range rds {
		rr, err := rule.Parse(rd, devs.Name2Index)
		if err != nil {
			continue
		}
		rules.Insert(rd.ID, rr)
	}
	ctx, cancel := context.WithCancel(context.Background())
	m := &Manager{
		config: config,
		rules:  rules,
		file:   file,

		devs: devs,

		ctx:       ctx,
		cancel:    cancel,
		flushReqs: make(chan struct{}, 20),
	}
	go m.handleSave()
	return m, nil
}

func (m *Manager) Close() error {
	m.cancel()
	return m.file.Close()
}

func (m *Manager) list() []*rule.Data {
	rrs := m.rules.Values()
	return funcs.Transform(rrs, func(rr *rule.Rule) *rule.Data {
		return rr.Data(m.devs.Index2Name)
	})
}

func (m *Manager) handleSave() {
	for {
		select {
		case <-m.flushReqs:
			// 先消耗所有请求
		drain:
			for {
				select {
				case <-m.flushReqs:
				default:
					// ✅ 跳出 for 循环
					break drain
				}
			}
			// 开始写入
			_ = m.save()
		case <-m.ctx.Done():
			// 要退出了，重新看看有没有写入请求
			select {
			case <-m.flushReqs:
				// 开始写入
				_ = m.save()
			default:
			}
			return
		}
	}
}

func (m *Manager) save() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	rrs := m.list()
	buf, err := json.Marshal(rrs)
	if err != nil {
		return fmt.Errorf("failed to encode rule file: %w", err)
	}
	if err := m.file.Write(buf); err != nil {
		return fmt.Errorf("failed to write rule file: %w", err)
	}
	return nil
}

func (m *Manager) Create(ro *rule.Option) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	rr, err := rule.New(ro, m.devs.Name2Index)
	if err != nil {
		return "", err
	}
	// 插入
	m.rules.Insert(rr.ID(), rr)
	m.flushReqs <- struct{}{}
	return rr.ID(), nil
}

func (m *Manager) Update(id string, ro *rule.Option) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	rr, exists := m.rules.Search(id)
	if !exists {
		return fmt.Errorf("rule %s not exists", id)
	}
	if err := rr.Update(ro, m.devs.Name2Index); err != nil {
		return err
	}
	m.flushReqs <- struct{}{}
	return nil
}

func (m *Manager) Delete(id string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	// 获取
	ok := m.rules.Delete(id)
	if !ok {
		return fmt.Errorf("rule %s not exists", id)
	}
	m.flushReqs <- struct{}{}
	return nil
}

func (m *Manager) Search(id string) *rule.Data {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	//
	rr, exists := m.rules.Search(id)
	if !exists {
		return nil
	}
	return rr.Data(m.devs.Index2Name)
}

func (m *Manager) List() []*rule.Data {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.list()
}

func (m *Manager) Enable(id string, enable bool) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	rr, exists := m.rules.Search(id)
	if !exists {
		return fmt.Errorf("rule %s not exists", id)
	}
	ro := &rule.Option{
		Enable: new(enable),
	}
	if err := rr.Update(ro, m.devs.Name2Index); err != nil {
		return fmt.Errorf("failed to update rule: %w", err)
	}
	m.flushReqs <- struct{}{}
	return nil
}

func (m *Manager) Match(f *flow.Flow) bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	// 匹配
	rr, exist := m.rules.First(func(r *rule.Rule) bool {
		return r.Match(f)
	})
	// 匹配失败使用默认策略
	if !exist {
		return m.config.DefaultAccept
	}
	// 规则策略
	return rr.Accept()
}
