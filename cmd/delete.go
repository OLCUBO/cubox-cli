package cmd

import (
	"fmt"
	"os"

	"github.com/OLCUBO/cubox-cli/internal/client"
	"github.com/OLCUBO/cubox-cli/internal/config"
	"github.com/spf13/cobra"
)

const detailPreviewThreshold = 3

var (
	deleteIDs    []string
	deleteDryRun bool
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete cards",
	Long: `Delete one or more cards by ID.

Use --dry-run to preview which cards would be deleted without actually
removing them. This is especially recommended when scripting or when
an AI Agent is performing deletions on behalf of a user.

When deleting up to 3 cards, the CLI fetches card details (title, URL)
for the preview. For larger batches, only the count and IDs are shown
to avoid expensive per-card API calls.

Examples:
  cubox-cli delete --id 7435692934957108160
  cubox-cli delete --id 7435692934957108160,7435691601617225646
  cubox-cli delete --id 7435692934957108160 --dry-run`,
	RunE: runDelete,
}

func init() {
	deleteCmd.Flags().StringSliceVar(&deleteIDs, "id", nil, "card IDs to delete (comma-separated, required)")
	deleteCmd.MarkFlagRequired("id")
	deleteCmd.Flags().BoolVar(&deleteDryRun, "dry-run", false, "preview cards to be deleted without actually deleting")

	rootCmd.AddCommand(deleteCmd)
}

type deletePreview struct {
	ID    string `json:"id"`
	Title string `json:"title,omitempty"`
	URL   string `json:"url,omitempty"`
}

type deleteResult struct {
	DryRun  bool            `json:"dry_run"`
	Count   int             `json:"count"`
	Cards   []deletePreview `json:"cards,omitempty"`
	Message string          `json:"message"`
}

func runDelete(cmd *cobra.Command, args []string) error {
	if len(deleteIDs) == 0 {
		return fmt.Errorf("at least one card ID is required")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c := client.New(cfg.BaseURL(), cfg.Token)

	if deleteDryRun {
		return outputDryRun(c, deleteIDs)
	}
	return executeDelete(c, deleteIDs)
}

func outputDryRun(c *client.Client, ids []string) error {
	count := len(ids)

	if count <= detailPreviewThreshold {
		previews := fetchDeletePreviews(c, ids)
		result := deleteResult{
			DryRun:  true,
			Count:   count,
			Cards:   previews,
			Message: fmt.Sprintf("[dry-run] Would delete %d card(s). No changes were made.", count),
		}
		if outputFormat == "text" {
			fmt.Fprintf(os.Stderr, "[dry-run] Would delete %d card(s):\n\n", count)
			printDeletePreviewText(previews)
			fmt.Fprintln(os.Stderr, "No changes were made. Remove --dry-run to delete.")
			return nil
		}
		printJSON(result)
		return nil
	}

	result := deleteResult{
		DryRun:  true,
		Count:   count,
		Message: fmt.Sprintf("[dry-run] Would delete %d card(s). No changes were made.", count),
	}
	if outputFormat == "text" {
		fmt.Fprintf(os.Stderr, "[dry-run] Would delete %d card(s). No changes were made.\n", count)
		fmt.Fprintf(os.Stderr, "Remove --dry-run to delete.\n")
		return nil
	}
	printJSON(result)
	return nil
}

func executeDelete(c *client.Client, ids []string) error {
	count := len(ids)

	if outputFormat == "text" && count <= detailPreviewThreshold {
		previews := fetchDeletePreviews(c, ids)
		fmt.Fprintf(os.Stderr, "Deleting %d card(s):\n\n", count)
		printDeletePreviewText(previews)
	}

	if err := c.DeleteCards(ids); err != nil {
		return err
	}

	result := deleteResult{
		DryRun:  false,
		Count:   count,
		Message: fmt.Sprintf("Successfully deleted %d card(s).", count),
	}

	if outputFormat == "text" {
		fmt.Fprintf(os.Stderr, "Successfully deleted %d card(s).\n", count)
		return nil
	}
	printJSON(result)
	return nil
}

// fetchDeletePreviews fetches card details for a small set of IDs to build a
// human-readable preview. Only called when len(ids) <= detailPreviewThreshold.
// Failures are tolerated — the card is still listed with its ID only.
func fetchDeletePreviews(c *client.Client, ids []string) []deletePreview {
	previews := make([]deletePreview, 0, len(ids))
	for _, id := range ids {
		p := deletePreview{ID: id}
		if detail, err := c.GetCardDetail(id); err == nil {
			p.Title = detail.Title
			p.URL = detail.URL
		}
		previews = append(previews, p)
	}
	return previews
}

func printDeletePreviewText(previews []deletePreview) {
	for _, p := range previews {
		title := p.Title
		if title == "" {
			title = "(unknown)"
		}
		fmt.Fprintf(os.Stderr, "  - %s  %s\n", p.ID, title)
		if p.URL != "" {
			fmt.Fprintf(os.Stderr, "    %s\n", p.URL)
		}
	}
	fmt.Fprintln(os.Stderr)
}
