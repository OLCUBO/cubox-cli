package cmd

import (
	"fmt"

	"github.com/OLCUBO/cubox-cli/internal/client"
	"github.com/OLCUBO/cubox-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	updateID          string
	updateTitle       string
	updateDescription string
	updateStar        bool
	updateUnstar      bool
	updateRead        bool
	updateUnread      bool
	updateArchive     bool
	updateFolder      string
	updateTags        []string
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a card's properties",
	Long: `Update a card — star/unstar, mark read/unread, archive,
move to a folder, set tags, update title or description.

Folders and tags are specified by name (including nested paths
like "parent/child"), not by ID.

Examples:
  cubox-cli update --id 7431974288044854951 --star
  cubox-cli update --id 7431974288044854951 --read --folder "Reading List"
  cubox-cli update --id 7431974288044854951 --folder ""
  cubox-cli update --id 7431974288044854951 --tag tech,AI/LLM --title "New Title"`,
	RunE: runUpdate,
}

func init() {
	updateCmd.Flags().StringVar(&updateID, "id", "", "card ID (required)")
	updateCmd.MarkFlagRequired("id")
	updateCmd.Flags().BoolVar(&updateStar, "star", false, "star the card")
	updateCmd.Flags().BoolVar(&updateUnstar, "unstar", false, "unstar the card")
	updateCmd.Flags().BoolVar(&updateRead, "read", false, "mark as read")
	updateCmd.Flags().BoolVar(&updateUnread, "unread", false, "mark as unread")
	updateCmd.Flags().BoolVar(&updateArchive, "archive", false, "archive the card")
	updateCmd.Flags().StringVar(&updateFolder, "folder", "", "move to folder by name (e.g. \"parent/child\"; \"\" = Uncategorized)")
	updateCmd.Flags().StringSliceVar(&updateTags, "tag", nil, "set tag names (comma-separated, supports nested like \"parent/child\")")
	updateCmd.Flags().StringVar(&updateTitle, "title", "", "update title")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "update description")

	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c := client.New(cfg.BaseURL(), cfg.Token)

	req := &client.CardUpdateRequest{
		ID:             updateID,
		TagNestedNames: updateTags,
	}

	if cmd.Flags().Changed("folder") {
		req.FolderNestedName = &updateFolder
	}

	if updateStar {
		v := true
		req.Starred = &v
	}
	if updateUnstar {
		v := false
		req.Starred = &v
	}
	if updateRead {
		v := true
		req.Read = &v
	}
	if updateUnread {
		v := false
		req.Read = &v
	}
	if updateArchive {
		v := true
		req.Archive = &v
	}
	if updateTitle != "" {
		req.Title = updateTitle
	}
	if updateDescription != "" {
		req.Description = updateDescription
	}

	if err := c.UpdateCard(req); err != nil {
		return err
	}

	fmt.Println("Card updated successfully.")
	return nil
}
