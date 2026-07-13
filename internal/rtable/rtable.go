package rtable

import (
	"YH-FireWall/internal/model/flow"
	"YH-FireWall/internal/pkg/flock"
	"YH-FireWall/internal/pkg/skiplist"
	"YH-FireWall/internal/rule"
	"encoding/json"
	"fmt"
	"log"
	"sync"
)

var (
	rules *skiplist.SkipList[string, *rule.Rule]

	mutex  sync.RWMutex
	lf     *flock.LockedFile
	config Config
)

func init() {
	// 设置默认规则
	rules = skiplist.New[string](func(a, b *rule.Rule) int {
		return a.Priority() - b.Priority()
	})
}

func Load(config_ Config) (err error) {
	mutex.Lock()
	defer mutex.Unlock()
	//
	config = config_
	//
	lf, err = flock.Open(config.Path)
	if err != nil {
		return fmt.Errorf("failed to open rule file: %w", err)
	}
	buf, err := lf.Read()
	if err != nil {
		return fmt.Errorf("failed to read rule file: %w", err)
	}
	ruleConfigs := make([]rule.Info, 0)
	if len(buf) > 0 {
		if err = json.Unmarshal(buf, &ruleConfigs); err != nil {
			return fmt.Errorf("failed to decode rule file: %w", err)
		}
	}
	// 加载配置
	for _, rc := range ruleConfigs {
		// 匹配
		if _, exists := rules.Search(rc.Id); exists {
			log.Printf("rule %s exists", rc.Id)
			continue
		}
		// 解析
		if rr, err := rule.New(&rc); err != nil {
			log.Printf("failed to parse rule %s: %v", rc.Id, err)
			continue
		} else {
			// 如果都没有，就添加
			rules.Insert(rc.Id, rr)
		}
	}
	return nil
}

func Close() error {
	mutex.Lock()
	defer mutex.Unlock()
	buf, err := json.MarshalIndent(getAll(), "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode rules: %w", err)
	}
	if err := lf.Write(buf); err != nil {
		return err
	}
	return lf.Close()
}

// Append 添加或更新规则
func Append(ro *rule.Option) (string, error) {
	// uuid不可能重复
	rc := ro.Build()
	//
	rr, err := rule.New(rc)
	if err != nil {
		return "", err
	}
	//
	mutex.Lock()
	defer mutex.Unlock()
	// 如果都没有，就添加
	rules.Insert(rr.Id(), rr)
	//
	return rr.Id(), nil
}

// Update 更新规则
func Update(id string, ro *rule.Option) error {
	mutex.Lock()
	defer mutex.Unlock()

	rr, exists := rules.Search(id)
	if !exists {
		return fmt.Errorf("rule %s not exists", id)
	}
	if err := rr.Update(*ro); err != nil {
		return fmt.Errorf("update rule %s error: %v", id, err)
	}
	return nil
}

// Delete 删除规则
func Delete(id string) error {
	mutex.Lock()
	defer mutex.Unlock()
	// 获取
	ok := rules.Delete(id)
	if !ok {
		return fmt.Errorf("rule %s not exists", id)
	}
	return nil
}

// Match 匹配：按优先级从高到低
func Match(flow *flow.Flow) bool {
	mutex.RLock()
	defer mutex.RUnlock()
	// 匹配
	_, rr, exist := rules.First(func(r *rule.Rule) bool {
		return r.Match(flow)
	})
	if !exist {
		return config.DefaultAccept
	}
	return rr.Accept()
}

func Search(rid string) *rule.Info {
	mutex.RLock()
	defer mutex.RUnlock()
	if rr, exists := rules.Search(rid); !exists {
		return nil
	} else {
		return rr.Info()
	}
}

func SearchAll() []rule.Info {
	mutex.RLock()
	defer mutex.RUnlock()
	return getAll()
}

func getAll() []rule.Info {
	ruleConfigs := make([]rule.Info, rules.Len())
	rules.Range(func(key string, rr *rule.Rule) {
		ruleConfigs = append(ruleConfigs, *rr.Info())
	})
	return ruleConfigs
}

func Enable(id string, enable bool) bool {
	mutex.Lock()
	defer mutex.Unlock()
	rr, exists := rules.Search(id)
	if !exists {
		return false
	}
	rr.SetEnable(enable)
	return true
}
