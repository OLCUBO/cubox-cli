package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var outputFormat string

var rootCmd = &cobra.Command{
	Use:   "cubox-cli",
	Short: "The official Cubox CLI tool, built for humans and AI Agents",
	Long: `cubox-cli — manage your Cubox bookmarks from the terminal.

Supports listing groups (folders), tags, filtering cards, and reading
card content. Designed for both interactive human use and AI Agent
automation with structured JSON output.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json", "output format: json, pretty, text")
}

func Execute() error {
	return rootCmd.Execute()
}

func printJSON(v interface{}) {
	switch outputFormat {
	case "pretty":
		data, _ := json.MarshalIndent(v, "", "  ")
		fmt.Println(string(data))
	default:
		data, _ := json.Marshal(v)
		fmt.Println(string(data))
	}
}

func exitError(msg string) {
	fmt.Fprintln(os.Stderr, "Error:", msg)
	os.Exit(1)
}
