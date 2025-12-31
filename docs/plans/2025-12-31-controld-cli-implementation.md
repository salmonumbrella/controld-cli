# ControlD CLI Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Build a full-featured CLI for the ControlD DNS management API, following airwallex-cli design patterns.

**Architecture:** Thin wrapper around controld-go library. Cobra for CLI, keyring for secrets, context-based dependency injection for testability.

**Tech Stack:** Go 1.21+, github.com/baptistecdr/controld-go, github.com/spf13/cobra, github.com/99designs/keyring

---

## Phase 1: Project Scaffold

### Task 1: Initialize Go Module

**Files:**
- Create: `go.mod`
- Create: `go.sum`

**Step 1: Initialize module**

Run: `go mod init github.com/salmonumbrella/controld-cli`
Expected: Creates go.mod

**Step 2: Add dependencies**

Run:
```bash
go get github.com/baptistecdr/controld-go@latest
go get github.com/spf13/cobra@latest
go get github.com/99designs/keyring@latest
go mod tidy
```
Expected: Downloads deps, updates go.mod and go.sum

**Step 3: Commit**

```bash
git add go.mod go.sum
git commit -m "chore: initialize go module with dependencies"
```

---

### Task 2: Create Config Package

**Files:**
- Create: `internal/config/paths.go`

**Step 1: Write paths.go**

```go
package config

const (
	AppName    = "controld-cli"
	EnvPrefix  = "CONTROLD"
	EnvToken   = "CONTROLD_API_TOKEN"
	EnvOutput  = "CONTROLD_OUTPUT"
	EnvColor   = "CONTROLD_COLOR"
)
```

**Step 2: Commit**

```bash
git add internal/config/paths.go
git commit -m "feat(config): add app name and env var constants"
```

---

### Task 3: Create Debug Package

**Files:**
- Create: `internal/debug/debug.go`

**Step 1: Write debug.go**

```go
package debug

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey struct{}

func SetupLogger(enabled bool) {
	level := slog.LevelInfo
	if enabled {
		level = slog.LevelDebug
	}
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	slog.SetDefault(slog.New(handler))
}

func WithDebug(ctx context.Context, enabled bool) context.Context {
	return context.WithValue(ctx, ctxKey{}, enabled)
}

func IsDebug(ctx context.Context) bool {
	v, _ := ctx.Value(ctxKey{}).(bool)
	return v
}
```

**Step 2: Commit**

```bash
git add internal/debug/debug.go
git commit -m "feat(debug): add debug logging setup"
```

---

### Task 4: Create UI Package

**Files:**
- Create: `internal/ui/ui.go`

**Step 1: Write ui.go**

```go
package ui

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/term"
)

type ctxKey struct{}

type UI struct {
	color string
}

func New(color string) *UI {
	return &UI{color: color}
}

func WithUI(ctx context.Context, u *UI) context.Context {
	return context.WithValue(ctx, ctxKey{}, u)
}

func FromContext(ctx context.Context) *UI {
	u, _ := ctx.Value(ctxKey{}).(*UI)
	if u == nil {
		return New("auto")
	}
	return u
}

func (u *UI) useColor() bool {
	switch u.color {
	case "always":
		return true
	case "never":
		return false
	default:
		return term.IsTerminal(int(os.Stdout.Fd()))
	}
}

func (u *UI) Success(msg string) {
	if u.useColor() {
		fmt.Fprintf(os.Stderr, "\033[32m✓\033[0m %s\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "✓ %s\n", msg)
	}
}

func (u *UI) Error(msg string) {
	if u.useColor() {
		fmt.Fprintf(os.Stderr, "\033[31m✗\033[0m %s\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "✗ %s\n", msg)
	}
}

func (u *UI) Info(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", msg)
}

func (u *UI) Warn(msg string) {
	if u.useColor() {
		fmt.Fprintf(os.Stderr, "\033[33m!\033[0m %s\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "! %s\n", msg)
	}
}
```

**Step 2: Run go mod tidy**

Run: `go mod tidy`

**Step 3: Commit**

```bash
git add internal/ui/ui.go go.mod go.sum
git commit -m "feat(ui): add colored terminal output"
```

---

### Task 5: Create Output Format Package

**Files:**
- Create: `internal/outfmt/outfmt.go`

**Step 1: Write outfmt.go**

```go
package outfmt

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"text/tabwriter"
)

type formatKey struct{}
type queryKey struct{}
type yesKey struct{}
type limitKey struct{}
type sortByKey struct{}
type descKey struct{}

func WithFormat(ctx context.Context, format string) context.Context {
	return context.WithValue(ctx, formatKey{}, format)
}

func WithQuery(ctx context.Context, query string) context.Context {
	return context.WithValue(ctx, queryKey{}, query)
}

func WithYes(ctx context.Context, yes bool) context.Context {
	return context.WithValue(ctx, yesKey{}, yes)
}

func WithLimit(ctx context.Context, limit int) context.Context {
	return context.WithValue(ctx, limitKey{}, limit)
}

func WithSortBy(ctx context.Context, sortBy string) context.Context {
	return context.WithValue(ctx, sortByKey{}, sortBy)
}

func WithDesc(ctx context.Context, desc bool) context.Context {
	return context.WithValue(ctx, descKey{}, desc)
}

func IsJSON(ctx context.Context) bool {
	format, _ := ctx.Value(formatKey{}).(string)
	return format == "json"
}

func GetYes(ctx context.Context) bool {
	yes, _ := ctx.Value(yesKey{}).(bool)
	return yes
}

func GetLimit(ctx context.Context) int {
	limit, _ := ctx.Value(limitKey{}).(int)
	return limit
}

func WriteJSON(w io.Writer, v interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

type Formatter struct {
	w      io.Writer
	format string
	query  string
}

func FromContext(ctx context.Context) *Formatter {
	format, _ := ctx.Value(formatKey{}).(string)
	query, _ := ctx.Value(queryKey{}).(string)
	return &Formatter{w: os.Stdout, format: format, query: query}
}

func (f *Formatter) Output(v interface{}) error {
	return WriteJSON(f.w, v)
}

func NewTabWriter(w io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
}
```

