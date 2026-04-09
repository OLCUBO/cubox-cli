package cmd

import (
	"fmt"
	"strings"

	"github.com/OLCUBO/cubox-cli/internal/client"
	"github.com/OLCUBO/cubox-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	markColors    []string
	markLastID    string
	markLimit     int
	markKeyword   string
	markStartTime string
	markEndTime   string
	markAll       bool
)

var markCmd = &cobra.Command{
	Use:   "mark",
	Short: "Manage highlights/annotations",
}

var markListCmd = &cobra.Command{
	Use:   "list",
	Short: "List highlights/annotations",
	Long: `Filter and list highlights (annotations) across all cards.

Examples:
  cubox-cli mark list
  cubox-cli mark list --color Yellow,Purple --limit 10
  cubox-cli mark list --keyword "important" --all`,
	RunE: runMarkList,
}

func init() {
	markListCmd.Flags().StringSliceVar(&markColors, "color", nil, "filter by color: Yellow,Green,Blue,Pink,Purple")
	markListCmd.Flags().StringVar(&markLastID, "last-id", "", "last highlight ID for cursor pagination")
	markListCmd.Flags().IntVar(&markLimit, "limit", 50, "page size")
	markListCmd.Flags().StringVar(&markKeyword, "keyword", "", "search keyword")
	markListCmd.Flags().StringVar(&markStartTime, "start-time", "", "filter by start time")
	markListCmd.Flags().StringVar(&markEndTime, "end-time", "", "filter by end time")
	markListCmd.Flags().BoolVar(&markAll, "all", false, "auto-paginate to fetch all results")

	markCmd.AddCommand(markListCmd)
	rootCmd.AddCommand(markCmd)
}

func runMarkList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c := client.New(cfg.BaseURL(), cfg.Token)

	req := &client.MarkFilterRequest{
		Colors:          markColors,
		LastHighlightID: markLastID,
		Limit:           markLimit,
		Keyword:         markKeyword,
		StartTime:       markStartTime,
		EndTime:         markEndTime,
	}

	if markAll {
		return runMarkListAll(c, req)
	}

	marks, err := c.FilterMarks(req)
	if err != nil {
		return err
	}

	if outputFormat == "text" {
		printMarksText(marks)
		return nil
	}
	printJSON(marks)
	return nil
}

func runMarkListAll(c *client.Client, req *client.MarkFilterRequest) error {
	var allMarks []client.Mark
	for {
		marks, err := c.FilterMarks(req)
		if err != nil {
			return err
		}
		if len(marks) == 0 {
			break
		}
		allMarks = append(allMarks, marks...)
		req.LastHighlightID = marks[len(marks)-1].ID
	}

	if outputFormat == "text" {
		printMarksText(allMarks)
		return nil
	}
	printJSON(allMarks)
	return nil
}

func printMarksText(marks []client.Mark) {
	for _, m := range marks {
		colorTag := strings.ToLower(m.Color)
		text := m.Text
		if text == "" && m.ImageURL != "" {
			text = "[image] " + m.ImageURL
		}
		fmt.Printf("[%s] %s\n", colorTag, text)
		if m.Note != "" {
			fmt.Printf("  note: %s\n", m.Note)
		}
		fmt.Printf("  card: %s  id: %s\n\n", m.CardID, m.ID)
	}
}
