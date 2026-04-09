package cmd

import (
	"fmt"

	"github.com/OLCUBO/cubox-cli/internal/client"
	"github.com/OLCUBO/cubox-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	saveGroupID string
	saveTagIDs  []string
)

var saveCmd = &cobra.Command{
	Use:   "save [urls...]",
	Short: "Save web page URLs as bookmarks",
	Long: `Save one or more URLs to your Cubox collection.

Examples:
  cubox-cli save https://example.com
  cubox-cli save https://a.com https://b.com --group 7295054384872820655
  cubox-cli save https://example.com --tag 7247925099053977508,7295070793040398540`,
	Args: cobra.MinimumNArgs(1),
	RunE: runSave,
}

func init() {
	saveCmd.Flags().StringVar(&saveGroupID, "group", "", "target group/folder ID")
	saveCmd.Flags().StringSliceVar(&saveTagIDs, "tag", nil, "tag IDs to apply (comma-separated)")

	rootCmd.AddCommand(saveCmd)
}

func runSave(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c := client.New(cfg.BaseURL(), cfg.Token)

	req := &client.SaveURLsRequest{
		URLs:    args,
		GroupID: saveGroupID,
		TagIDs:  saveTagIDs,
	}

	if err := c.SaveURLs(req); err != nil {
		return err
	}

	fmt.Printf("Saved %d URL(s) successfully.\n", len(args))
	return nil
}
