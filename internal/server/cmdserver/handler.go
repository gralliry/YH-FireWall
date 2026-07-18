package cmdserver

import (
	"YH-FireWall/internal/model/rule"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

type HandleFunc = func(cmd *cobra.Command, args []string) error

func handleRuleList(handler Handler) HandleFunc {
	return func(cmd *cobra.Command, args []string) error {
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
	}
}

func handleRuleAppend(handler Handler) HandleFunc {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Usage: append {config}")
		}
		var ro rule.Option
		if err := json.Unmarshal([]byte(args[0]), &ro); err != nil {
			return fmt.Errorf("invalid config: %w", err)
		}
		if _, err := handler.CreateRule(&ro); err != nil {
			return err
		}
		cmd.Println("Rule Appended")
		return nil
	}
}

func handleRuleRemove(handler Handler) HandleFunc {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Usage: remove {id}")
		}
		if err := handler.DeleteRule(args[0]); err != nil {
			return err
		}
		cmd.Println("Rule Removed")
		return nil
	}
}

func handleRuleChange(handler Handler) HandleFunc {
	return func(cmd *cobra.Command, args []string) error {
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
	}
}

func handleRuleEnable(handler Handler) HandleFunc {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Usage: enable {id}")
		}
		if err := handler.EnableRule(args[0], true); err != nil {
			return fmt.Errorf("No such rule: %w", err)
		}
		cmd.Println("ok")
		return nil
	}
}

func handleRuleDisable(handler Handler) HandleFunc {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return fmt.Errorf("Usage: disable {id}")
		}
		if err := handler.EnableRule(args[0], false); err != nil {
			return fmt.Errorf("No such rule: %w", err)
		}
		cmd.Println("ok")
		return nil
	}
}

func handleRuleSet(handler Handler) HandleFunc {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) != 3 {
			return fmt.Errorf("Usage: set {id} {key} {value}")
		}
		ro, err := rule.Set(args[1], args[2])
		if err != nil {
			return err
		}
		if err := handler.UpdateRule(args[0], ro); err != nil {
			return err
		}
		cmd.Println("ok")
		return nil
	}
}

func handleConfigGet(handler Handler) HandleFunc {
	return func(cmd *cobra.Command, args []string) error {
		cmd.Printf("Config File at %s\n", handler.GetConfigPath())
		cmd.Println(handler.GetConfig())
		return nil
	}
}

func handleVersion(handler Handler) HandleFunc {
	return func(cmd *cobra.Command, args []string) error {
		cmd.Println(handler.Version())
		return nil
	}
}

func handleInterfaceList(handler Handler) HandleFunc {
	return func(cmd *cobra.Command, args []string) error {
		ifaces, err := handler.ListInterfaces()
		if err != nil {
			return err
		}
		for _, i := range ifaces {
			cmd.Printf("%-4d %-8s %-18s %-5d %v\n", i.Index, i.Name, i.MAC, i.MTU, i.Flags)
		}
		return nil
	}
}

func handleProtocolList(handler Handler) HandleFunc {
	return func(cmd *cobra.Command, args []string) error {
		for _, name := range handler.ListProtocols() {
			cmd.Println(name)
		}
		return nil
	}
}
