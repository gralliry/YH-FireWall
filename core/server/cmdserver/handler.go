package cmdserver

import (
	"YH-FireWall/core/rule"
	"encoding/json"
	"fmt"
	"strings"
)

type Handler interface {
	Version() string

	Stop() error

	AppendRule(ro *rule.Option) error
	UpdateRule(id string, ro *rule.Option) error
	DeleteRule(id string) error
	GetRule(id string) *rule.Config
	GetRules() []rule.Config
	EnableRule(id string, enable bool) bool

	GetConfig() (string, error)
}

var handler Handler

// Use 4 spaces instead of \t
const tips = `
Usage: %s <command> [<args>]
Commands:
    - start         Start the firewall service.               Example. yfw start
    - stop          Stop the firewall service.                Example. yfw stop
    - status        Check the status of the firewall service. Example. yfw status
    - r/rule
        - l/ls/list               Example. yfw rule list
            - {id}                Example. yfw rule list DSHUIC90
        - a/add/append {config}   Example. yfw rule add '{"tar_net":"0.0.0.0/0"}'
        - r/remove {id}           Example. yfw rule remove 12345678
        - c/change {id} {config}  Example. yfw rule change 12345678 '{"tar_net":"0.0.0.0/0"}'
        - e/enable {id}
        - d/disable {id}
    - c/cfg/config
    - h/help       Display this help message
    - v/version
`

func handleArgs(args []string) string {
	if len(args) == 0 {
		return fmt.Sprintf(tips, "No Any Param")
	}
	if len(args) == 1 {
		return fmt.Sprintf(tips, args[0])
	}
	switch args[1] {
	case "start":
		// yfw start
		return "How do you get there? ! It is impossible!"
	case "status":
		// yfw status
		return handleStatus(args[2:])
	case "stop":
		return handleStop(args[2:])
	case "r", "rule":
		// yfw r/rule
		return handleRule(args[2:])
	case "c", "cfg", "config":
		if data, err := handler.GetConfig(); err != nil {
			return err.Error()
		} else {
			return data
		}
	case "h", "help":
		// yfw h/help
		return fmt.Sprintf(tips, args[0])
	case "v", "version":
		// yfw h/help
		return handler.Version()
	default:
		return fmt.Sprintf("Unknown Command {%s}. Use help.", args[1])
	}
}

func handleStatus(_ []string) string {
	// yfw status
	return "ok"
}

// Stop 停止服务
func handleStop(_ []string) string {
	// yfw stop
	if err := handler.Stop(); err != nil {
		return "YFW has tried to stop but error occurred. Use status to check the status"
	}
	return "YFW Stopped"
}

// -----------------------------------------------------------------
func handleRule(args []string) string {
	if len(args) == 0 {
		return "Rule: Missing subcommand"
	}
	switch args[0] {
	case "l", "ls", "list":
		return handleRuleList(args[1:])
	case "a", "add", "append":
		return handleRuleAppend(args[1:])
	case "r", "remove":
		return handleRuleRemove(args[1:])
	case "c", "change":
		return handleRuleChange(args[1:])
	case "e", "enable":
		return handleRuleEnable(args[1:])
	case "d", "disable":
		return handleRuleDisable(args[1:])
	default:
		return "Unknown rule subcommand"
	}
}

func handleRuleList(args []string) string {
	if len(args) == 0 {
		rules := handler.GetRules()
		var sb strings.Builder
		for _, r := range rules {
			sb.WriteString(r.String())
		}
		return sb.String()
	} else {
		r := handler.GetRule(args[0])
		if r == nil {
			return "No such rule"
		}
		return r.String()
	}
}

func handleRuleAppend(args []string) string {
	if len(args) == 0 {
		return "Usage: append {config}"
	}
	var ro rule.Option
	if err := json.Unmarshal([]byte(args[0]), &ro); err != nil {
		return err.Error()
	}
	if err := handler.AppendRule(&ro); err != nil {
		return err.Error()
	}
	return "ok"
}

func handleRuleRemove(args []string) string {
	if len(args) == 0 {
		return "Usage: remove {id}"
	}
	if err := handler.DeleteRule(args[0]); err != nil {
		return err.Error()
	}
	return "ok"
}

func handleRuleChange(args []string) string {
	if len(args) != 2 {
		return "Usage: change {id} {config}"
	}
	var ro rule.Option
	if err := json.Unmarshal([]byte(args[1]), &ro); err != nil {
		return err.Error()
	}
	if err := handler.UpdateRule(args[0], &ro); err != nil {
		return err.Error()
	}
	return "ok"
}

func handleRuleEnable(args []string) string {
	if len(args) == 0 {
		return "Usage: enable {id}"
	}
	if handler.EnableRule(args[0], true) {
		return "ok"
	}
	return "No such rule"
}

func handleRuleDisable(args []string) string {
	if len(args) == 0 {
		return "Usage: disable {id}"
	}
	if handler.EnableRule(args[0], false) {
		return "ok"
	}
	return "No such rule"
}
