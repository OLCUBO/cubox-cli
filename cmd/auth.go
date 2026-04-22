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

var servers = []struct {
	Domain string
	Label  string
}{
	{"cubox.pro", "cubox.pro"},
	{"cubox.cc", "cubox.cc (international)"},
}

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
	if server == "" {
		fmt.Println("? Select Cubox server:")
		for i, s := range servers {
			fmt.Printf("  %d. %s\n", i+1, s.Label)
		}
		fmt.Print("Enter choice [1]: ")
		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)
		switch line {
		case "", "1":
			server = servers[0].Domain
		case "2":
			server = servers[1].Domain
		default:
			return fmt.Errorf("invalid choice: %s", line)
		}
	}

	if server != "cubox.pro" && server != "cubox.cc" {
		return fmt.Errorf("invalid server %q, must be cubox.pro or cubox.cc", server)
	}

	var token string
	switch {
	case flagTokenStdin:
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("reading token from stdin: %w", err)
		}
		token = extractToken(strings.TrimSpace(strings.SplitN(string(data), "\n", 2)[0]))
	case flagToken != "":
		token = extractToken(flagToken)
	default:
		extensionsURL := fmt.Sprintf("https://%s/web/settings/extensions", server)
		fmt.Printf("\nOpen %s\nCopy your API link and paste it here:\n> ", extensionsURL)
		line, _ := reader.ReadString('\n')
		token = extractToken(strings.TrimSpace(line))
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

func extractToken(input string) string {
	input = strings.TrimSpace(input)
	if input == "" {
		return ""
	}
	u, err := url.Parse(input)
	if err == nil && u.Scheme != "" {
		parts := strings.Split(strings.TrimRight(u.Path, "/"), "/")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
	}
	return input
}
