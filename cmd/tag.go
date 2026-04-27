package cmd

import (
	"fmt"
	"strings"

	"github.com/OLCUBO/cubox-cli/internal/client"
	"github.com/OLCUBO/cubox-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	tagUpdateID      string
	tagUpdateNewName string

	tagDeleteIDs []string

	tagMergeSourceIDs []string
	tagMergeTargetID  string
)

var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Manage tags",
}

var tagListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tags",
	RunE:  runTagList,
}

var tagUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Rename a tag",
	Long: `Rename a tag by ID. The new name applies to the leaf segment
only — nested children stay attached and are reachable through the
new path automatically.

To find a tag's ID, run "cubox-cli tag list".

Examples:
  cubox-cli tag update --id 7295070793040398540 --new-name link
  cubox-cli tag update --id 7295070793040398540 --new-name "Reading List"`,
	RunE: runTagUpdate,
}

var tagDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete one or more tags (batch)",
	Long: `Delete one or more tags by ID. Cards previously tagged with the
deleted tag(s) are kept; only the tag-card association is removed.

To merge tags into another tag instead of dropping them, use
"cubox-cli tag merge".

Examples:
  cubox-cli tag delete --id 7444025677600260245
  cubox-cli tag delete --id 7444025677600260245,7443973659296793971`,
	RunE: runTagDelete,
}

var tagMergeCmd = &cobra.Command{
	Use:   "merge",
	Short: "Merge tags into a target tag (batch)",
	Long: `Merge one or more source tags into a target tag. All cards
associated with the source tags are re-tagged onto the target, and
the source tags are then deleted.

Examples:
  cubox-cli tag merge --source 7342187912403881105,7342187917722258501 --target 7247925099053977508`,
	RunE: runTagMerge,
}

func init() {
	tagUpdateCmd.Flags().StringVar(&tagUpdateID, "id", "", "tag ID to rename (required)")
	tagUpdateCmd.MarkFlagRequired("id")
	tagUpdateCmd.Flags().StringVar(&tagUpdateNewName, "new-name", "", "new leaf name for the tag (required)")
	tagUpdateCmd.MarkFlagRequired("new-name")

	tagDeleteCmd.Flags().StringSliceVar(&tagDeleteIDs, "id", nil, "tag IDs to delete (comma-separated, required)")
	tagDeleteCmd.MarkFlagRequired("id")

	tagMergeCmd.Flags().StringSliceVar(&tagMergeSourceIDs, "source", nil, "source tag IDs to merge (comma-separated, required)")
	tagMergeCmd.MarkFlagRequired("source")
	tagMergeCmd.Flags().StringVar(&tagMergeTargetID, "target", "", "target tag ID to merge into (required)")
	tagMergeCmd.MarkFlagRequired("target")

	tagCmd.AddCommand(tagListCmd, tagUpdateCmd, tagDeleteCmd, tagMergeCmd)
	rootCmd.AddCommand(tagCmd)
}

func runTagList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c := client.New(cfg.BaseURL(), cfg.Token)
	tags, err := c.ListTags()
	if err != nil {
		return err
	}

	if outputFormat == "text" {
		printTagsText(tags)
		return nil
	}
	printJSON(tags)
	return nil
}

type tagOpResult struct {
	Count   int    `json:"count"`
	Message string `json:"message"`
}

func runTagUpdate(cmd *cobra.Command, args []string) error {
	if strings.Contains(tagUpdateNewName, "/") {
		return fmt.Errorf("--new-name must be a single leaf name without \"/\"; nested paths are not supported")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c := client.New(cfg.BaseURL(), cfg.Token)

	if err := c.UpdateTag(&client.TagUpdateRequest{
		ID:   tagUpdateID,
		Name: tagUpdateNewName,
	}); err != nil {
		return err
	}

	msg := fmt.Sprintf("Successfully renamed tag %s to %q.", tagUpdateID, tagUpdateNewName)
	if outputFormat == "text" {
		fmt.Println(msg)
		return nil
	}
	printJSON(tagOpResult{Count: 1, Message: msg})
	return nil
}

func runTagDelete(cmd *cobra.Command, args []string) error {
	if len(tagDeleteIDs) == 0 {
		return fmt.Errorf("at least one tag ID is required")
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c := client.New(cfg.BaseURL(), cfg.Token)

	if err := c.DeleteTags(tagDeleteIDs); err != nil {
		return err
	}

	count := len(tagDeleteIDs)
	msg := fmt.Sprintf("Successfully deleted %d tag(s).", count)
	if outputFormat == "text" {
		fmt.Println(msg)
		return nil
	}
	printJSON(tagOpResult{Count: count, Message: msg})
	return nil
}

func runTagMerge(cmd *cobra.Command, args []string) error {
	if len(tagMergeSourceIDs) == 0 {
		return fmt.Errorf("at least one source tag ID is required")
	}
	for _, id := range tagMergeSourceIDs {
		if id == tagMergeTargetID {
			return fmt.Errorf("--target tag ID %s cannot also appear in --source", id)
		}
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c := client.New(cfg.BaseURL(), cfg.Token)

	if err := c.MergeTags(&client.TagMergeRequest{
		SourceTagIDs: tagMergeSourceIDs,
		TargetTagID:  tagMergeTargetID,
	}); err != nil {
		return err
	}

	count := len(tagMergeSourceIDs)
	msg := fmt.Sprintf("Successfully merged %d source tag(s) into %s.", count, tagMergeTargetID)
	if outputFormat == "text" {
		fmt.Println(msg)
		return nil
	}
	printJSON(tagOpResult{Count: count, Message: msg})
	return nil
}

func printTagsText(tags []client.Tag) {
	for _, t := range tags {
		depth := strings.Count(t.NestedName, "/")
		indent := strings.Repeat("  ", depth)
		fmt.Printf("%s%s  [%s]\n", indent, t.Name, t.ID)
	}
}
