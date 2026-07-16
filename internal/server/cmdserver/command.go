package cmdserver

import (
	"YH-FireWall/internal/model/rule"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func str(ss ...string) string {
	return strings.Join(ss, "\n")
}

type Handler interface {
	Version() string
	//
	AppendRule(ro *rule.Option) (string, error)
	UpdateRule(id string, ro *rule.Option) error
	DeleteRule(id string) error
	SearchRule(id string) *rule.Option
	ListRules() []*rule.Option
	EnableRule(id string, enable bool) bool
	//
	GetConfig() string
}

func newCmd(handler Handler) *cobra.Command {
	cmdRoot := &cobra.Command{
		Use:   "yfw",
		Short: "YH Firewall CLI tool",
		Long:  "YH Firewall is a powerful firewall service that can be controlled via command line interface.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("Usage: yfw <command> [<args>]")
		},
	}
	// rule
	cmdRule := &cobra.Command{
		Use:     "rule",
		Aliases: []string{"r"},
		Short:   "Manage firewall rules",
		Long:    "Manage firewall rules including listing, adding, removing, and modifying rules.",
		Example: str(
			"yfw rule list",
			`yfw rule append '{"action":"allow","protocol":"tcp","port":80}'`,
			"yfw rule remove RULE_ID",
			`yfw rule change RULE_ID '{"action":"deny","protocol":"tcp","port":80}'`,
			"yfw rule enable RULE_ID",
			"yfw rule disable RULE_ID",
		),
	}
	// rule list
	cmdRuleList := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "List all firewall rules or a specific rule",
		Long:    "List all firewall rules if no argument provided, or display a specific rule by ID.",
		Example: str(
			"yfw rule list",
			"yfw rule list RULE_ID",
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				rules := handler.ListRules()
				for _, r := range rules {
					cmd.Print(r.String())
				}
			} else {
				r := handler.SearchRule(args[0])
				if r == nil {
					return fmt.Errorf("No such rule")
				}
				cmd.Print(r.String())
			}
			return nil
		},
	}
	// rule append
	cmdRuleAppend := &cobra.Command{
		Use:     "append",
		Aliases: []string{"a", "add"},
		Short:   "Add a new firewall rule",
		Long:    "Add a new firewall rule with the provided configuration in JSON format.",
		Example: str(
			`yfw rule append '{"action":"allow","protocol":"tcp","port":80}'`,
			`yfw rule a '{"action":"deny","protocol":"udp","port":53}'`,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("Usage: append {config}")
			}
			var ro rule.Option
			if err := json.Unmarshal([]byte(args[0]), &ro); err != nil {
				return fmt.Errorf("invalid config: %w", err)
			}
			if _, err := handler.AppendRule(&ro); err != nil {
				return err
			}
			cmd.Println("Rule Appended")
			return nil
		},
	}
	// rule remove
	cmdRuleRemove := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"r", "rm", "del", "delete"},
		Short:   "Remove a firewall rule",
		Long:    "Remove a firewall rule by its ID.",
		Example: str(
			"yfw rule remove RULE_ID",
			"yfw rule rm RULE_ID",
			"yfw rule del RULE_ID",
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("Usage: remove {id}")
			}
			if err := handler.DeleteRule(args[0]); err != nil {
				return err
			}
			cmd.Println("Rule Removed")
			return nil
		},
	}
	// rule change
	cmdRuleChange := &cobra.Command{
		Use:     "change",
		Aliases: []string{"c", "modify", "update"},
		Short:   "Modify an existing firewall rule",
		Long:    "UpdateByPush an existing firewall rule with new configuration by providing the rule ID and new JSON configuration.",
		Example: str(
			`yfw rule change RULE_ID '{"action":"deny","protocol":"tcp","port":80}'`,
			`yfw rule update RULE_ID '{"action":"allow","protocol":"udp","port":53}'`,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 2 {
				return fmt.Errorf("Usage: change {id} {config}")
			}
			var ro rule.Option
			if err := json.Unmarshal([]byte(args[1]), &ro); err != nil {
				return fmt.Errorf("invalid config: %w", err)
			}
			if err := handler.UpdateRule(args[0], &ro); err != nil {
				return err
			}
			cmd.Println("Rule Updated")
			return nil
		},
	}
	// rule enable
	cmdRuleEnable := &cobra.Command{
		Use:     "enable",
		Aliases: []string{"e", "en"},
		Short:   "Enable a firewall rule",
		Long:    "Enable a previously disabled firewall rule by its ID.",
		Example: str(
			"yfw rule enable RULE_ID",
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("Usage: enable {id}")
			}
			if !handler.EnableRule(args[0], true) {
				return fmt.Errorf("No such rule")
			}
			cmd.Println("ok")
			return nil
		},
	}
	// rule disable
	cmdRuleDisable := &cobra.Command{
		Use:     "disable",
		Aliases: []string{"d", "dis"},
		Short:   "Disable a firewall rule",
		Long:    "Disable an active firewall rule by its ID.",
		Example: str(
			"yfw rule disable RULE_ID",
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("Usage: disable {id}")
			}
			if !handler.EnableRule(args[0], false) {
				return fmt.Errorf("No such rule")
			}
			cmd.Println("ok")
			return nil
		},
	}
	// config
	cmdConfig := &cobra.Command{
		Use:     "config",
		Aliases: []string{"c", "cfg"},
		Short:   "Get current configuration",
		Long:    "Display the current configuration of the YH Firewall service.",
		Example: str("yfw config"),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println(handler.GetConfig())
			return nil
		},
	}
	// version
	cmdVersion := &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Display version information",
		Long:    "Show the current version of the YH Firewall service.",
		Example: str(
			"yfw version",
			"yfw v",
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println(handler.Version())
			return nil
		},
	}
	cmdRule.AddCommand(
		cmdRuleList,
		cmdRuleAppend,
		cmdRuleRemove,
		cmdRuleChange,
		cmdRuleEnable,
		cmdRuleDisable,
	)
	cmdRoot.AddCommand(
		cmdRule,
		cmdConfig,
		cmdVersion,
	)
	return cmdRoot
}
