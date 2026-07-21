package rtable

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"YH-FireWall/internal/constant/itfdev"
	"YH-FireWall/internal/constant/protocol"
	"YH-FireWall/internal/model/flow"
	"YH-FireWall/internal/model/rule"
	"YH-FireWall/internal/pkg/cfile"
	"YH-FireWall/internal/pkg/funcs"
	"YH-FireWall/internal/pkg/indexedmap"
)

type Manager struct {
	mutex  sync.RWMutex
	config Config

	rules *indexedmap.IndexedMap[string, *rule.Rule]
	file  *cfile.CacheFile

	logger *slog.Logger
}

func New(config Config, logger *slog.Logger) (m *Manager, err error) {
	file, err := cfile.Open(config.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open rule file: %w", err)
	}
	defer func() {
		if err != nil {
			file.Close()
		}
	}()
	buf := file.Read()
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
	rules := indexedmap.New[string](func(a, b *rule.Rule) int {
		return a.Compare(b)
	})
	// 加载配置
	for _, rd := range rds {
		if rd.ID == "" {
			logger.Warn("rtable: skip rule", slog.String("error", "rule ID is empty"))
			continue
		}
		rr, err := rule.Parse(rd, itfdev.Name2Index, protocol.Name2Protocol)
		if err != nil {
			logger.Warn("rtable: skip rule", slog.String("id", rd.ID), slog.String("error", err.Error()))
			continue
		}
		rules.Insert(rd.ID, rr)
	}
	return &Manager{
		config: config,
		rules:  rules,
		file:   file,

		logger: logger,
	}, nil
}

func (m *Manager) Close() error {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.file.Close()
	return nil
}

func (m *Manager) list() []*rule.Data {
	rrs := m.rules.Values()
	return funcs.Transform(rrs, func(rr *rule.Rule) *rule.Data {
		return rr.Data(itfdev.Index2Name, protocol.Protocol2Name)
	})
}

func (m *Manager) save() {
	rds := m.list()
	buf, err := json.Marshal(rds)
	if err != nil {
		m.logger.Error("rtable: failed to marshal rules", slog.String("error", err.Error()))
		return
	}
	if err := m.file.Write(buf); err != nil {
		m.logger.Error("rtable: failed to persist rules", slog.String("error", err.Error()))
		return
	}
}

func (m *Manager) Create(ro *rule.Option) (string, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	rr, err := rule.New(ro, itfdev.Name2Index, protocol.Name2Protocol)
	if err != nil {
		return "", err
	}
	// 插入
	m.rules.Insert(rr.ID(), rr)
	// 持久化
	m.save()
	return rr.ID(), nil
}

func (m *Manager) Update(id string, ro *rule.Option) error {
	m.mutex.Lock()
	// 移除
	rr, exists := m.rules.Delete(id)
	if !exists {
		return fmt.Errorf("rule %s not exists", id)
	}
	// 更新
	rr, err := rr.Update(ro, itfdev.Name2Index, protocol.Name2Protocol)
	// 再添加
	m.rules.Insert(id, rr)
	//
	m.mutex.Unlock()
	// 错误返回
	if err != nil {
		return err
	}
	// 持久化
	m.save()
	return nil
}

func (m *Manager) Delete(id string) error {
	m.mutex.Lock()
	// 获取
	_, ok := m.rules.Delete(id)
	m.mutex.Unlock()
	if !ok {
		return fmt.Errorf("rule %s not exists", id)
	}
	// 持久化
	m.save()
	return nil
}

func (m *Manager) Search(id string) *rule.Data {
	m.mutex.RLock()
	rr, exists := m.rules.Search(id)
	m.mutex.RUnlock()
	if !exists {
		return nil
	}
	return rr.Data(itfdev.Index2Name, protocol.Protocol2Name)
}

func (m *Manager) List() []*rule.Data {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.list()
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
