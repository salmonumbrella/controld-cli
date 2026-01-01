package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

	"github.com/salmonumbrella/controld-cli/internal/auth"
	"github.com/salmonumbrella/controld-cli/internal/outfmt"
	"github.com/salmonumbrella/controld-cli/internal/secrets"
	"github.com/salmonumbrella/controld-cli/internal/ui"
)

func newAuthCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Manage authentication",
	}
	cmd.AddCommand(newAuthLoginCmd())
	cmd.AddCommand(newAuthLogoutCmd())
	cmd.AddCommand(newAuthListCmd())
	cmd.AddCommand(newAuthStatusCmd())
	return cmd
}

func newAuthLoginCmd() *cobra.Command {
	var name string
	var token string
	var noBrowser bool

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with ControlD API",
		Long: `Authenticate with ControlD API.

By default, opens a browser window for interactive authentication.
Use --no-browser for terminal-only authentication.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			u := ui.FromContext(cmd.Context())

			store, err := secrets.OpenDefault()
			if err != nil {
				return fmt.Errorf("failed to open keyring: %w", err)
			}

			// Use browser flow if no flags provided and we're in a terminal
			useBrowser := !noBrowser && token == "" && term.IsTerminal(int(syscall.Stdin))

			if useBrowser {
				u.Info("Opening browser for authentication...")
				fmt.Fprintln(os.Stderr, "If the browser doesn't open, navigate to the URL shown.")
				fmt.Fprintln(os.Stderr)

				server, err := auth.NewSetupServer(store)
				if err != nil {
					return fmt.Errorf("failed to start auth server: %w", err)
				}

				result, err := server.Start(cmd.Context())
				if err != nil {
					return fmt.Errorf("authentication failed: %w", err)
				}

				u.Success(fmt.Sprintf("Authenticated as %s", result.AccountName))
				return nil
			}

			// Terminal flow
			if name == "" {
				name = "default"
			}

			if token == "" {
				fmt.Fprint(os.Stderr, "API Token: ")
				if term.IsTerminal(int(syscall.Stdin)) {
					tokenBytes, err := term.ReadPassword(int(syscall.Stdin))
					if err != nil {
						return fmt.Errorf("failed to read token: %w", err)
					}
					fmt.Fprintln(os.Stderr)
					token = string(tokenBytes)
				} else {
					reader := bufio.NewReader(os.Stdin)
					line, err := reader.ReadString('\n')
					if err != nil {
						return fmt.Errorf("failed to read token: %w", err)
					}
					token = strings.TrimSpace(line)
				}
			}

			if token == "" {
				return fmt.Errorf("token is required")
			}

			if err := store.Set(name, token); err != nil {
				return fmt.Errorf("failed to save credentials: %w", err)
			}

			u.Success(fmt.Sprintf("Authenticated as %s", name))
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Account name for terminal auth (default: default)")
	cmd.Flags().StringVar(&token, "token", "", "API token (skips interactive prompt)")
	cmd.Flags().BoolVar(&noBrowser, "no-browser", false, "Use terminal-only authentication")
	return cmd
}

func newAuthLogoutCmd() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Remove stored credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			u := ui.FromContext(cmd.Context())

			if name == "" {
				name = "default"
			}

			store, err := secrets.OpenDefault()
			if err != nil {
				return fmt.Errorf("failed to open keyring: %w", err)
			}

			if err := store.Delete(name); err != nil {
				return fmt.Errorf("failed to delete credentials: %w", err)
			}

			u.Success(fmt.Sprintf("Logged out from %s", name))
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Account name (default: default)")
	return cmd
}

func newAuthListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List stored accounts",
		RunE: func(cmd *cobra.Command, args []string) error {
			store, err := secrets.OpenDefault()
			if err != nil {
				return fmt.Errorf("failed to open keyring: %w", err)
			}

			creds, err := store.List()
			if err != nil {
				return fmt.Errorf("failed to list accounts: %w", err)
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, creds)
			}

			if len(creds) == 0 {
				fmt.Println("No accounts configured. Run: controld auth login")
				return nil
			}

			tw := outfmt.NewTabWriter(os.Stdout)
			_, _ = fmt.Fprintln(tw, "NAME\tCREATED")
			for _, c := range creds {
				_, _ = fmt.Fprintf(tw, "%s\t%s\n", c.Name, c.CreatedAt.Format("2006-01-02"))
			}
			return tw.Flush()
		},
	}
}

func newAuthStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current authentication status",
		RunE: func(cmd *cobra.Command, args []string) error {
			u := ui.FromContext(cmd.Context())

			store, err := secrets.OpenDefault()
			if err != nil {
				return fmt.Errorf("failed to open keyring: %w", err)
			}

			creds, err := store.List()
			if err != nil {
				return fmt.Errorf("failed to list accounts: %w", err)
			}

			if len(creds) == 0 {
				u.Warn("Not authenticated. Run: controld auth login")
				return nil
			}

			u.Success(fmt.Sprintf("Authenticated with %d account(s)", len(creds)))
			for _, c := range creds {
				fmt.Printf("  - %s (added %s)\n", c.Name, c.CreatedAt.Format("2006-01-02"))
			}
			return nil
		},
	}
}