**Step 2: Commit**

```bash
git add internal/outfmt/outfmt.go
git commit -m "feat(outfmt): add output formatting with JSON support"
```

---

### Task 6: Create Root Command

**Files:**
- Create: `internal/cmd/root.go`

**Step 1: Write root.go**

```go
package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/controld-cli/internal/config"
	"github.com/salmonumbrella/controld-cli/internal/debug"
	"github.com/salmonumbrella/controld-cli/internal/outfmt"
	"github.com/salmonumbrella/controld-cli/internal/ui"
)

type rootFlags struct {
	Token   string
	Account string
	Output  string
	Color   string
	Debug   bool
	Query   string
	Yes     bool
	Limit   int
	SortBy  string
	Desc    bool
}

var flags rootFlags

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "controld",
		Short:        "ControlD CLI for DNS management",
		Long:         "A command-line interface for the ControlD DNS management API.",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if flags.Desc && flags.SortBy == "" {
				return fmt.Errorf("--desc requires --sort-by to be specified")
			}

			debug.SetupLogger(flags.Debug)
			ctx := debug.WithDebug(cmd.Context(), flags.Debug)

			u := ui.New(flags.Color)
			ctx = ui.WithUI(ctx, u)

			ctx = outfmt.WithFormat(ctx, flags.Output)
			ctx = outfmt.WithQuery(ctx, flags.Query)
			ctx = outfmt.WithYes(ctx, flags.Yes)
			ctx = outfmt.WithLimit(ctx, flags.Limit)
			ctx = outfmt.WithSortBy(ctx, flags.SortBy)
			ctx = outfmt.WithDesc(ctx, flags.Desc)

			cmd.SetContext(ctx)
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&flags.Token, "token", "", "API token (overrides keyring and env)")
	cmd.PersistentFlags().StringVar(&flags.Account, "account", os.Getenv(config.EnvToken), "Account name from keyring")
	cmd.PersistentFlags().StringVar(&flags.Output, "output", getEnvOrDefault(config.EnvOutput, "text"), "Output format: text|json")
	cmd.PersistentFlags().StringVar(&flags.Color, "color", getEnvOrDefault(config.EnvColor, "auto"), "Color output: auto|always|never")
	cmd.PersistentFlags().BoolVar(&flags.Debug, "debug", false, "Enable debug output")
	cmd.PersistentFlags().StringVar(&flags.Query, "query", "", "JQ filter expression for JSON output")

	cmd.PersistentFlags().BoolVarP(&flags.Yes, "yes", "y", false, "Skip confirmation prompts")
	cmd.PersistentFlags().IntVar(&flags.Limit, "limit", 0, "Limit number of results")
	cmd.PersistentFlags().StringVar(&flags.SortBy, "sort-by", "", "Field name to sort by")
	cmd.PersistentFlags().BoolVar(&flags.Desc, "desc", false, "Sort descending")

	cmd.AddCommand(newVersionCmd())

	return cmd
}

func getEnvOrDefault(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func Execute(args []string) error {
	cmd := NewRootCmd()
	cmd.SetArgs(args)
	return cmd.Execute()
}

func ExecuteContext(ctx context.Context, args []string) error {
	cmd := NewRootCmd()
	cmd.SetArgs(args)
	return cmd.ExecuteContext(ctx)
}
```

**Step 2: Commit**

```bash
git add internal/cmd/root.go
git commit -m "feat(cmd): add root command with global flags"
```

---

### Task 7: Create Version Command

**Files:**
- Create: `internal/cmd/version.go`

**Step 1: Write version.go**

```go
package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/controld-cli/internal/outfmt"
)

var (
	Version   = "dev"
	Commit    = "none"
	BuildDate = "unknown"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			info := map[string]string{
				"version":    Version,
				"commit":     Commit,
				"build_date": BuildDate,
				"go_version": runtime.Version(),
				"os":         runtime.GOOS,
				"arch":       runtime.GOARCH,
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(cmd.OutOrStdout(), info)
			}

			fmt.Printf("controld %s\n", Version)
			fmt.Printf("  commit:     %s\n", Commit)
			fmt.Printf("  built:      %s\n", BuildDate)
			fmt.Printf("  go version: %s\n", runtime.Version())
			fmt.Printf("  platform:   %s/%s\n", runtime.GOOS, runtime.GOARCH)
			return nil
		},
	}
}
```

**Step 2: Commit**

```bash
git add internal/cmd/version.go
git commit -m "feat(cmd): add version command"
```

---

### Task 8: Create Main Entry Point

**Files:**
- Create: `cmd/controld/main.go`

**Step 1: Write main.go**

```go
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/salmonumbrella/controld-cli/internal/cmd"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	if err := cmd.ExecuteContext(ctx, os.Args[1:]); err != nil {
		os.Exit(1)
	}
}
```

**Step 2: Run go mod tidy and build**

Run:
```bash
go mod tidy
go build -o controld ./cmd/controld
./controld version
```
Expected: Prints version info

**Step 3: Commit**

```bash
git add cmd/controld/main.go go.mod go.sum
git commit -m "feat: add main entry point"
```

---

## Phase 2: Authentication

### Task 9: Create Secrets Store

**Files:**
- Create: `internal/secrets/store.go`

**Step 1: Write store.go**

