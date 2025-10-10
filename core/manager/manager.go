package manager

import (
	"YH-FireWall/core/config"
	"YH-FireWall/core/rule"
	"sort"
	"sync"
	"time"
)

var (
	cfg   *config.Config
	mutex sync.RWMutex
)

func Init(c *config.Config) error {
	mutex.Lock()
	defer mutex.Unlock()
	cfg = c
	// --- 处理规则 ---
	// 清空
	ruleList = make([]*rule.Rule, 0)
	ruleMap = make(map[string]*rule.Rule)
	// 加载配置
	var err error
	var rr *rule.Rule
	for _, rc := range cfg.Rules {
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

func Sync(cfg *config.Config) error {
	mutex.Lock()
	defer mutex.Unlock()
	// 更新时间
	cfg.LastUpdateDate = time.Now().Format("2006-01-02 15:04:05")
	// 存储配置文件
	cfg.Rules = getRules()

	return nil
}
