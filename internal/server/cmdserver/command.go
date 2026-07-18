package cmdserver

import (
	"fmt"

	"github.com/spf13/cobra"

	"YH-FireWall/internal/model/rule"
)

type Handler interface {
	Version() string
	//
	CreateRule(ro *rule.Option) (string, error)
	UpdateRule(id string, ro *rule.Option) error
	DeleteRule(id string) error

	SearchRule(id string) *rule.Data
	ListRules() []*rule.Data

	EnableRule(id string, enable bool) error

	//
	GetConfig() string
	GetConfigPath() string
}

func newCmd(handler Handler) *cobra.Command {
	cmdRoot := &cobra.Command{
		Use:   "yfw",
		Short: "YH Firewall CLI tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("Usage: yfw <command> [<args>]")
		},
	}

	cmdRule := &cobra.Command{
		Use:     "rule",
		Aliases: []string{"r"},
		Short:   "Manage firewall rules",
	}
	cmdRule.AddCommand(
		&cobra.Command{
			Use:     "list",
			Aliases: []string{"ls", "l"},
			Short:   "List all firewall rules or a specific rule",
			RunE:    handlerRuleList(handler),
		},
		&cobra.Command{
			Use:     "append",
			Aliases: []string{"a", "add"},
			Short:   "Add a new firewall rule",
			RunE:    handlerRuleAppend(handler),
		},
		&cobra.Command{
			Use:     "remove",
			Aliases: []string{"r", "rm", "del", "delete"},
			Short:   "Remove a firewall rule",
			RunE:    handlerRuleRemove(handler),
		},
		&cobra.Command{
			Use:     "change",
			Aliases: []string{"c", "modify", "update"},
			Short:   "Modify an existing firewall rule",
			RunE:    handlerRuleChange(handler),
		},
		&cobra.Command{
			Use:   "set",
			Short: "Set a single field by key=value",
			Long: "Usage: set {id} {key} {value}\n" +
				"Keys (case-insensitive, lowercase): group, comment, accept, enable, priority," +
				" inDevs, outDevs, protocols, srcNets, dstNets, srcPorts, dstPorts",
			RunE: handlerRuleSet(handler),
		},
		&cobra.Command{
			Use:     "enable",
			Aliases: []string{"e", "en"},
			Short:   "Enable a firewall rule",
			RunE:    handlerRuleEnable(handler),
		},
		&cobra.Command{
			Use:     "disable",
			Aliases: []string{"d", "dis"},
			Short:   "Disable a firewall rule",
			RunE:    handlerRuleDisable(handler),
		},
	)

	cmdRoot.AddCommand(
		cmdRule,
		&cobra.Command{
			Use:     "config",
			Aliases: []string{"c", "cfg"},
			Short:   "Get current configuration",
			RunE:    handlerConfigGet(handler),
		},
		&cobra.Command{
			Use:     "version",
			Aliases: []string{"v"},
			Short:   "Display version information",
			RunE:    handlerVersion(handler),
		},
	)
	return cmdRoot
}