```go
package secrets

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/99designs/keyring"

	"github.com/salmonumbrella/controld-cli/internal/config"
)

type Store interface {
	Set(name string, token string) error
	Get(name string) (Credentials, error)
	Delete(name string) error
	List() ([]Credentials, error)
	Keys() ([]string, error)
}

type Credentials struct {
	Name      string    `json:"name"`
	Token     string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
}

type storedCredentials struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

type KeyringStore struct {
	ring keyring.Keyring
}

func OpenDefault() (Store, error) {
	ring, err := keyring.Open(keyring.Config{
		ServiceName: config.AppName,
	})
	if err != nil {
		return nil, err
	}
	return &KeyringStore{ring: ring}, nil
}

func (s *KeyringStore) Keys() ([]string, error) {
	return s.ring.Keys()
}

func (s *KeyringStore) Set(name string, token string) error {
	name = normalize(name)
	if name == "" {
		return fmt.Errorf("missing account name")
	}
	if token == "" {
		return fmt.Errorf("missing token")
	}

	payload, err := json.Marshal(storedCredentials{
		Token:     token,
		CreatedAt: time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	return s.ring.Set(keyring.Item{
		Key:  credentialKey(name),
		Data: payload,
	})
}

func (s *KeyringStore) Get(name string) (Credentials, error) {
	name = normalize(name)
	if name == "" {
		return Credentials{}, fmt.Errorf("missing account name")
	}

	item, err := s.ring.Get(credentialKey(name))
	if err != nil {
		return Credentials{}, err
	}

	var stored storedCredentials
	if err := json.Unmarshal(item.Data, &stored); err != nil {
		return Credentials{}, err
	}

	return Credentials{
		Name:      name,
		Token:     stored.Token,
		CreatedAt: stored.CreatedAt,
	}, nil
}

func (s *KeyringStore) Delete(name string) error {
	name = normalize(name)
	if name == "" {
		return fmt.Errorf("missing account name")
	}
	return s.ring.Remove(credentialKey(name))
}

func (s *KeyringStore) List() ([]Credentials, error) {
	keys, err := s.Keys()
	if err != nil {
		return nil, err
	}

	var out []Credentials
	for _, k := range keys {
		name, ok := parseCredentialKey(k)
		if !ok {
			continue
		}
		creds, err := s.Get(name)
		if err != nil {
			continue
		}
		out = append(out, creds)
	}
	return out, nil
}

func credentialKey(name string) string {
	return fmt.Sprintf("account:%s", name)
}

func parseCredentialKey(k string) (string, bool) {
	const prefix = "account:"
	if !strings.HasPrefix(k, prefix) {
		return "", false
	}
	return strings.TrimPrefix(k, prefix), true
}

func normalize(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
```

**Step 2: Commit**

```bash
git add internal/secrets/store.go
git commit -m "feat(secrets): add keyring-based credential storage"
```

---

### Task 10: Create Auth Commands

**Files:**
- Create: `internal/cmd/auth.go`

**Step 1: Write auth.go**

```go
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"

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

	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with ControlD API",
		RunE: func(cmd *cobra.Command, args []string) error {
			u := ui.FromContext(cmd.Context())

			if name == "" {
				name = "default"
			}

			if token == "" {
				fmt.Print("API Token: ")
				if term.IsTerminal(int(syscall.Stdin)) {
					tokenBytes, err := term.ReadPassword(int(syscall.Stdin))
					if err != nil {
						return fmt.Errorf("failed to read token: %w", err)
					}
					fmt.Println()
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

			store, err := secrets.OpenDefault()
			if err != nil {
				return fmt.Errorf("failed to open keyring: %w", err)
			}

			if err := store.Set(name, token); err != nil {
				return fmt.Errorf("failed to save credentials: %w", err)
			}

			u.Success(fmt.Sprintf("Authenticated as %s", name))
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Account name (default: default)")
	cmd.Flags().StringVar(&token, "token", "", "API token")
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
			fmt.Fprintln(tw, "NAME\tCREATED")
			for _, c := range creds {
				fmt.Fprintf(tw, "%s\t%s\n", c.Name, c.CreatedAt.Format("2006-01-02"))
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
```

**Step 2: Update root.go to add auth command**

Add after `cmd.AddCommand(newVersionCmd())`:
```go
cmd.AddCommand(newAuthCmd())
```

**Step 3: Run go mod tidy and test**

Run:
```bash
go mod tidy
go build -o controld ./cmd/controld
./controld auth --help
```
Expected: Shows auth subcommands

**Step 4: Commit**

```bash
git add internal/cmd/auth.go internal/cmd/root.go go.mod go.sum
git commit -m "feat(auth): add login/logout/list/status commands"
```

---

## Phase 3: API Client

### Task 11: Create API Client Wrapper

**Files:**
- Create: `internal/api/client.go`

**Step 1: Write client.go**

```go
package api

import (
	"context"
	"fmt"
	"os"

	controld "github.com/baptistecdr/controld-go"

	"github.com/salmonumbrella/controld-cli/internal/config"
	"github.com/salmonumbrella/controld-cli/internal/debug"
	"github.com/salmonumbrella/controld-cli/internal/secrets"
)

type ClientConfig struct {
	Token   string
	Account string
}

func NewClient(ctx context.Context, cfg ClientConfig) (*controld.API, error) {
	token := resolveToken(cfg)
	if token == "" {
		return nil, fmt.Errorf("no API token found. Set %s or run: controld auth login", config.EnvToken)
	}

	opts := []controld.Option{}
	if debug.IsDebug(ctx) {
		opts = append(opts, controld.Debug(true))
	}

	return controld.New(token, opts...)
}

func resolveToken(cfg ClientConfig) string {
	// 1. Explicit token flag
	if cfg.Token != "" {
		return cfg.Token
	}

	// 2. Environment variable
	if token := os.Getenv(config.EnvToken); token != "" {
		return token
	}

	// 3. Keyring
	store, err := secrets.OpenDefault()
	if err != nil {
		return ""
	}

	// If account specified, use it
	if cfg.Account != "" {
		creds, err := store.Get(cfg.Account)
		if err != nil {
			return ""
		}
		return creds.Token
	}

	// Auto-select if only one account
	creds, err := store.List()
	if err != nil || len(creds) != 1 {
		return ""
	}
	return creds[0].Token
}
```

