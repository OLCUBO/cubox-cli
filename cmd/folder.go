package cmd

import (
	"fmt"
	"strings"

	"github.com/OLCUBO/cubox-cli/internal/client"
	"github.com/OLCUBO/cubox-cli/internal/config"
	"github.com/spf13/cobra"
)

var folderCmd = &cobra.Command{
	Use:   "folder",
	Short: "Manage folders",
}

var folderListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all folders",
	RunE:  runFolderList,
}

func init() {
	folderCmd.AddCommand(folderListCmd)
	rootCmd.AddCommand(folderCmd)
}

func runFolderList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c := client.New(cfg.BaseURL(), cfg.Token)
	folders, err := c.ListFolders()
	if err != nil {
		return err
	}

	if outputFormat == "text" {
		printFoldersText(folders)
		return nil
	}
	printJSON(folders)
	return nil
}

func printFoldersText(folders []client.Folder) {
	for _, f := range folders {
		depth := strings.Count(f.NestedName, "/")
		indent := strings.Repeat("  ", depth)
		label := f.Name
		if f.Uncategorized {
			label += " (uncategorized)"
		}
		fmt.Printf("%s%s  [%s]\n", indent, label, f.ID)
	}
}
