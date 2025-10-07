package core

import (
	"YH-FireWall/internal/core/manager"
	"YH-FireWall/internal/core/rule"
)

// AppendRule 添加或更新规则
func AppendRule(ro *rule.Option) error {
	return manager.AppendRule(ro)
}

// UpdateRule 更新规则
func UpdateRule(id string, ro *rule.Option) error {
	return manager.UpdateRule(id, ro)
}

// DeleteRule 删除规则
func DeleteRule(id string) error {
	return manager.DeleteRule(id)
}

func GetRule(id string) *rule.Config {
	return manager.GetRule(id)
}

func GetRules() []rule.Config {
	return manager.GetRules()
}

func EnableRule(id string, enable bool) bool {
	return manager.EnableRule(id, enable)
}

func EnableGroup(group string, enable bool) bool {
	return manager.EnableGroup(group, enable)
}