**Step 2: Commit**

```bash
git add internal/api/client.go
git commit -m "feat(api): add client wrapper with token resolution"
```

---

### Task 12: Add Client Helper to Commands

**Files:**
- Modify: `internal/cmd/root.go`

**Step 1: Add getClient function to root.go**

Add at the end of root.go:
```go
import (
	// ... existing imports
	controld "github.com/baptistecdr/controld-go"
	"github.com/salmonumbrella/controld-cli/internal/api"
)

func getClient(ctx context.Context) (*controld.API, error) {
	return api.NewClient(ctx, api.ClientConfig{
		Token:   flags.Token,
		Account: flags.Account,
	})
}
```

**Step 2: Commit**

```bash
git add internal/cmd/root.go
git commit -m "feat(cmd): add getClient helper for commands"
```

---

## Phase 4: Core Commands

### Task 13: Create Users Command

**Files:**
- Create: `internal/cmd/users.go`

**Step 1: Write users.go**

```go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/controld-cli/internal/outfmt"
)

func newUsersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "users",
		Short: "Show account information",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			user, err := client.ListUser(cmd.Context())
			if err != nil {
				return err
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, user)
			}

			tw := outfmt.NewTabWriter(os.Stdout)
			fmt.Fprintf(tw, "email\t%s\n", user.Email)
			fmt.Fprintf(tw, "status\t%v\n", user.Status)
			fmt.Fprintf(tw, "resolver_ip\t%s\n", user.ResolverIP)
			fmt.Fprintf(tw, "stats_endpoint\t%s\n", user.StatsEndpoint)
			fmt.Fprintf(tw, "2fa\t%v\n", user.Twofa)
			return tw.Flush()
		},
	}
}
```

**Step 2: Add to root.go**

Add after auth command: `cmd.AddCommand(newUsersCmd())`

**Step 3: Commit**

```bash
git add internal/cmd/users.go internal/cmd/root.go
git commit -m "feat(cmd): add users command for account info"
```

---

### Task 14: Create Devices Commands

**Files:**
- Create: `internal/cmd/devices.go`

**Step 1: Write devices.go**

```go
package cmd

import (
	"fmt"
	"os"

	controld "github.com/baptistecdr/controld-go"
	"github.com/spf13/cobra"

	"github.com/salmonumbrella/controld-cli/internal/outfmt"
	"github.com/salmonumbrella/controld-cli/internal/ui"
)

func newDevicesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "devices",
		Short: "Manage DNS resolver devices",
	}
	cmd.AddCommand(newDevicesListCmd())
	cmd.AddCommand(newDevicesGetCmd())
	cmd.AddCommand(newDevicesCreateCmd())
	cmd.AddCommand(newDevicesDeleteCmd())
	cmd.AddCommand(newDevicesTypesCmd())
	return cmd
}

func newDevicesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all devices",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			devices, err := client.ListDevices(cmd.Context())
			if err != nil {
				return err
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, devices)
			}

			if len(devices) == 0 {
				fmt.Println("No devices found")
				return nil
			}

			tw := outfmt.NewTabWriter(os.Stdout)
			fmt.Fprintln(tw, "DEVICE_ID\tNAME\tSTATUS\tPROFILE")
			for _, d := range devices {
				status := "pending"
				switch d.Status {
				case controld.Active:
					status = "active"
				case controld.SoftDisabled:
					status = "soft-disabled"
				case controld.HardDisabled:
					status = "hard-disabled"
				}
				fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", d.DeviceID, d.Name, status, d.Profile.Name)
			}
			return tw.Flush()
		},
	}
}

func newDevicesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <device-id>",
		Short: "Get device details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			devices, err := client.ListDevices(cmd.Context())
			if err != nil {
				return err
			}

			deviceID := args[0]
			var device *controld.Device
			for i := range devices {
				if devices[i].DeviceID == deviceID {
					device = &devices[i]
					break
				}
			}

			if device == nil {
				return fmt.Errorf("device not found: %s", deviceID)
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, device)
			}

			tw := outfmt.NewTabWriter(os.Stdout)
			fmt.Fprintf(tw, "device_id\t%s\n", device.DeviceID)
			fmt.Fprintf(tw, "name\t%s\n", device.Name)
			fmt.Fprintf(tw, "status\t%d\n", device.Status)
			fmt.Fprintf(tw, "profile\t%s\n", device.Profile.Name)
			fmt.Fprintf(tw, "doh\t%s\n", device.Resolvers.DoH)
			fmt.Fprintf(tw, "dot\t%s\n", device.Resolvers.DoT)
			if device.Icon != nil {
				fmt.Fprintf(tw, "icon\t%s\n", *device.Icon)
			}
			return tw.Flush()
		},
	}
}

func newDevicesCreateCmd() *cobra.Command {
	var name string
	var profileID string
	var icon string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new device",
		RunE: func(cmd *cobra.Command, args []string) error {
			u := ui.FromContext(cmd.Context())
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			device, err := client.CreateDevice(cmd.Context(), controld.CreateDeviceParams{
				Name:      name,
				ProfileID: profileID,
				Icon:      controld.IconName(icon),
			})
			if err != nil {
				return err
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, device)
			}

			u.Success(fmt.Sprintf("Created device: %s (%s)", device.Name, device.DeviceID))
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Device name (required)")
	cmd.Flags().StringVar(&profileID, "profile-id", "", "Profile ID (required)")
	cmd.Flags().StringVar(&icon, "icon", "router", "Device icon")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("profile-id")
	return cmd
}

func newDevicesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <device-id>",
		Short: "Delete a device",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			u := ui.FromContext(cmd.Context())
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			if !outfmt.GetYes(cmd.Context()) {
				fmt.Printf("Delete device %s? [y/N]: ", args[0])
				var confirm string
				fmt.Scanln(&confirm)
				if confirm != "y" && confirm != "Y" {
					fmt.Println("Cancelled")
					return nil
				}
			}

			_, err = client.DeleteDevice(cmd.Context(), controld.DeleteDeviceParams{
				DeviceID: args[0],
			})
			if err != nil {
				return err
			}

			u.Success(fmt.Sprintf("Deleted device: %s", args[0]))
			return nil
		},
	}
}

func newDevicesTypesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "types",
		Short: "List available device types",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			types, err := client.ListDeviceType(cmd.Context())
			if err != nil {
				return err
			}

			return outfmt.WriteJSON(os.Stdout, types)
		},
	}
}
```

