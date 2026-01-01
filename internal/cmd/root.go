package cmd

import (
	"context"
	"fmt"
	"os"

	controld "github.com/baptistecdr/controld-go"
	"github.com/spf13/cobra"

	"github.com/salmonumbrella/controld-cli/internal/api"
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
	Yes     bool
}

var flags rootFlags

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "controld",
		Short:        "ControlD CLI for DNS management",
		Long:         "A command-line interface for the ControlD DNS management API.",
		SilenceUsage: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if flags.Output != "" && flags.Output != "json" && flags.Output != "text" {
				return fmt.Errorf("invalid output format %q: must be 'json' or 'text'", flags.Output)
			}

			debug.SetupLogger(flags.Debug)
			ctx := debug.WithDebug(cmd.Context(), flags.Debug)

			u := ui.New(flags.Color)
			ctx = ui.WithUI(ctx, u)

			ctx = outfmt.WithFormat(ctx, flags.Output)
			ctx = outfmt.WithYes(ctx, flags.Yes)

			cmd.SetContext(ctx)
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&flags.Token, "token", "", "API token (overrides keyring and env)")
	cmd.PersistentFlags().StringVar(&flags.Account, "account", os.Getenv(config.EnvToken), "Account name from keyring")
	cmd.PersistentFlags().StringVar(&flags.Output, "output", getEnvOrDefault(config.EnvOutput, "text"), "Output format: text|json")
	cmd.PersistentFlags().StringVar(&flags.Color, "color", getEnvOrDefault(config.EnvColor, "auto"), "Color output: auto|always|never")
	cmd.PersistentFlags().BoolVar(&flags.Debug, "debug", false, "Enable debug output")
	cmd.PersistentFlags().BoolVarP(&flags.Yes, "yes", "y", false, "Skip confirmation prompts")

	cmd.AddCommand(newVersionCmd())
	cmd.AddCommand(newAuthCmd())
	cmd.AddCommand(newUsersCmd())
	cmd.AddCommand(newDevicesCmd())
	cmd.AddCommand(newProfilesCmd())
	cmd.AddCommand(newServicesCmd())
	cmd.AddCommand(newNetworkCmd())
	cmd.AddCommand(newAccessCmd())
	cmd.AddCommand(newCompletionCmd())

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

func getClient(ctx context.Context) (*controld.API, error) {
	return api.NewClient(ctx, api.ClientConfig{
		Token:   flags.Token,
		Account: flags.Account,
	})
}
