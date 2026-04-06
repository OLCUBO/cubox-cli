package cmd

import (
	"fmt"
	"strings"

	"github.com/OLCUBO/cubox-cli/internal/client"
	"github.com/OLCUBO/cubox-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	cardGroupFilter   []string
	cardTypeFilter    []string
	cardTagFilter     []string
	cardStarred       bool
	cardRead          bool
	cardUnread        bool
	cardAnnotated     bool
	cardLimit         int
	cardCursor        string
	cardAll           bool
	cardContentID     string
)

var cardCmd = &cobra.Command{
	Use:   "card",
	Short: "Manage cards (bookmarks)",
}

var cardListCmd = &cobra.Command{
	Use:   "list",
	Short: "List and filter cards",
	Long: `Filter and list bookmark cards.

Examples:
  cubox-cli card list
  cubox-cli card list --type Article,Snippet --starred --limit 10
  cubox-cli card list --group 7230156249357091393 --all`,
	RunE: runCardList,
}

var cardContentCmd = &cobra.Command{
	Use:   "content",
	Short: "Get card content in markdown",
	Long: `Retrieve the full article content of a card in markdown format.

Examples:
  cubox-cli card content --id 7247925101516031380
  cubox-cli card content --id 7247925101516031380 -o json`,
	RunE: runCardContent,
}

func init() {
	cardListCmd.Flags().StringSliceVar(&cardGroupFilter, "group", nil, "filter by group IDs (comma-separated)")
	cardListCmd.Flags().StringSliceVar(&cardTypeFilter, "type", nil, "filter by type: Article,Snippet,Memo,Image,Audio,Video,File")
	cardListCmd.Flags().StringSliceVar(&cardTagFilter, "tag", nil, "filter by tag IDs (comma-separated, empty string = no tag)")
	cardListCmd.Flags().BoolVar(&cardStarred, "starred", false, "only starred cards")
	cardListCmd.Flags().BoolVar(&cardRead, "read", false, "only read cards")
	cardListCmd.Flags().BoolVar(&cardUnread, "unread", false, "only unread cards")
	cardListCmd.Flags().BoolVar(&cardAnnotated, "annotated", false, "only annotated cards")
	cardListCmd.Flags().IntVar(&cardLimit, "limit", 50, "page size")
	cardListCmd.Flags().StringVar(&cardCursor, "cursor", "", "pagination cursor: CARD_ID,UPDATE_TIME")
	cardListCmd.Flags().BoolVar(&cardAll, "all", false, "auto-paginate to fetch all results")

	cardContentCmd.Flags().StringVar(&cardContentID, "id", "", "card ID (required)")
	cardContentCmd.MarkFlagRequired("id")

	cardCmd.AddCommand(cardListCmd, cardContentCmd)
	rootCmd.AddCommand(cardCmd)
}

func runCardList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c := client.New(cfg.BaseURL(), cfg.Token)

	req := &client.CardFilterRequest{
		GroupFilters: cardGroupFilter,
		TypeFilters:  cardTypeFilter,
		TagFilters:   cardTagFilter,
		Limit:        cardLimit,
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

	if cardCursor != "" {
		parts := strings.SplitN(cardCursor, ",", 2)
		if len(parts) == 2 {
			req.LastCardID = parts[0]
			req.LastCardUpdate = parts[1]
		} else {
			return fmt.Errorf("invalid cursor format, expected CARD_ID,UPDATE_TIME")
		}
	}

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
	for {
		cards, err := c.FilterCards(req)
		if err != nil {
			return err
		}
		if len(cards) == 0 {
			break
		}
		allCards = append(allCards, cards...)
		last := cards[len(cards)-1]
		req.LastCardID = last.ID
		req.LastCardUpdate = last.UpdateTime
	}

	if outputFormat == "text" {
		printCardsText(allCards)
		return nil
	}
	printJSON(allCards)
	return nil
}

func runCardContent(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c := client.New(cfg.BaseURL(), cfg.Token)

	content, err := c.GetCardContent(cardContentID)
	if err != nil {
		return err
	}

	if outputFormat == "text" || outputFormat == "json" {
		fmt.Print(content)
		return nil
	}
	printJSON(map[string]string{"id": cardContentID, "content": content})
	return nil
}

func printCardsText(cards []client.Card) {
	for _, c := range cards {
		tags := ""
		if len(c.Tags) > 0 {
			tags = " [" + strings.Join(c.Tags, ", ") + "]"
		}
		fmt.Printf("%s  %s  (%s)%s\n", c.ID, c.Title, c.Type, tags)
		if c.Description != "" {
			fmt.Printf("    %s\n", c.Description)
		}
		if c.URL != "" {
			fmt.Printf("    %s\n", c.URL)
		}
		if len(c.Highlights) > 0 {
			fmt.Printf("    %d highlight(s)\n", len(c.Highlights))
		}
		fmt.Println()
	}
}
