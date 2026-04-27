package cmd

import (
	"fmt"

	"github.com/OLCUBO/cubox-cli/internal/client"
	"github.com/OLCUBO/cubox-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	archiveIDs []string

	unarchiveIDs    []string
	unarchiveFolder string
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "Archive one or more cards (batch)",
	Long: `Archive one or more cards by ID. Archived cards are hidden from
the default card list (use "card list --archived" to view them).

Use "cubox-cli unarchive" to restore archived cards back to a folder.

Examples:
  cubox-cli archive --id 7444025677600260245
  cubox-cli archive --id 7444025677600260245,7443973659296793971`,
	RunE: runArchive,
}

var unarchiveCmd = &cobra.Command{
	Use:   "unarchive",
	Short: "Restore archived cards to a folder (batch)",
	Long: `Restore one or more archived cards by moving them into a non-archived
folder.

A destination folder is required. Specify it by name (including nested
paths like "parent/child"), or pass an empty string to move into the
"Uncategorized" folder.

Examples:
  cubox-cli unarchive --id 7444025677600260245 --folder "Reading List"
  cubox-cli unarchive --id 7444025677600260245,7443973659296793971 --folder "parent/child"
  cubox-cli unarchive --id 7444025677600260245 --folder ""`,
	RunE: runUnarchive,
}

func init() {
	archiveCmd.Flags().StringSliceVar(&archiveIDs, "id", nil, "card IDs to archive (comma-separated, required)")
	archiveCmd.MarkFlagRequired("id")

	unarchiveCmd.Flags().StringSliceVar(&unarchiveIDs, "id", nil, "card IDs to unarchive (comma-separated, required)")
	unarchiveCmd.MarkFlagRequired("id")
	unarchiveCmd.Flags().StringVar(&unarchiveFolder, "folder", "", "destination folder name (required; \"\" = Uncategorized)")
	unarchiveCmd.MarkFlagRequired("folder")

	rootCmd.AddCommand(archiveCmd)
	rootCmd.AddCommand(unarchiveCmd)
}

type archiveResult struct {
	Count   int    `json:"count"`
	Message string `json:"message"`
}

func runArchive(cmd *cobra.Command, args []string) error {
	if len(archiveIDs) == 0 {
		return fmt.Errorf("at least one card ID is required")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c := client.New(cfg.BaseURL(), cfg.Token)

	if err := c.ArchiveCards(archiveIDs); err != nil {
		return err
	}

	count := len(archiveIDs)
	result := archiveResult{
		Count:   count,
		Message: fmt.Sprintf("Successfully archived %d card(s).", count),
	}

	if outputFormat == "text" {
		fmt.Printf("Successfully archived %d card(s).\n", count)
		return nil
	}
	printJSON(result)
	return nil
}

func runUnarchive(cmd *cobra.Command, args []string) error {
	if len(unarchiveIDs) == 0 {
		return fmt.Errorf("at least one card ID is required")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c := client.New(cfg.BaseURL(), cfg.Token)

	folderID, err := resolveFolderID(c, unarchiveFolder)
	if err != nil {
		return err
	}

	req := &client.MoveCardsRequest{
		FolderID: folderID,
		CardIDs:  unarchiveIDs,
	}
	if err := c.MoveCards(req); err != nil {
		return err
	}

	count := len(unarchiveIDs)
	result := archiveResult{
		Count:   count,
		Message: fmt.Sprintf("Successfully unarchived %d card(s).", count),
	}

	if outputFormat == "text" {
		fmt.Printf("Successfully unarchived %d card(s).\n", count)
		return nil
	}
	printJSON(result)
	return nil
}

// resolveFolderID looks up a folder by its nested_name (e.g. "parent/child"
// or "" for the user's Uncategorized folder) and returns its ID.
func resolveFolderID(c *client.Client, nestedName string) (string, error) {
	folders, err := c.ListFolders()
	if err != nil {
		return "", fmt.Errorf("listing folders: %w", err)
	}

	if nestedName == "" {
		for _, f := range folders {
			if f.Uncategorized {
				return f.ID, nil
			}
		}
		return "", fmt.Errorf("could not find Uncategorized folder; pass --folder NAME with an explicit folder name")
	}

	for _, f := range folders {
		if f.NestedName == nestedName {
			return f.ID, nil
		}
	}
	return "", fmt.Errorf("folder not found: %q (use \"cubox-cli folder list\" to see available folders)", nestedName)
}
