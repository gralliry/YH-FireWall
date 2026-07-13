package cmdserver

import (
	"YH-FireWall/internal/rule"
	"encoding/json"

	"github.com/spf13/cobra"
)

type Handler interface {
	Version() string

	AppendRule(ro *rule.Option) (string, error)
	UpdateRule(id string, ro *rule.Option) error
	DeleteRule(id string) error
	SearchRule(id string) *rule.Info
	SearchRules() []rule.Info
	EnableRule(id string, enable bool) bool

	GetConfig() string
}

func newCommand(handler Handler) *cobra.Command {
	cmdRoot := &cobra.Command{
		Use:   "yfw",
		Short: "YH Firewall CLI tool",
		Long:  "YH Firewall is a powerful firewall service that can be controlled via command line interface.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.PrintErrln("Usage: yfw <command> [<args>]")
		},
	}
	//  start
	cmdStart := &cobra.Command{
		Use:   "start",
		Short: "Start the YH Firewall service",
		Long:  "Start the YH Firewall service in the background.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("yfw has started")
		},
		Example: `
yfw start
`,
	}
	cmdRoot.AddCommand(cmdStart)
	//  status
	cmdStatus := &cobra.Command{
		Use:   "status",
		Short: "Check the status of YH Firewall service",
		Long:  "Display the current running status of the YH Firewall service.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("ok")
		},
		Example: `
yfw status
`,
	}
	cmdRoot.AddCommand(cmdStatus)
	// rule
	cmdRule := &cobra.Command{
		Use:     "rule",
		Aliases: []string{"r"},
		Short:   "Manage firewall rules",
		Long:    "Manage firewall rules including listing, adding, removing, and modifying rules.",
		Example: `
yfw rule list
yfw rule append '{"action":"allow","protocol":"tcp","port":80}'
yfw rule remove RULE_ID
yfw rule change RULE_ID '{"action":"deny","protocol":"tcp","port":80}'
yfw rule enable RULE_ID
yfw rule disable RULE_ID
`,
	}
	cmdRoot.AddCommand(cmdRule)
	// rule list
	cmdRuleList := &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls", "l"},
		Short:   "List all firewall rules or a specific rule",
		Long:    "List all firewall rules if no argument provided, or display a specific rule by ID.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				rules := handler.SearchRules()
				for _, r := range rules {
					cmd.Print(r.String())
				}
			} else {
				r := handler.SearchRule(args[0])
				if r == nil {
					cmd.PrintErrln("No such rule")
				} else {
					cmd.Print(r.String())
				}
			}
		},
		Example: `
yfw rule list
yfw rule list RULE_ID
`,
	}
	cmdRule.AddCommand(cmdRuleList)
	// rule append
	cmdRuleAppend := &cobra.Command{
		Use:     "append",
		Aliases: []string{"a", "add"},
		Short:   "Add a new firewall rule",
		Long:    "Add a new firewall rule with the provided configuration in JSON format.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.PrintErrln("Usage: append {config}")
				return
			}
			var ro rule.Option
			if err := json.Unmarshal([]byte(args[0]), &ro); err != nil {
				cmd.PrintErrln(err)
				return
			}
			if _, err := handler.AppendRule(&ro); err != nil {
				cmd.PrintErrln(err)
				return
			}
			cmd.Println("Rule Appended")
		},
		Example: `
yfw rule append '{"action":"allow","protocol":"tcp","port":80}'
yfw rule a '{"action":"deny","protocol":"udp","port":53}'
`,
	}
	cmdRule.AddCommand(cmdRuleAppend)
	// rule remove
	cmdRuleRemove := &cobra.Command{
		Use:     "remove",
		Aliases: []string{"r", "rm", "del", "delete"},
		Short:   "Remove a firewall rule",
		Long:    "Remove a firewall rule by its ID.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.PrintErrln("Usage: remove {id}")
				return
			}
			if err := handler.DeleteRule(args[0]); err != nil {
				cmd.PrintErrln(err)
				return
			}
			cmd.Println("Rule Removed")
		},
		Example: `
yfw rule remove RULE_ID
yfw rule rm RULE_ID
yfw rule del RULE_ID
`,
	}
	cmdRule.AddCommand(cmdRuleRemove)
	//
	cmdRuleChange := &cobra.Command{
		Use:     "change",
		Aliases: []string{"c", "modify", "update"},
		Short:   "Modify an existing firewall rule",
		Long:    "UpdateByPush an existing firewall rule with new configuration by providing the rule ID and new JSON configuration.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) != 2 {
				cmd.PrintErrln("Usage: change {id} {config}")
				return
			}
			var ro rule.Option
			if err := json.Unmarshal([]byte(args[1]), &ro); err != nil {
				cmd.PrintErrln(err)
				return
			}
			if err := handler.UpdateRule(args[0], &ro); err != nil {
				cmd.PrintErrln(err)
				return
			}
			cmd.Println("Rule Updated")
		},
		Example: `
yfw rule change RULE_ID '{"action":"deny","protocol":"tcp","port":80}'
yfw rule update RULE_ID '{"action":"allow","protocol":"udp","port":53}'
`,
	}
	cmdRule.AddCommand(cmdRuleChange)
	//
	cmdRuleEnable := &cobra.Command{
		Use:     "enable",
		Aliases: []string{"e", "en"},
		Short:   "Enable a firewall rule",
		Long:    "Enable a previously disabled firewall rule by its ID.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.PrintErrln("Usage: enable {id}")
				return
			}
			if !handler.EnableRule(args[0], true) {
				cmd.PrintErrln("No such rule")
				return
			}
			cmd.Println("ok")
		},
		Example: `
yfw rule enable RULE_ID
`,
	}
	cmdRule.AddCommand(cmdRuleEnable)
	//
	cmdRuleDisable := &cobra.Command{
		Use:     "disable",
		Aliases: []string{"d", "dis"},
		Short:   "Disable a firewall rule",
		Long:    "Disable an active firewall rule by its ID.",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				cmd.PrintErrln("Usage: disable {id}")
				return
			}
			if !handler.EnableRule(args[0], false) {
				cmd.PrintErrln("No such rule")
				return
			}
			cmd.Println("ok")
		},
		Example: `
yfw rule disable RULE_ID
`,
	}
	cmdRule.AddCommand(cmdRuleDisable)
	// config
	cmdConfig := &cobra.Command{
		Use:     "config",
		Aliases: []string{"c", "cfg"},
		Short:   "Get current configuration",
		Long:    "Display the current configuration of the YH Firewall service.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(handler.GetConfig())
		},
		Example: `
yfw config
`,
	}
	cmdRoot.AddCommand(cmdConfig)
	// version
	cmdVersion := &cobra.Command{
		Use:     "version",
		Aliases: []string{"v"},
		Short:   "Display version information",
		Long:    "Show the current version of the YH Firewall service.",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(handler.Version())
		},
		Example: `
yfw version
yfw v
`,
	}
	cmdRoot.AddCommand(cmdVersion)
	return cmdRoot
}
