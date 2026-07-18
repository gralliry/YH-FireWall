package cmdserver

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"YH-FireWall/internal/model/itf"
	"YH-FireWall/internal/model/rule"
)

func str(ss ...string) string {
	return strings.Join(ss, "\n")
}

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
	//
	ListInterfaces() ([]itf.Itf, error)
	ListProtocols() []string
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
		Example: str(
			"yfw rule list",
			`yfw rule append '{"accept":true,"srcNets":"10.0.0.0/8"}'`,
			"yfw rule remove RULE_ID",
			`yfw rule change RULE_ID '{"accept":false}'`,
			"yfw rule enable RULE_ID",
			"yfw rule disable RULE_ID",
		),
	}
	cmdRuleList := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "List all firewall rules or a specific rule",
		Example: str(
			"yfw rule list",
			"yfw rule list RULE_ID",
		),
		RunE: handleRuleList(handler),
	}
	cmdRuleAppend := &cobra.Command{
		Use:     "append",
		Aliases: []string{"a", "add"},
		Short:   "Add a new firewall rule",
		Example: str(
			`yfw rule append '{"accept":true,"srcNets":"10.0.0.0/8"}'`,
			`yfw rule a '{"accept":false,"protocols":"tcp"}'`,
		),
		RunE: handleRuleAppend(handler),
	}
	cmdRuleRemove := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"r", "rm", "del", "delete"},
		Short:   "Remove a firewall rule",
		Example: str(
			"yfw rule remove RULE_ID",
			"yfw rule rm RULE_ID",
		),
		RunE: handleRuleRemove(handler),
	}
	cmdRuleChange := &cobra.Command{
		Use:     "change",
		Aliases: []string{"c", "modify", "update"},
		Short:   "Modify an existing firewall rule",
		Example: str(
			`yfw rule change RULE_ID '{"accept":false}'`,
			`yfw rule update RULE_ID '{"protocols":"udp"}'`,
		),
		RunE: handleRuleChange(handler),
	}
	cmdRuleSet := &cobra.Command{
		Use:   "set",
		Short: "Set a single field by key=value",
		Long: "Keys (case-insensitive, lowercase): group, comment, accept, enable, priority," +
			" inDevs, outDevs, protocols, srcNets, dstNets, srcPorts, dstPorts",
		Example: str(
			"yfw rule set RULE_ID accept false",
			"yfw rule set RULE_ID srcNets 10.0.0.0/8",
		),
		RunE: handleRuleSet(handler),
	}
	cmdRuleEnable := &cobra.Command{
		Use:     "enable",
		Aliases: []string{"e", "en"},
		Short:   "Enable a firewall rule",
		Example: str("yfw rule enable RULE_ID"),
		RunE:    handleRuleEnable(handler),
	}
	cmdRuleDisable := &cobra.Command{
		Use:     "disable",
		Aliases: []string{"d", "dis"},
		Short:   "Disable a firewall rule",
		Example: str("yfw rule disable RULE_ID"),
		RunE:    handleRuleDisable(handler),
	}
	cmdRule.AddCommand(
		cmdRuleList,
		cmdRuleAppend,
		cmdRuleRemove,
		cmdRuleChange,
		cmdRuleSet,
		cmdRuleEnable,
		cmdRuleDisable,
	)

	cmdConfig := &cobra.Command{
		Use:     "config",
		Aliases: []string{"c", "cfg"},
		Short:   "Get current configuration",
		Example: str("yfw config"),
		RunE:    handleConfigGet(handler),
	}
	cmdVersion := &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Display version information",
		Example: str(
			"yfw version",
			"yfw v",
		),
		RunE: handleVersion(handler),
	}
	cmdInterfaces := &cobra.Command{
		Use:     "interfaces",
		Aliases: []string{"iface", "i"},
		Short:   "List available network interfaces",
		Example: str(
			"yfw interfaces",
			"yfw iface",
		),
		RunE: handleInterfaceList(handler),
	}
	cmdProtocols := &cobra.Command{
		Use:     "protocols",
		Aliases: []string{"proto", "p"},
		Short:   "List supported protocols",
		Example: str(
			"yfw protocols",
			"yfw proto",
		),
		RunE: handleProtocolList(handler),
	}
	cmdRoot.AddCommand(
		cmdRule,
		cmdConfig,
		cmdVersion,
		cmdInterfaces,
		cmdProtocols,
	)
	return cmdRoot
}
