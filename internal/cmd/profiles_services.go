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
	cmd.AddCommand(newProfilesServicesDisableCmd())
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
			_, _ = fmt.Fprintln(tw, "SERVICE_ID\tNAME\tCATEGORY\tACTION\tSTATUS")
			for _, s := range services {
				action := actionToString(s.Action.Do)
				status := "disabled"
				if s.Action.Status {
					status = "enabled"
				}
				_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\t%s\t%s\n", s.PK, s.Name, s.Category, action, status)
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
		Long: `Set action for a service.

Actions:
  block   - Block access to this service
  bypass  - Allow access, bypassing any filters
  spoof   - Use proxy/redirect for geo-unblocking`,
		Args: cobra.ExactArgs(2),
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

func newProfilesServicesDisableCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "disable <profile-id> <service-id>",
		Short: "Disable a service rule (remove custom action)",
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
				Do:        controld.Block,
				Status:    controld.IntBool(false),
			})
			if err != nil {
				return err
			}

			u.Success(fmt.Sprintf("Disabled service rule: %s", args[1]))
			return nil
		},
	}
}