**Step 2: Add to root.go**

Add: `cmd.AddCommand(newDevicesCmd())`

**Step 3: Build and test**

Run:
```bash
go build -o controld ./cmd/controld
./controld devices --help
```

**Step 4: Commit**

```bash
git add internal/cmd/devices.go internal/cmd/root.go
git commit -m "feat(cmd): add devices commands"
```

---

### Task 15: Create Profiles Commands

**Files:**
- Create: `internal/cmd/profiles.go`

**Step 1: Write profiles.go**

```go
package cmd

import (
	"fmt"
	"os"

	controld "github.com/baptistecdr/controld-go"
	"github.com/spf13/cobra"

	"github.com/salmonumbrella/controld-cli/internal/outfmt"
	"github.com/salmonumbrella/controld-cli/internal/ui"
)

func newProfilesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profiles",
		Short: "Manage DNS filtering profiles",
	}
	cmd.AddCommand(newProfilesListCmd())
	cmd.AddCommand(newProfilesCreateCmd())
	cmd.AddCommand(newProfilesDeleteCmd())
	return cmd
}

func newProfilesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all profiles",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			profiles, err := client.ListProfiles(cmd.Context())
			if err != nil {
				return err
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, profiles)
			}

			if len(profiles) == 0 {
				fmt.Println("No profiles found")
				return nil
			}

			tw := outfmt.NewTabWriter(os.Stdout)
			fmt.Fprintln(tw, "PROFILE_ID\tNAME\tUPDATED")
			for _, p := range profiles {
				fmt.Fprintf(tw, "%s\t%s\t%s\n", p.PK, p.Name, p.Updated.Format("2006-01-02"))
			}
			return tw.Flush()
		},
	}
}

func newProfilesCreateCmd() *cobra.Command {
	var name string
	var cloneFrom string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			u := ui.FromContext(cmd.Context())
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			params := controld.CreateProfileParams{Name: name}
			if cloneFrom != "" {
				params.CloneProfileID = &cloneFrom
			}

			profiles, err := client.CreateProfile(cmd.Context(), params)
			if err != nil {
				return err
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, profiles)
			}

			if len(profiles) > 0 {
				u.Success(fmt.Sprintf("Created profile: %s (%s)", profiles[0].Name, profiles[0].PK))
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Profile name (required)")
	cmd.Flags().StringVar(&cloneFrom, "clone-from", "", "Clone from existing profile ID")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func newProfilesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <profile-id>",
		Short: "Delete a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			u := ui.FromContext(cmd.Context())
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			if !outfmt.GetYes(cmd.Context()) {
				fmt.Printf("Delete profile %s? [y/N]: ", args[0])
				var confirm string
				fmt.Scanln(&confirm)
				if confirm != "y" && confirm != "Y" {
					fmt.Println("Cancelled")
					return nil
				}
			}

			_, err = client.DeleteProfile(cmd.Context(), controld.DeleteProfileParams{
				ProfileID: args[0],
			})
			if err != nil {
				return err
			}

			u.Success(fmt.Sprintf("Deleted profile: %s", args[0]))
			return nil
		},
	}
}
```

**Step 2: Add to root.go**

Add: `cmd.AddCommand(newProfilesCmd())`

**Step 3: Commit**

```bash
git add internal/cmd/profiles.go internal/cmd/root.go
git commit -m "feat(cmd): add profiles commands"
```

---

## Phase 5: Profile Sub-Commands

### Task 16: Create Profiles Rules Commands

**Files:**
- Create: `internal/cmd/profiles_rules.go`

**Step 1: Write profiles_rules.go**

