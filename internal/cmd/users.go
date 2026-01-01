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
			_, _ = fmt.Fprintf(tw, "email\t%s\n", user.Email)
			_, _ = fmt.Fprintf(tw, "status\t%v\n", user.Status)
			_, _ = fmt.Fprintf(tw, "resolver_ip\t%s\n", user.ResolverIP)
			_, _ = fmt.Fprintf(tw, "stats_endpoint\t%s\n", user.StatsEndpoint)
			_, _ = fmt.Fprintf(tw, "2fa\t%v\n", user.Twofa)
			return tw.Flush()
		},
	}
}
