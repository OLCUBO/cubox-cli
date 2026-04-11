package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/OLCUBO/cubox-cli/internal/update"
	"github.com/spf13/cobra"
)

var outputFormat string

// updateCh receives the result of the background update check.
var updateCh = make(chan *update.Info, 1)

var rootCmd = &cobra.Command{
	Use:   "cubox-cli",
	Short: "The official Cubox CLI tool, built for humans and AI Agents",
	Long: `cubox-cli — manage your Cubox bookmarks from the terminal.

Supports listing groups (folders), tags, filtering cards, and reading
card content. Designed for both interactive human use and AI Agent
automation with structured JSON output.`,
	SilenceUsage:  true,
	SilenceErrors: true,

	// Start the update check concurrently so it overlaps with command execution.
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		go func() {
			updateCh <- update.Check(Version)
		}()
	},

	// After each command, collect the update result (non-blocking) and emit
	// a notice if a newer version is available.
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		var info *update.Info
		select {
		case info = <-updateCh:
		default:
		}
		if info == nil {
			return
		}
		if outputFormat == "text" {
			fmt.Fprintf(os.Stderr, "\nUpdate available: %s -> %s\nRun: %s\n",
				info.Current, info.Latest, info.Command)
		}
		// In JSON modes the notice is injected by printJSON; store it for use there.
		pendingUpdate = info
	},
}

// pendingUpdate is set by PersistentPostRun for the JSON printer to consume.
// It is only accessed from the main goroutine (cobra runs hooks sequentially).
var pendingUpdate *update.Info

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json", "output format: json, pretty, text")
}

func Execute() error {
	return rootCmd.Execute()
}

// printJSON serialises v as the command output. When a newer version has been
// detected, the payload is wrapped in {"data":...,"_notice":{"update":...}}.
// Without a pending update the raw value is printed directly (no wrapper),
// preserving backward compatibility for callers that parse the output.
func printJSON(v interface{}) {
	if pendingUpdate != nil {
		wrapped := map[string]interface{}{
			"data": v,
			"_notice": map[string]interface{}{
				"update": pendingUpdate,
			},
		}
		pendingUpdate = nil // consume so it is only printed once
		printRaw(wrapped)
		return
	}
	printRaw(v)
}

func printRaw(v interface{}) {
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
