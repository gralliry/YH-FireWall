package rtable

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"YH-FireWall/internal/model/flow"
	"YH-FireWall/internal/model/rule"
	"YH-FireWall/internal/pkg/bimap"
	"YH-FireWall/internal/pkg/funcs"
	"YH-FireWall/internal/pkg/lfile"
	"YH-FireWall/internal/pkg/skiplist"

	"github.com/google/gopacket/layers"
)


type Manager struct {
	mutex  sync.RWMutex
	config Config

	rules *skiplist.SkipList[string, *rule.Rule]
	file  *lfile.LockedFile

	devMap *bimap.Map[uint32, string]
	proMap *bimap.Map[layers, string]
}

func New(config Config, devMap ) (*Manager, error) {
	file, err := lfile.Open(config.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open rule file: %w", err)
	}
	buf, err := file.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read rule file: %w", err)
	}
	ros := make([]*rule.Option, 0)
	if len(buf) > 0 {
		if err := json.Unmarshal(buf, &ros); err != nil {
			return nil, fmt.Errorf("failed to decode rule file: %w", err)
		}
	}
	// 对其去重
	ros = funcs.Distinct(ros, func(ro *rule.Option) string {
		return ro.ID
	})
	// 加载配置
	rules := skiplist.New[string](func(a, b *rule.Rule) int {
		return a.Compare(b)
	})
	for _, ro := range ros {
		rr := rule.New(ro.ID)
		if err := rr.Update(ro); err != nil {
			continue
		}
		rules.Insert(ro.ID, rr)
	}
	return &Manager{
		config: config,
		rules:  rules,
		file:   file,
	}, nil
}

func (m *Manager) Close() error{
	return nil
}

func (m *Manager) Create(ro *rule.Option) (string, error) {
	if ro.ID == "" {
		return "", fmt.Errorf("rule ID is empty")
	}
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if _, ok := m.rules.Search(ro.ID); ok {
		return "", fmt.Errorf("rule %s already exists", ro.ID)
	}
	rr := rule.New(ro.ID)
	if err := rr.Update(ro, Protocol2Name); err != nil {
		return "", err
	}
	// 插入
	m.rules.Insert(ro.ID, rr)
	return ro.ID, nil
}

func (m *Manager) Update(ro *rule.Option) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	rr, exists := m.rules.Search(ro.ID)
	if !exists {
		return fmt.Errorf("rule %s not exists", ro.ID)
	}
	if err := rr.Update(ro, Protocol2Name); err != nil {
		return err
	}
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
	return nil
}

func (m *Manager) Search(id string) *rule.Option {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	rr, exists := m.rules.Search(id)
	if !exists {
		return nil
	}
	return rr.Option(Protocol2Name)
}

func (m *Manager) Select() []*rule.Option {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.select_()
}

func (m *Manager) select_() []*rule.Option {
	rrs := m.rules.Values()
	return funcs.Transform(rrs, func(rr *rule.Rule) *rule.Option {
		return rr.Option(Protocol2Name)
	})
}

func (m *Manager) Enable(id string, enable bool) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	rr, exists := m.rules.Search(id)
	if !exists {
		return false
	}
	rr.SetEnable(enable)
	return true
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