```go
package cmd

import (
	"fmt"
	"os"

	controld "github.com/baptistecdr/controld-go"
	"github.com/spf13/cobra"

	"github.com/salmonumbrella/controld-cli/internal/outfmt"
	"github.com/salmonumbrella/controld-cli/internal/ui"
)

func newProfilesRulesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rules",
		Short: "Manage custom domain rules",
	}
	cmd.AddCommand(newProfilesRulesListCmd())
	cmd.AddCommand(newProfilesRulesCreateCmd())
	cmd.AddCommand(newProfilesRulesDeleteCmd())
	return cmd
}

func newProfilesRulesListCmd() *cobra.Command {
	var folderID string

	cmd := &cobra.Command{
		Use:   "list <profile-id>",
		Short: "List custom rules for a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			if folderID == "" {
				folderID = "0"
			}

			rules, err := client.ListProfileCustomRules(cmd.Context(), controld.ListProfileCustomRulesParams{
				ProfileID: args[0],
				FolderID:  folderID,
			})
			if err != nil {
				return err
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, rules)
			}

			if len(rules) == 0 {
				fmt.Println("No custom rules found")
				return nil
			}

			tw := outfmt.NewTabWriter(os.Stdout)
			fmt.Fprintln(tw, "DOMAIN\tACTION\tSTATUS")
			for _, r := range rules {
				action := actionToString(r.Action.Do)
				status := "disabled"
				if r.Action.Status {
					status = "enabled"
				}
				fmt.Fprintf(tw, "%s\t%s\t%s\n", r.PK, action, status)
			}
			return tw.Flush()
		},
	}

	cmd.Flags().StringVar(&folderID, "folder", "", "Folder ID (default: 0)")
	return cmd
}

func newProfilesRulesCreateCmd() *cobra.Command {
	var do string
	var hostnames []string

	cmd := &cobra.Command{
		Use:   "create <profile-id>",
		Short: "Create custom rules",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			u := ui.FromContext(cmd.Context())
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			doType := stringToAction(do)

			rules, err := client.CreateProfileCustomRule(cmd.Context(), controld.CreateProfileCustomRuleParams{
				ProfileID: args[0],
				Do:        doType,
				Status:    controld.IntBool(true),
				Hostnames: hostnames,
			})
			if err != nil {
				return err
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, rules)
			}

			u.Success(fmt.Sprintf("Created %d rule(s)", len(rules)))
			return nil
		},
	}

	cmd.Flags().StringVar(&do, "action", "block", "Action: block|bypass|spoof|redirect")
	cmd.Flags().StringSliceVar(&hostnames, "hostname", nil, "Hostname(s) to add (required)")
	_ = cmd.MarkFlagRequired("hostname")
	return cmd
}

func newProfilesRulesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <profile-id> <hostname>",
		Short: "Delete a custom rule",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			u := ui.FromContext(cmd.Context())
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			_, err = client.DeleteProfileCustomRule(cmd.Context(), controld.DeleteProfileCustomRuleParams{
				ProfileID: args[0],
				Hostname:  args[1],
			})
			if err != nil {
				return err
			}

			u.Success(fmt.Sprintf("Deleted rule: %s", args[1]))
			return nil
		},
	}
}

func actionToString(do controld.DoType) string {
	switch do {
	case controld.Block:
		return "block"
	case controld.Bypass:
		return "bypass"
	case controld.Spoof:
		return "spoof"
	case controld.Redirect:
		return "redirect"
	default:
		return "unknown"
	}
}

func stringToAction(s string) controld.DoType {
	switch s {
	case "bypass":
		return controld.Bypass
	case "spoof":
		return controld.Spoof
	case "redirect":
		return controld.Redirect
	default:
		return controld.Block
	}
}
```

**Step 2: Add to profiles.go**

In `newProfilesCmd()`, add: `cmd.AddCommand(newProfilesRulesCmd())`

**Step 3: Commit**

```bash
git add internal/cmd/profiles_rules.go internal/cmd/profiles.go
git commit -m "feat(cmd): add profiles rules sub-commands"
```

---

### Task 17: Create Profiles Filters Commands

**Files:**
- Create: `internal/cmd/profiles_filters.go`

**Step 1: Write profiles_filters.go**

```go
package cmd

import (
	"fmt"
	"os"

	controld "github.com/baptistecdr/controld-go"
	"github.com/spf13/cobra"

	"github.com/salmonumbrella/controld-cli/internal/outfmt"
	"github.com/salmonumbrella/controld-cli/internal/ui"
)

func newProfilesFiltersCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "filters",
		Short: "Manage DNS filters",
	}
	cmd.AddCommand(newProfilesFiltersListCmd())
	cmd.AddCommand(newProfilesFiltersEnableCmd())
	cmd.AddCommand(newProfilesFiltersDisableCmd())
	return cmd
}

func newProfilesFiltersListCmd() *cobra.Command {
	var external bool

	cmd := &cobra.Command{
		Use:   "list <profile-id>",
		Short: "List available filters",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			params := controld.ListProfileFiltersParams{ProfileID: args[0]}

			var filters []controld.Filter
			if external {
				filters, err = client.ListProfileExternalFilters(cmd.Context(), params)
			} else {
				filters, err = client.ListProfileNativeFilters(cmd.Context(), params)
			}
			if err != nil {
				return err
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, filters)
			}

			if len(filters) == 0 {
				fmt.Println("No filters found")
				return nil
			}

			tw := outfmt.NewTabWriter(os.Stdout)
			fmt.Fprintln(tw, "FILTER_ID\tNAME\tSTATUS")
			for _, f := range filters {
				status := "disabled"
				if f.Status {
					status = "enabled"
				}
				fmt.Fprintf(tw, "%s\t%s\t%s\n", f.PK, f.Name, status)
			}
			return tw.Flush()
		},
	}

	cmd.Flags().BoolVar(&external, "external", false, "Show external/community filters")
	return cmd
}

func newProfilesFiltersEnableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "enable <profile-id> <filter-id>",
		Short: "Enable a filter",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			u := ui.FromContext(cmd.Context())
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			_, err = client.UpdateProfileFilter(cmd.Context(), controld.UpdateProfileFilterParams{
				ProfileID: args[0],
				Filter:    args[1],
				Status:    controld.IntBool(true),
			})
			if err != nil {
				return err
			}

			u.Success(fmt.Sprintf("Enabled filter: %s", args[1]))
			return nil
		},
	}
}

func newProfilesFiltersDisableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disable <profile-id> <filter-id>",
		Short: "Disable a filter",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			u := ui.FromContext(cmd.Context())
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			_, err = client.UpdateProfileFilter(cmd.Context(), controld.UpdateProfileFilterParams{
				ProfileID: args[0],
				Filter:    args[1],
				Status:    controld.IntBool(false),
			})
			if err != nil {
				return err
			}

			u.Success(fmt.Sprintf("Disabled filter: %s", args[1]))
			return nil
		},
	}
}
```

**Step 2: Add to profiles.go**

Add: `cmd.AddCommand(newProfilesFiltersCmd())`

**Step 3: Commit**

```bash
git add internal/cmd/profiles_filters.go internal/cmd/profiles.go
git commit -m "feat(cmd): add profiles filters sub-commands"
```

---

### Task 18: Create Profiles Services Commands

**Files:**
- Create: `internal/cmd/profiles_services.go`

**Step 1: Write profiles_services.go**

