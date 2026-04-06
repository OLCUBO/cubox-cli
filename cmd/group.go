package cmd

import (
	"fmt"
	"strings"

	"github.com/OLCUBO/cubox-cli/internal/client"
	"github.com/OLCUBO/cubox-cli/internal/config"
	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage groups (folders)",
}

var groupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all groups",
	RunE:  runGroupList,
}

func init() {
	groupCmd.AddCommand(groupListCmd)
	rootCmd.AddCommand(groupCmd)
}

func runGroupList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c := client.New(cfg.BaseURL(), cfg.Token)
	groups, err := c.ListGroups()
	if err != nil {
		return err
	}

	if outputFormat == "text" {
		printGroupsText(groups)
		return nil
	}
	printJSON(groups)
	return nil
}

func printGroupsText(groups []client.Group) {
	for _, g := range groups {
		depth := strings.Count(g.NestedName, "/")
		indent := strings.Repeat("  ", depth)
		label := g.Name
		if g.Uncategorized {
			label += " (uncategorized)"
		}
		fmt.Printf("%s%s  [%s]\n", indent, label, g.ID)
	}
}
