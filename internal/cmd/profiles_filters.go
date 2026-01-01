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
			_, _ = fmt.Fprintln(tw, "FILTER_ID\tNAME\tSTATUS")
			for _, f := range filters {
				status := "disabled"
				if f.Status {
					status = "enabled"
				}
				_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\n", f.PK, f.Name, status)
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
