package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/OLCUBO/cubox-cli/internal/client"
	"github.com/OLCUBO/cubox-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	saveFolderID string
	saveTagIDs  []string
	saveTitle   string
	saveDesc    string
	saveJSON    string
)

var saveCmd = &cobra.Command{
	Use:   "save [urls...]",
	Short: "Save web pages as bookmarks",
	Long: `Save one or more web pages to your Cubox collection.

Three input modes:

1. URL arguments (simple — just URLs):
   cubox-cli save https://example.com https://another.com

2. Single card with metadata:
   cubox-cli save --url https://example.com --title "Example" --desc "A description"

3. Batch via JSON (full control):
   cubox-cli save --json '[{"url":"https://a.com","title":"A"},{"url":"https://b.com"}]'

All modes support --folder and --tag flags.

Examples:
  cubox-cli save https://example.com
  cubox-cli save https://a.com https://b.com --folder 7295054384872820655
  cubox-cli save --url https://example.com --title "My Page" --desc "Interesting read"
  cubox-cli save --json '[{"url":"https://a.com","title":"Title A"}]' --tag 7247925099053977508`,
	RunE: runSave,
}

func init() {
	saveCmd.Flags().StringVar(&saveFolderID, "folder", "", "target folder ID")
	saveCmd.Flags().StringSliceVar(&saveTagIDs, "tag", nil, "tag IDs to apply (comma-separated)")
	saveCmd.Flags().StringVar(&saveTitle, "title", "", "title for the saved page (use with --url)")
	saveCmd.Flags().StringVar(&saveDesc, "desc", "", "description for the saved page (use with --url)")
	saveCmd.Flags().StringVar(&saveJSON, "json", "", `batch card entries as JSON array: [{"url":"...","title":"...","description":"..."}]`)

	rootCmd.AddCommand(saveCmd)
}

func runSave(cmd *cobra.Command, args []string) error {
	cards, err := buildSaveCards(args)
	if err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c := client.New(cfg.BaseURL(), cfg.Token)

	req := &client.SaveCardsRequest{
		Cards:    cards,
		FolderID: saveFolderID,
		TagIDs:  saveTagIDs,
	}

	if err := c.SaveCards(req); err != nil {
		return err
	}

	fmt.Printf("Saved %d card(s) successfully.\n", len(cards))
	return nil
}

func buildSaveCards(args []string) ([]client.SaveCardEntry, error) {
	hasJSON := saveJSON != ""
	hasArgs := len(args) > 0
	hasMeta := saveTitle != "" || saveDesc != ""

	if hasJSON && hasArgs {
		return nil, fmt.Errorf("cannot use both --json and URL arguments")
	}
	if hasJSON && hasMeta {
		return nil, fmt.Errorf("cannot use --title/--desc with --json")
	}
	if !hasJSON && !hasArgs {
		return nil, fmt.Errorf("provide URLs as arguments, or use --json for batch input")
	}

	if hasJSON {
		var cards []client.SaveCardEntry
		if err := json.Unmarshal([]byte(saveJSON), &cards); err != nil {
			return nil, fmt.Errorf("invalid --json: %w", err)
		}
		for i, c := range cards {
			if c.URL == "" {
				return nil, fmt.Errorf("--json entry %d: url is required", i)
			}
		}
		return cards, nil
	}

	if hasMeta && len(args) > 1 {
		return nil, fmt.Errorf("--title/--desc apply to a single URL; pass one URL or use --json for batch")
	}

	cards := make([]client.SaveCardEntry, 0, len(args))
	for _, u := range args {
		cards = append(cards, client.SaveCardEntry{
			URL:         u,
			Title:       saveTitle,
			Description: saveDesc,
		})
	}
	return cards, nil
}
