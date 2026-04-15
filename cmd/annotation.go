package cmd

import (
	"fmt"

	"github.com/OLCUBO/cubox-cli/internal/client"
	"github.com/OLCUBO/cubox-cli/internal/config"
	"github.com/OLCUBO/cubox-cli/internal/timefmt"
	"github.com/spf13/cobra"
)

var (
	annotationColors    []string
	annotationKeyword   string
	annotationLimit     int
	annotationLastID    string
	annotationAll       bool
	annotationStartTime string
	annotationEndTime   string
)

var annotationCmd = &cobra.Command{
	Use:   "annotation",
	Short: "Manage annotations",
}

var annotationListCmd = &cobra.Command{
	Use:   "list",
	Short: "List and search annotations",
	Long: `List and search annotations across all cards.

Examples:
  cubox-cli annotation list
  cubox-cli annotation list --color Yellow,Green
  cubox-cli annotation list --keyword "database"
  cubox-cli annotation list --all`,
	RunE: runAnnotationList,
}

func init() {
	annotationListCmd.Flags().StringSliceVar(&annotationColors, "color", nil, "filter by colors (comma-separated)")
	annotationListCmd.Flags().StringVar(&annotationKeyword, "keyword", "", "search annotations by keyword")
	annotationListCmd.Flags().IntVar(&annotationLimit, "limit", 50, "page size")
	annotationListCmd.Flags().StringVar(&annotationLastID, "last-id", "", "last annotation ID for cursor pagination")
	annotationListCmd.Flags().BoolVar(&annotationAll, "all", false, "auto-paginate to fetch all results")
	annotationListCmd.Flags().StringVar(&annotationStartTime, "start-time", "", "filter start time (today, yesterday, 7d, 2006-01-02, or full timestamp)")
	annotationListCmd.Flags().StringVar(&annotationEndTime, "end-time", "", "filter end time (today, yesterday, 7d, 2006-01-02, or full timestamp)")

	annotationCmd.AddCommand(annotationListCmd)
	rootCmd.AddCommand(annotationCmd)
}

func buildAnnotationFilterRequest() (*client.AnnotationFilterRequest, error) {
	startTime, err := timefmt.Parse(annotationStartTime)
	if err != nil {
		return nil, fmt.Errorf("--start-time: %w", err)
	}
	endTime, err := timefmt.ParseEnd(annotationEndTime)
	if err != nil {
		return nil, fmt.Errorf("--end-time: %w", err)
	}

	return &client.AnnotationFilterRequest{
		Colors:           annotationColors,
		LastAnnotationID: annotationLastID,
		Limit:            annotationLimit,
		Keyword:          annotationKeyword,
		StartTime:        startTime,
		EndTime:          endTime,
	}, nil
}

func runAnnotationList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	c := client.New(cfg.BaseURL(), cfg.Token)

	req, err := buildAnnotationFilterRequest()
	if err != nil {
		return err
	}

	if annotationAll {
		return runAnnotationListAll(c, req)
	}

	annotations, err := c.FilterAnnotations(req)
	if err != nil {
		return err
	}

	printJSON(annotations)
	return nil
}

func runAnnotationListAll(c *client.Client, req *client.AnnotationFilterRequest) error {
	var allAnnotations []client.Annotation

	for {
		annotations, err := c.FilterAnnotations(req)
		if err != nil {
			return err
		}
		if len(annotations) == 0 {
			break
		}
		allAnnotations = append(allAnnotations, annotations...)
		req.LastAnnotationID = annotations[len(annotations)-1].ID
	}

	printJSON(allAnnotations)
	return nil
}
