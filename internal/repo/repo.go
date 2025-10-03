package repo

import (
	"YH-FireWall/internal/rule"
	"encoding/json"
	"os"
	"sort"
	"sync"
)

var (
	filepath = ".config.json"
	isModify = false
	rules    = make(map[uint32]*rule.Config) // Id → Rule 	// 控制 JSON 写入顺序
	mutex    sync.RWMutex
)

// Start 读取 JSON 文件到内存，并初始化 Index
func Start() error {
	mutex.Lock()
	defer mutex.Unlock()

	data, err := os.ReadFile(filepath)
	if err != nil {
		// 如果文件不存在，可以认为初始化为空
		if !os.IsNotExist(err) {
			return err
		}
		rules = make(map[uint32]*rule.Config)
		return nil
	}

	var rls []*rule.Config
	if err = json.Unmarshal(data, &rls); err != nil {
		return err
	}

	rules = make(map[uint32]*rule.Config)
	for _, r := range rls {
		rules[r.Id] = r
	}
	return nil
}

// Close 将内存规则写回 JSON 文件
func Close() error {
	mutex.Lock()
	defer mutex.Unlock()

	if !isModify {
		return nil
	}
	// 根据 Index 排序写回文件
	var rls []*rule.Config
	for _, r := range rules {
		rls = append(rls, r)
	}
	sort.Slice(rls, func(i, j int) bool {
		return rls[i].Priority < rls[j].Priority
	})

	data, err := json.MarshalIndent(rls, "", "    ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)
}
