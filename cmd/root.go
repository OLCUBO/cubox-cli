package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/OLCUBO/cubox-cli/internal/update"
	"github.com/spf13/cobra"
)

var outputFormat string

const (
	emptyCacheRefreshWait = 800 * time.Millisecond
	staleCacheRefreshWait = 300 * time.Millisecond
	refreshHardDeadline   = 5 * time.Second
)

var (
	notifier            *update.Notifier
	updateInfo          *update.Info
	updateNoticeEmitted atomic.Bool
	updateRefreshDone   chan struct{}
	updateRefreshWait   time.Duration
	updateRefreshCancel context.CancelFunc
)

var rootCmd = &cobra.Command{
	Use:   "cubox-cli",
	Short: "The official Cubox CLI tool, built for humans and AI Agents",
	Long: `cubox-cli — manage your Cubox bookmarks from the terminal.

Supports listing groups (folders), tags, filtering cards, and reading
card content. Designed for both interactive human use and AI Agent
automation with structured JSON output.`,
	SilenceUsage:  true,
	SilenceErrors: true,

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		prepareUpdateNotice()
	},
}

// isCompletionInvocation reports whether os.Args looks like a Cobra shell
// completion invocation. We must suppress all update-check side effects
// (network I/O, cache writes) on these paths to keep tab completion fast and
// to avoid corrupting machine-parseable completion output.
func isCompletionInvocation(args []string) bool {
	for _, a := range args {
		switch a {
		case "completion", "__complete", "__completeNoDesc":
			return true
		}
	}
	return false
}

func prepareUpdateNotice() {
	updateInfo = nil
	updateRefreshDone = nil
	updateRefreshWait = 0
	updateNoticeEmitted.Store(false)
	if updateRefreshCancel != nil {
		updateRefreshCancel()
		updateRefreshCancel = nil
	}

	if isCompletionInvocation(os.Args) {
		return
	}

	notifier = update.Default(Version)
	state := notifier.Lookup()
	updateInfo = state.Info
	if !state.NeedsRefresh {
		return
	}

	if state.CacheEmpty {
		updateRefreshWait = emptyCacheRefreshWait
	} else {
		updateRefreshWait = staleCacheRefreshWait
	}
	ctx, cancel := context.WithTimeout(context.Background(), refreshHardDeadline)
	updateRefreshCancel = cancel
	updateRefreshDone = make(chan struct{})
	go func() {
		defer close(updateRefreshDone)
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "update check panic: %v\n", r)
			}
		}()
		notifier.RefreshCache(ctx)
	}()
}

func waitForUpdateRefresh() {
	if updateRefreshDone == nil || updateRefreshWait <= 0 {
		return
	}

	select {
	case <-updateRefreshDone:
		if notifier != nil {
			updateInfo = notifier.Lookup().Info
		}
	case <-time.After(updateRefreshWait):
	}
}

// stopBackgroundRefresh cancels any in-flight RefreshCache goroutine. It is
// idempotent and safe to call after waitForUpdateRefresh.
func stopBackgroundRefresh() {
	if updateRefreshCancel != nil {
		updateRefreshCancel()
	}
}

func emitUpdateNotice() {
	if updateInfo == nil || updateNoticeEmitted.Load() {
		return
	}
	fmt.Fprintf(os.Stderr, "\nUpdate available: %s -> %s\nRun: %s\n",
		updateInfo.Current, updateInfo.Latest, updateInfo.Command)
	updateNoticeEmitted.Store(true)
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "json", "output format: json, pretty, text")
}

func Execute() error {
	err := rootCmd.Execute()
	waitForUpdateRefresh()
	stopBackgroundRefresh()
	if err != nil {
		emitErrorOutput(err)
	}
	emitUpdateNotice()
	return err
}

// emitErrorOutput writes a top-level error envelope to stderr in the active
// output format. When an update notice is known, JSON/pretty output gains a
// `_notice.update` sibling so agents see the upgrade hint even on the error
// path.
func emitErrorOutput(err error) {
	switch outputFormat {
	case "json", "pretty":
		envelope := map[string]interface{}{"error": err.Error()}
		if updateInfo != nil {
			envelope["_notice"] = map[string]interface{}{"update": updateInfo}
			updateNoticeEmitted.Store(true)
		}
		var data []byte
		if outputFormat == "pretty" {
			data, _ = json.MarshalIndent(envelope, "", "  ")
		} else {
			data, _ = json.Marshal(envelope)
		}
		fmt.Fprintln(os.Stderr, string(data))
	default:
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

// printJSON serialises v as the command output. When an update is known from
// the startup cache lookup, the payload is wrapped in
// {"data":...,"_notice":{"update":...}}. Without an update the raw value is
// printed directly, preserving backward compatibility.
func printJSON(v interface{}) {
	if updateInfo != nil {
		wrapped := map[string]interface{}{
			"data": v,
			"_notice": map[string]interface{}{
				"update": updateInfo,
			},
		}
		printRaw(wrapped)
		updateNoticeEmitted.Store(true)
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
