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
	updateGroupID     string
	updateAddTagIDs   []string
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a card's properties",
	Long: `Update a card — star/unstar, mark read/unread, archive,
move to a group, or add tags, title, description.

Examples:
  cubox-cli update --id 7431974288044854951 --star
  cubox-cli update --id 7431974288044854951 --read --group 7295054384872820655
  cubox-cli update --id 7431974288044854951 --add-tag 7247925099053977508 --title "New Title" --description "New Description"`,
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
	updateCmd.Flags().StringVar(&updateGroupID, "group", "", "move to group/folder ID")
	updateCmd.Flags().StringSliceVar(&updateAddTagIDs, "add-tag", nil, "tag IDs to add (comma-separated)")
	updateCmd.Flags().StringVar(&updateTitle, "title", "", "title to update")
	updateCmd.Flags().StringVar(&updateDescription, "description", "", "description to update")

	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c := client.New(cfg.BaseURL(), cfg.Token)

	req := &client.CardUpdateRequest{
		ID:        updateID,
		GroupID:   updateGroupID,
		AddTagIDs: updateAddTagIDs,
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