```go
package cmd

import (
	"fmt"
	"os"

	controld "github.com/baptistecdr/controld-go"
	"github.com/spf13/cobra"

	"github.com/salmonumbrella/controld-cli/internal/outfmt"
	"github.com/salmonumbrella/controld-cli/internal/ui"
)

func newProfilesServicesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "services",
		Short: "Manage service blocking/unblocking",
	}
	cmd.AddCommand(newProfilesServicesListCmd())
	cmd.AddCommand(newProfilesServicesSetCmd())
	return cmd
}

func newProfilesServicesListCmd() *cobra.Command {
	var category string

	cmd := &cobra.Command{
		Use:   "list <profile-id>",
		Short: "List services for a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			services, err := client.ListProfileServices(cmd.Context(), controld.ListProfileServicesParams{
				ProfileID: args[0],
			})
			if err != nil {
				return err
			}

			// Filter by category if specified
			if category != "" {
				var filtered []controld.ProfileService
				for _, s := range services {
					if s.Category == category {
						filtered = append(filtered, s)
					}
				}
				services = filtered
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, services)
			}

			if len(services) == 0 {
				fmt.Println("No services found")
				return nil
			}

			tw := outfmt.NewTabWriter(os.Stdout)
			fmt.Fprintln(tw, "SERVICE_ID\tNAME\tCATEGORY\tACTION\tSTATUS")
			for _, s := range services {
				action := actionToString(s.Action.Do)
				status := "disabled"
				if s.Action.Status {
					status = "enabled"
				}
				fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", s.PK, s.Name, s.Category, action, status)
			}
			return tw.Flush()
		},
	}

	cmd.Flags().StringVar(&category, "category", "", "Filter by category")
	return cmd
}

func newProfilesServicesSetCmd() *cobra.Command {
	var action string

	cmd := &cobra.Command{
		Use:   "set <profile-id> <service-id>",
		Short: "Set action for a service",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			u := ui.FromContext(cmd.Context())
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			_, err = client.UpdateProfileService(cmd.Context(), controld.UpdateProfileServiceParams{
				ProfileID: args[0],
				Service:   args[1],
				Do:        stringToAction(action),
				Status:    controld.IntBool(true),
			})
			if err != nil {
				return err
			}

			u.Success(fmt.Sprintf("Set %s to %s", args[1], action))
			return nil
		},
	}

	cmd.Flags().StringVar(&action, "action", "block", "Action: block|bypass|spoof")
	return cmd
}
```

**Step 2: Add to profiles.go**

Add: `cmd.AddCommand(newProfilesServicesCmd())`

**Step 3: Commit**

```bash
git add internal/cmd/profiles_services.go internal/cmd/profiles.go
git commit -m "feat(cmd): add profiles services sub-commands"
```

---

## Phase 6: Reference Commands

### Task 19: Create Services Command

**Files:**
- Create: `internal/cmd/services.go`

**Step 1: Write services.go**

```go
package cmd

import (
	"fmt"
	"os"

	controld "github.com/baptistecdr/controld-go"
	"github.com/spf13/cobra"

	"github.com/salmonumbrella/controld-cli/internal/outfmt"
)

func newServicesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "services",
		Short: "List available services (reference)",
	}
	cmd.AddCommand(newServicesListCmd())
	cmd.AddCommand(newServicesCategoriesCmd())
	return cmd
}

func newServicesListCmd() *cobra.Command {
	var category string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all available services",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			params := controld.ListServicesParams{}
			if category != "" {
				params.Category = &category
			}

			services, err := client.ListServices(cmd.Context(), params)
			if err != nil {
				return err
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, services)
			}

			if len(services) == 0 {
				fmt.Println("No services found")
				return nil
			}

			tw := outfmt.NewTabWriter(os.Stdout)
			fmt.Fprintln(tw, "SERVICE_ID\tNAME\tCATEGORY")
			for _, s := range services {
				fmt.Fprintf(tw, "%s\t%s\t%s\n", s.PK, s.Name, s.Category)
			}
			return tw.Flush()
		},
	}

	cmd.Flags().StringVar(&category, "category", "", "Filter by category")
	return cmd
}

func newServicesCategoriesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "categories",
		Short: "List service categories",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			categories, err := client.ListServiceCategories(cmd.Context())
			if err != nil {
				return err
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, categories)
			}

			tw := outfmt.NewTabWriter(os.Stdout)
			fmt.Fprintln(tw, "CATEGORY_ID\tNAME\tCOUNT")
			for _, c := range categories {
				fmt.Fprintf(tw, "%s\t%s\t%d\n", c.PK, c.Name, c.Count)
			}
			return tw.Flush()
		},
	}
}
```

**Step 2: Add to root.go**

Add: `cmd.AddCommand(newServicesCmd())`

**Step 3: Commit**

```bash
git add internal/cmd/services.go internal/cmd/root.go
git commit -m "feat(cmd): add services reference commands"
```

---

### Task 20: Create Network Command

**Files:**
- Create: `internal/cmd/network.go`

**Step 1: Write network.go**

```go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/controld-cli/internal/outfmt"
)

func newNetworkCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "network",
		Short: "Network status and information",
	}
	cmd.AddCommand(newNetworkStatusCmd())
	cmd.AddCommand(newNetworkIPCmd())
	return cmd
}

func newNetworkStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show ControlD network status",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			network, err := client.ListNetwork(cmd.Context())
			if err != nil {
				return err
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, network)
			}

			tw := outfmt.NewTabWriter(os.Stdout)
			fmt.Fprintln(tw, "POP\tCITY\tCOUNTRY\tAPI\tDNS\tPROXY")
			for _, n := range network {
				fmt.Fprintf(tw, "%s\t%s\t%s\t%d\t%d\t%d\n",
					n.IataCode, n.CityName, n.CountryName,
					n.Status.API, n.Status.DNS, n.Status.Pxy)
			}
			return tw.Flush()
		},
	}
}

func newNetworkIPCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ip",
		Short: "Show your current IP information",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			ip, err := client.ListIP(cmd.Context())
			if err != nil {
				return err
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, ip)
			}

			tw := outfmt.NewTabWriter(os.Stdout)
			fmt.Fprintf(tw, "ip\t%s\n", ip.IP)
			fmt.Fprintf(tw, "type\t%s\n", ip.Type)
			fmt.Fprintf(tw, "country\t%s\n", ip.Country)
			fmt.Fprintf(tw, "org\t%s\n", ip.Org)
			fmt.Fprintf(tw, "pop\t%s\n", ip.Pop)
			return tw.Flush()
		},
	}
}
```

