package cmd

import (
	"fmt"
	"strings"

	"github.com/OLCUBO/cubox-cli/internal/client"
	"github.com/OLCUBO/cubox-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	cardGroupFilter []string
	cardTagFilter   []string
	cardStarred     bool
	cardRead        bool
	cardUnread      bool
	cardAnnotated   bool
	cardLimit       int
	cardLastID      string
	cardAll         bool
	cardKeyword     string
	cardPage        int
	cardStartTime   string
	cardEndTime     string
	cardDetailID    string
)

var cardCmd = &cobra.Command{
	Use:   "card",
	Short: "Manage cards (bookmarks)",
}

var cardListCmd = &cobra.Command{
	Use:   "list",
	Short: "List and filter cards",
	Long: `Filter and list bookmark cards. Supports keyword search.

When using --keyword for search, pagination uses --page (1-based).
Without --keyword, pagination uses --last-id (cursor-based).

Examples:
  cubox-cli card list
  cubox-cli card list --starred --limit 10
  cubox-cli card list --group 7230156249357091393 --all
  cubox-cli card list --keyword "AI agent" --page 1
  cubox-cli card list --start-time "2026-01-01T00:00:00:000+08:00"`,
	RunE: runCardList,
}

var cardDetailCmd = &cobra.Command{
	Use:   "detail",
	Short: "Get full card detail with content, highlights, and AI insight",
	Long: `Retrieve the full detail of a card including article content (markdown),
highlights, and AI-generated insight (summary + Q&A).

Examples:
  cubox-cli card detail --id 7247925101516031380
  cubox-cli card detail --id 7247925101516031380 -o pretty`,
	RunE: runCardDetail,
}

func init() {
	cardListCmd.Flags().StringSliceVar(&cardGroupFilter, "group", nil, "filter by group IDs (comma-separated)")
	cardListCmd.Flags().StringSliceVar(&cardTagFilter, "tag", nil, "filter by tag IDs (comma-separated, empty string = no tag)")
	cardListCmd.Flags().BoolVar(&cardStarred, "starred", false, "only starred cards")
	cardListCmd.Flags().BoolVar(&cardRead, "read", false, "only read cards")
	cardListCmd.Flags().BoolVar(&cardUnread, "unread", false, "only unread cards")
	cardListCmd.Flags().BoolVar(&cardAnnotated, "annotated", false, "only annotated cards")
	cardListCmd.Flags().IntVar(&cardLimit, "limit", 50, "page size")
	cardListCmd.Flags().StringVar(&cardLastID, "last-id", "", "last card ID for cursor pagination (non-search)")
	cardListCmd.Flags().BoolVar(&cardAll, "all", false, "auto-paginate to fetch all results")
	cardListCmd.Flags().StringVar(&cardKeyword, "keyword", "", "search keyword")
	cardListCmd.Flags().IntVar(&cardPage, "page", 0, "page number for search pagination (1-based)")
	cardListCmd.Flags().StringVar(&cardStartTime, "start-time", "", "filter by create time start (e.g. 2026-01-01T00:00:00:000+08:00)")
	cardListCmd.Flags().StringVar(&cardEndTime, "end-time", "", "filter by create time end")

	cardDetailCmd.Flags().StringVar(&cardDetailID, "id", "", "card ID (required)")
	cardDetailCmd.MarkFlagRequired("id")

	cardCmd.AddCommand(cardListCmd, cardDetailCmd)
	rootCmd.AddCommand(cardCmd)
}

func buildCardFilterRequest() *client.CardFilterRequest {
	req := &client.CardFilterRequest{
		GroupFilters: cardGroupFilter,
		TagFilters:   cardTagFilter,
		Limit:        cardLimit,
		Keyword:      cardKeyword,
		StartTime:    cardStartTime,
		EndTime:      cardEndTime,
	}

	if cardStarred {
		v := true
		req.Starred = &v
	}
	if cardRead {
		v := true
		req.Read = &v
	}
	if cardUnread {
		v := false
		req.Read = &v
	}
	if cardAnnotated {
		v := true
		req.Annotated = &v
	}

	if cardKeyword != "" {
		if cardPage > 0 {
			req.Page = cardPage
		} else {
			req.Page = 1
		}
	} else {
		req.LastCardID = cardLastID
	}

	return req
}

func runCardList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c := client.New(cfg.BaseURL(), cfg.Token)

	req := buildCardFilterRequest()

	if cardAll {
		return runCardListAll(c, req)
	}

	cards, err := c.FilterCards(req)
	if err != nil {
		return err
	}

	if outputFormat == "text" {
		printCardsText(cards)
		return nil
	}
	printJSON(cards)
	return nil
}

func runCardListAll(c *client.Client, req *client.CardFilterRequest) error {
	var allCards []client.Card

	if req.Keyword != "" {
		for page := 1; ; page++ {
			req.Page = page
			cards, err := c.FilterCards(req)
			if err != nil {
				return err
			}
			if len(cards) == 0 {
				break
			}
			allCards = append(allCards, cards...)
		}
	} else {
		for {
			cards, err := c.FilterCards(req)
			if err != nil {
				return err
			}
			if len(cards) == 0 {
				break
			}
			allCards = append(allCards, cards...)
			req.LastCardID = cards[len(cards)-1].ID
		}
	}

	if outputFormat == "text" {
		printCardsText(allCards)
		return nil
	}
	printJSON(allCards)
	return nil
}

func runCardDetail(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c := client.New(cfg.BaseURL(), cfg.Token)

	detail, err := c.GetCardDetail(cardDetailID)
	if err != nil {
		return err
	}

	if outputFormat == "text" {
		fmt.Print(detail.Content)
		return nil
	}
	printJSON(detail)
	return nil
}

func printCardsText(cards []client.Card) {
	for _, c := range cards {
		star := " "
		if c.Starred {
			star = "*"
		}
		readMark := " "
		if c.Read {
			readMark = "R"
		}
		tags := ""
		if len(c.Tags) > 0 {
			tags = " [" + strings.Join(c.Tags, ", ") + "]"
		}
		fmt.Printf("[%s%s] %s  %s%s\n", star, readMark, c.ID, c.Title, tags)
		if c.URL != "" {
			fmt.Printf("     %s\n", c.URL)
		}
		fmt.Println()
	}
}
