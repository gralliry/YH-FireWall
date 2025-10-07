package manager

import (
	"YH-FireWall/internal/core/packet"
	"YH-FireWall/internal/core/rule"
	"fmt"
	"sort"
	"sync"
	"sync/atomic"
)

var (
	// 索引
	ruleList = make([]*rule.Rule, 0)
	ruleMap  = make(map[string]*rule.Rule)
	//
	mutex sync.RWMutex
	dirty atomic.Bool
)

// RegisterRule 注册规则
func RegisterRule(rr *rule.Rule) error {
	mutex.Lock()
	defer mutex.Unlock()
	//
	if _, exists := ruleMap[rr.Id()]; exists {
		return fmt.Errorf("rule %s exists", rr.Id())
	}
	// 标记
	dirty.Store(true)
	// 如果都没有，就添加
	ruleMap[rr.Id()] = rr
	ruleList = append(ruleList, rr)
	return nil
}

// AppendRule 添加或更新规则
func AppendRule(ro *rule.Option) error {
	rc := ro.Default()
	//
	mutex.Lock()
	defer mutex.Unlock()
	// 注意性能，小心死循环
	for {
		// 检测ID是否被使用
		if _, exists := ruleMap[rc.Id]; !exists {
			break
		}
		// 刷新，避免ID冲突
		rc.Refresh()
	}
	//
	r, err := rule.Parse(*rc)
	if err != nil {
		return err
	}
	// 标记
	dirty.Store(true)
	// 如果都没有，就添加
	ruleMap[r.Id()] = r
	ruleList = append(ruleList, r)
	return nil
}

// UpdateRule 更新规则
func UpdateRule(id string, ro *rule.Option) error {
	mutex.RLock()
	defer mutex.RUnlock()
	//
	rr, exists := ruleMap[id]
	if !exists {
		return fmt.Errorf("rule %s not exists", id)
	}
	//
	if err := rr.Update(*ro); err != nil {
		return err
	}
	return nil
}

// DeleteRule 删除规则
func DeleteRule(id string) error {
	mutex.Lock()
	defer mutex.Unlock()
	//
	if _, exists := ruleMap[id]; !exists {
		return fmt.Errorf("rule %s not exists", id)
	}
	//
	dirty.Store(true)
	//
	delete(ruleMap, id)
	// 重新构造
	ruleList = ruleList[:0]
	for _, r := range ruleList {
		if r.Id() != id {
			ruleList = append(ruleList, r)
		}
	}
	return nil
}

// Match 匹配：按优先级从高到低
func Match(p *packet.Packet) (bool, bool) {
	//
	if dirty.Load() {
		mutex.Lock()
		if dirty.Load() {
			sort.SliceStable(ruleList, func(i, j int) bool {
				return ruleList[i].Priority() < ruleList[j].Priority()
			})
		}
		dirty.Store(false)
		mutex.Unlock()
	}
	//
	mutex.Lock()
	defer mutex.Unlock()
	// 匹配
	for _, r := range ruleList {
		if r.Match(p) {
			return true, r.Accept()
		}
	}
	return false, false
}

func GetRule(rid string) *rule.Config {
	mutex.RLock()
	defer mutex.RUnlock()
	rr, exists := ruleMap[rid]
	if !exists {
		return nil
	}
	return rr.Unparse()
}

func GetRules() []rule.Config {
	mutex.RLock()
	defer mutex.RUnlock()
	//
	rules := make([]rule.Config, 0)
	for _, r := range ruleList {
		rules = append(rules, *r.Unparse())
	}
	return rules
}

func EnableRule(id string, enable bool) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	rr, exists := ruleMap[id]
	if !exists {
		return false
	}
	rr.SetEnable(enable)
	return true
}

func EnableGroup(group string, enable bool) (effect bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	for _, rr := range ruleList {
		if rr.Group() == group {
			if !effect && rr.Enable() != enable {
				effect = true
			}
			rr.SetEnable(enable)
		}
	}
	return effect
}