**Step 2: Add to root.go**

Add: `cmd.AddCommand(newNetworkCmd())`

**Step 3: Commit**

```bash
git add internal/cmd/network.go internal/cmd/root.go
git commit -m "feat(cmd): add network status commands"
```

---

### Task 21: Create Access Command

**Files:**
- Create: `internal/cmd/access.go`

**Step 1: Write access.go**

```go
package cmd

import (
	"fmt"
	"net"
	"os"

	controld "github.com/baptistecdr/controld-go"
	"github.com/spf13/cobra"

	"github.com/salmonumbrella/controld-cli/internal/outfmt"
	"github.com/salmonumbrella/controld-cli/internal/ui"
)

func newAccessCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "access",
		Short: "Manage known IPs for devices",
	}
	cmd.AddCommand(newAccessListCmd())
	cmd.AddCommand(newAccessAddCmd())
	cmd.AddCommand(newAccessDeleteCmd())
	return cmd
}

func newAccessListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list <device-id>",
		Short: "List known IPs for a device",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			ips, err := client.ListKnownIPs(cmd.Context(), controld.ListKnownIPsParams{
				DeviceID: args[0],
			})
			if err != nil {
				return err
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, ips)
			}

			if len(ips) == 0 {
				fmt.Println("No known IPs")
				return nil
			}

			tw := outfmt.NewTabWriter(os.Stdout)
			fmt.Fprintln(tw, "IP\tCOUNTRY\tCITY\tISP")
			for _, ip := range ips {
				fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", ip.IP, ip.Country, ip.City, ip.ISP)
			}
			return tw.Flush()
		},
	}
}

func newAccessAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <device-id> <ip>...",
		Short: "Add known IPs to a device",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			u := ui.FromContext(cmd.Context())
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			var ips []net.IP
			for _, ipStr := range args[1:] {
				ip := net.ParseIP(ipStr)
				if ip == nil {
					return fmt.Errorf("invalid IP: %s", ipStr)
				}
				ips = append(ips, ip)
			}

			_, err = client.LearnNewIPs(cmd.Context(), controld.LearnNewIPsParams{
				DeviceID: args[0],
				IPs:      ips,
			})
			if err != nil {
				return err
			}

			u.Success(fmt.Sprintf("Added %d IP(s)", len(ips)))
			return nil
		},
	}
}

func newAccessDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <device-id> <ip>...",
		Short: "Delete known IPs from a device",
		Args:  cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			u := ui.FromContext(cmd.Context())
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			var ips []net.IP
			for _, ipStr := range args[1:] {
				ip := net.ParseIP(ipStr)
				if ip == nil {
					return fmt.Errorf("invalid IP: %s", ipStr)
				}
				ips = append(ips, ip)
			}

			_, err = client.DeleteLearnedIPs(cmd.Context(), controld.DeleteLearnedIPsParams{
				DeviceID: args[0],
				IPs:      ips,
			})
			if err != nil {
				return err
			}

			u.Success(fmt.Sprintf("Deleted %d IP(s)", len(ips)))
			return nil
		},
	}
}
```

**Step 2: Add to root.go**

Add: `cmd.AddCommand(newAccessCmd())`

**Step 3: Commit**

```bash
git add internal/cmd/access.go internal/cmd/root.go
git commit -m "feat(cmd): add access commands for known IPs"
```

---

### Task 22: Create Completion Command

**Files:**
- Create: `internal/cmd/completion.go`

**Step 1: Write completion.go**

```go
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

func newCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "completion [bash|zsh|fish|powershell]",
		Short: "Generate shell completion scripts",
		Long: `Generate shell completion scripts for controld.

To load completions:

Bash:
  $ source <(controld completion bash)

Zsh:
  $ source <(controld completion zsh)

Fish:
  $ controld completion fish | source
`,
		DisableFlagsInUseLine: true,
		ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
		Args:                  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			switch args[0] {
			case "bash":
				return cmd.Root().GenBashCompletion(os.Stdout)
			case "zsh":
				return cmd.Root().GenZshCompletion(os.Stdout)
			case "fish":
				return cmd.Root().GenFishCompletion(os.Stdout, true)
			case "powershell":
				return cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
			}
			return nil
		},
	}
	return cmd
}
```

**Step 2: Add to root.go**

Add: `cmd.AddCommand(newCompletionCmd())`

**Step 3: Final build and test**

Run:
```bash
go mod tidy
go build -o controld ./cmd/controld
./controld --help
```

**Step 4: Commit**

```bash
git add internal/cmd/completion.go internal/cmd/root.go go.mod go.sum
git commit -m "feat(cmd): add shell completion command"
```

---

### Task 23: Push All Changes

**Step 1: Push to GitHub**

Run:
```bash
git push origin main
```

**Step 2: Verify on GitHub**

Visit: https://github.com/salmonumbrella/controld-cli

---

## Summary

This plan implements the ControlD CLI with:

- **Phase 1:** Project scaffold (go.mod, config, debug, ui, outfmt, root command)
- **Phase 2:** Authentication (keyring store, login/logout/list/status)
- **Phase 3:** API client wrapper with token resolution
- **Phase 4:** Core commands (users, devices, profiles)
- **Phase 5:** Profile sub-commands (rules, filters, services)
- **Phase 6:** Reference commands (services, network, access, completion)

Each task is atomic and can be completed in 2-5 minutes.
