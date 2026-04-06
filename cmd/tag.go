package cmd

import (
	"fmt"
	"strings"

	"github.com/OLCUBO/cubox-cli/internal/client"
	"github.com/OLCUBO/cubox-cli/internal/config"
	"github.com/spf13/cobra"
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage tags",
}

var tagListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tags",
	RunE:  runTagList,
}

func init() {
	tagCmd.AddCommand(tagListCmd)
	rootCmd.AddCommand(tagCmd)
}

func runTagList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c := client.New(cfg.BaseURL(), cfg.Token)
	tags, err := c.ListTags()
	if err != nil {
		return err
	}

	if outputFormat == "text" {
		printTagsText(tags)
		return nil
	}
	printJSON(tags)
	return nil
}

func printTagsText(tags []client.Tag) {
	for _, t := range tags {
		depth := strings.Count(t.NestedName, "/")
		indent := strings.Repeat("  ", depth)
		fmt.Printf("%s%s  [%s]\n", indent, t.Name, t.ID)
	}
}
