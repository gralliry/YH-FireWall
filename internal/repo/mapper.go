package repo

import (
	"YH-FireWall/internal/rule"
	"errors"
	"fmt"
	"math/rand"
	"sort"
)

// GetAllRules 返回按 Index 排序的所有规则（拷贝）
func GetAllRules() ([]rule.Config, error) {
	mutex.RLock()
	defer mutex.RUnlock()

	rls := make([]*rule.Config, 0, len(rules))
	for _, r := range rules {
		rls = append(rls, r)
	}
	sort.Slice(rls, func(i, j int) bool {
		return rls[i].Protocol < rls[j].Protocol
	})

	res := make([]rule.Config, 0, len(rls))
	for _, r := range rls {
		res = append(res, *r)
	}
	return res, nil
}

// AppendRule 新增规则，Id 系统生成，Index 自动写在最后
func AppendRule(rule rule.Config) error {
	mutex.Lock()
	defer mutex.Unlock()

	// 生成随机 Id，直到不重复
	for i := 0; i < 10; i++ { // 尝试 10 次
		randID := rand.Uint32()
		if _, exists := rules[randID]; !exists {
			rule.Id = randID
			break
		}
		if i == 9 {
			return errors.New("无法生成唯一 Id")
		}
	}
	//
	isModify = true
	rules[rule.Id] = &rule
	return nil
}

// DeleteRule 删除规则
func DeleteRule(id uint32) error {
	mutex.Lock()
	defer mutex.Unlock()

	if _, exists := rules[id]; !exists {
		return fmt.Errorf("规则不存在，id=%d", id)
	}

	isModify = true
	delete(rules, id)
	return nil
}

// UpdateRule 更新规则（以 Id 为准）
func UpdateRule(rule rule.Config) error {
	mutex.Lock()
	defer mutex.Unlock()

	_, exists := rules[rule.Id]
	if !exists {
		return fmt.Errorf("规则不存在，id=%d", rule.Id)
	}
	isModify = true
	rules[rule.Id] = &rule
	return nil
}
