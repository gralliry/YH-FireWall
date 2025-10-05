package group

import (
	"YH-FireWall/internal/core/rule"
	"fmt"
	"strings"
)

type Config struct {
	Name    string        `json:"name"`
	Comment string        `json:"comment"`
	Qnum    uint16        `json:"qnum"`
	Enable  bool          `json:"enable"`
	Rules   []rule.Config `json:"rules"`
}

func (c *Config) String() string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("- Group %s: %s %t\n", c.Name, c.Comment, c.Enable))

	for _, r := range c.Rules {
		sb.WriteString(fmt.Sprintf("\t- Rule %s:\n", r.Name))
		for _, s := range r.StringList() {
			sb.WriteString("\t\t" + s + "\n")
		}
	}

	return sb.String()
}
