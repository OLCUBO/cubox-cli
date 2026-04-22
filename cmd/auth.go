package cmd

import (
	"bufio"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/OLCUBO/cubox-cli/internal/client"
	"github.com/OLCUBO/cubox-cli/internal/config"
	"github.com/spf13/cobra"
)

var (
	flagServer     string
	flagToken      string
	flagTokenStdin bool
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to Cubox",
	Long: `Log in to Cubox by selecting a server and providing your API key.

Interactive mode (for humans):
  cubox-cli auth login

Non-interactive modes (recommended for Agents / CI):

  # Preferred: transient environment variables (no token is persisted on disk)
  export CUBOX_SERVER=cubox.pro
  export CUBOX_TOKEN=...
  cubox-cli folder list

  # Persisted login, token piped via stdin (keeps it out of argv/ps/history)
  printf '%s' "$TOKEN" | cubox-cli auth login --server cubox.pro --token-stdin

The --token flag is still accepted for backwards compatibility, but note that
passing a token on the command line exposes it to shell history and the
process list; prefer CUBOX_TOKEN or --token-stdin in automated environments.`,
	RunE: runLogin,
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current login status",
	RunE:  runStatus,
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out and remove saved credentials",
	RunE:  runLogout,
}

func init() {
	loginCmd.Flags().StringVar(&flagServer, "server", "", "server domain: cubox.pro or cubox.cc")
	loginCmd.Flags().StringVar(&flagToken, "token", "", "API token on argv (leaks to shell history/ps; prefer CUBOX_TOKEN or --token-stdin)")
	loginCmd.Flags().BoolVar(&flagTokenStdin, "token-stdin", false, "read API token from stdin (first line)")

	authCmd.AddCommand(loginCmd, statusCmd, logoutCmd)
	rootCmd.AddCommand(authCmd)
}

func runLogin(cmd *cobra.Command, args []string) error {
	if flagTokenStdin && flagToken != "" {
		return fmt.Errorf("--token and --token-stdin are mutually exclusive")
	}

	reader := bufio.NewReader(os.Stdin)

	server := flagServer
	var token string

	switch {
	case flagTokenStdin:
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("reading token from stdin: %w", err)
		}
		firstLine := strings.SplitN(strings.TrimSpace(string(data)), "\n", 2)[0]
		server, token = resolveAPIInput(server, firstLine)
	case flagToken != "":
		server, token = resolveAPIInput(server, flagToken)
	default:
		fmt.Println("Sign in to Cubox and paste your API link below.")
		fmt.Println()
		fmt.Println("  1. Please sign in to the Cubox web associated with your account:")
		fmt.Println("     - For international .cc users: https://cubox.cc/web/settings/extensions")
		fmt.Println("     - For .pro users:              https://cubox.pro/web/settings/extensions")
		fmt.Println("  2. Go to Extensions, locate the API Extension, enable it, and copy")
		fmt.Println("     your unique link (e.g. https://cubox.pro/c/api/save/abcd12345).")
		fmt.Println("  3. Paste the link here.")
		fmt.Println()
		fmt.Print("> ")
		line, _ := reader.ReadString('\n')
		server, token = resolveAPIInput(server, line)
	}

	if server == "" {
		return fmt.Errorf("could not determine server. Paste the full API link (e.g. https://cubox.pro/c/api/save/...) or pass --server cubox.pro|cubox.cc")
	}
	if server != "cubox.pro" && server != "cubox.cc" {
		return fmt.Errorf("invalid server %q, must be cubox.pro or cubox.cc", server)
	}
	if token == "" {
		return fmt.Errorf("empty token provided")
	}

	cfg := &config.Config{
		Server: server,
		Token:  token,
	}

	c := client.New(cfg.BaseURL(), cfg.Token)
	if _, err := c.ListFolders(); err != nil {
		return fmt.Errorf("login verification failed: %w\nPlease check your token and try again", err)
	}

	if err := config.Save(cfg); err != nil {
		return err
	}

	fmt.Printf("Logged in to %s successfully.\n", server)
	return nil
}

func runStatus(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	masked := cfg.Token
	if len(masked) > 4 {
		masked = masked[:4] + strings.Repeat("*", len(masked)-4)
	}

	fmt.Printf("Server: %s\n", cfg.Server)
	fmt.Printf("Token:  %s\n", masked)

	c := client.New(cfg.BaseURL(), cfg.Token)
	if _, err := c.ListFolders(); err != nil {
		fmt.Println("Status: connection failed -", err)
	} else {
		fmt.Println("Status: connected")
	}
	return nil
}

func runLogout(cmd *cobra.Command, args []string) error {
	if err := config.Remove(); err != nil {
		return err
	}
	fmt.Println("Logged out successfully.")
	return nil
}

// resolveAPIInput interprets raw as either a bare API token or a full API
// link URL ("https://cubox.pro/c/api/save/abcd"). When raw parses as a URL,
// its host auto-populates server unless the caller already set one
// explicitly (for example via --server). Returns the (possibly updated)
// server and the extracted token.
func resolveAPIInput(server, raw string) (string, string) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return server, ""
	}
	if u, err := url.Parse(raw); err == nil && u.Scheme != "" && u.Host != "" {
		parts := strings.Split(strings.TrimRight(u.Path, "/"), "/")
		token := ""
		if len(parts) > 0 {
			token = parts[len(parts)-1]
		}
		if server == "" {
			server = strings.ToLower(u.Hostname())
		}
		return server, token
	}
	return server, raw
}
