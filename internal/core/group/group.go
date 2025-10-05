package group

import (
	"YH-FireWall/internal/core/packet"
	"YH-FireWall/internal/core/rule"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
)

type Group struct {
	name    string
	comment string
	qnum    uint16
	enable  bool
	// 索引
	ruleList []*rule.Rule
	ruleMap  map[string]*rule.Rule
	//
	mutex sync.RWMutex
	dirty atomic.Bool
}

func Load(config *Config) (*Group, error) {
	if config == nil {
		return nil, fmt.Errorf("group config is nil")
	}
	if config.Name == "" {
		return nil, fmt.Errorf("group name is empty")
	}
	g := &Group{
		name:     config.Name,
		comment:  config.Comment,
		qnum:     config.Qnum,
		enable:   config.Enable,
		ruleList: make([]*rule.Rule, 0, len(config.Rules)),
		ruleMap:  make(map[string]*rule.Rule),
	}
	// 脏
	g.dirty.Store(false)
	// 不再检查是否重复，如果重复，直接覆盖
	for _, rc := range config.Rules {
		r, err := rule.Parse(&rc)
		// 忽略错误
		if err != nil {
			continue
		}
		g.ruleMap[r.Name()] = r
	}
	for _, r := range g.ruleMap {
		g.ruleList = append(g.ruleList, r)
	}
	return g, nil
}

// Append 添加或更新规则
func (g *Group) Append(rc *rule.Config) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if _, exists := g.ruleMap[rc.Name]; exists {
		return fmt.Errorf("rule %s exists", rc.Name)
	}
	r, err := rule.Parse(rc)
	if err != nil {
		return fmt.Errorf("rule %s parse error: %v", rc.Name, err)
	}
	// 如果都没有，就添加
	g.ruleMap[rc.Name] = r
	g.ruleList = append(g.ruleList, r)
	return nil
}

func (g *Group) Update(rc *rule.Config) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if _, exists := g.ruleMap[rc.Name]; !exists {
		return fmt.Errorf("rule %s not exists", rc.Name)
	}
	r, err := rule.Parse(rc)
	if err != nil {
		return fmt.Errorf("rule %s parse error: %v", rc.Name, err)
	}
	g.ruleMap[rc.Name] = r
	for i, rr := range g.ruleList {
		if rr.Name() == rc.Name {
			g.ruleList[i] = r
			return nil
		}
	}
	return nil
}

// Delete 删除规则
func (g *Group) Delete(name string) error {
	g.mutex.Lock()
	defer g.mutex.Unlock()
	if _, exists := g.ruleMap[name]; !exists {
		return fmt.Errorf("rule %s not exists", name)
	}
	delete(g.ruleMap, name)
	// 重新构造
	g.ruleList = g.ruleList[:0]
	for _, r := range g.ruleList {
		if r.Name() != name {
			g.ruleList = append(g.ruleList, r)
		}
	}
	return nil
}

// Match 匹配：按优先级从高到低
func (g *Group) Match(p *packet.Packet) (bool, bool) {
	if g.dirty.Load() {
		g.mutex.Lock()
		if g.dirty.Load() {
			sort.SliceStable(g.ruleList, func(i, j int) bool {
				return g.ruleList[i].Priority() < g.ruleList[j].Priority()
			})
		}
		g.mutex.Unlock()
		g.dirty.Store(false)
	}

	g.mutex.Lock()
	defer g.mutex.Unlock()
	//
	if !g.enable {
		return false, false
	}
	// 匹配
	for _, r := range g.ruleList {
		if r.Match(p) {
			return true, r.Accept()
		}
	}
	return false, false
}

func (g *Group) Unparse() *Config {
	g.mutex.RLock()
	defer g.mutex.RUnlock()
	c := &Config{
		Name:    g.name,
		Comment: g.comment,
		Qnum:    g.qnum,
		Enable:  g.enable,
		Rules:   make([]rule.Config, 0, len(g.ruleList)),
	}
	for _, r := range g.ruleList {
		c.Rules = append(c.Rules, *r.Unparse())
	}
	return c
}

func (g *Group) Name() string {
	return g.name
}
