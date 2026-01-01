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
