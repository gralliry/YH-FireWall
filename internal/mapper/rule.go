package mapper

import (
	"YH-FireWall/internal/repo"
	"YH-FireWall/internal/rule"
)

func GetAllRules() ([]rule.Config, error) {
	return repo.GetAllRules()
}

func UpdateRule(rule rule.Config) error {
	return repo.UpdateRule(rule)
}

func AppendRule(rule rule.Config) error {
	return repo.AppendRule(rule)
}

func DeleteRule(id uint32) error {
	return repo.DeleteRule(id)
}
