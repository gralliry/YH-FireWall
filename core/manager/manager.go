package manager

import (
	"YH-FireWall/core/rule"
	"sort"
	"sync"
)

var (
	mutex sync.RWMutex
)

func Init(rules []rule.Config) error {
	mutex.Lock()
	defer mutex.Unlock()
	// --- 处理规则 ---
	ruleList = make([]*rule.Rule, 0)
	ruleMap = make(map[string]*rule.Rule)
	// 加载配置
	var err error
	var rr *rule.Rule
	for _, rc := range rules {
		// 解析规则
		rr, err = rule.Parse(rc)
		if err != nil {
			continue
		}
		// 检测ID是否被使用
		if _, exists := ruleMap[rr.Id()]; exists {
			continue
		}
		// 写入
		ruleMap[rr.Id()] = rr
		ruleList = append(ruleList, rr)
	}
	// 排序
	sort.SliceStable(ruleList, func(i, j int) bool {
		return ruleList[i].Priority() < ruleList[j].Priority()
	})
	// 标记为非脏数据
	ruleIsListDirty.Store(false)
	return nil
}
