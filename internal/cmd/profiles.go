package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/controld-cli/internal/controld"
	"github.com/salmonumbrella/controld-cli/internal/outfmt"
	"github.com/salmonumbrella/controld-cli/internal/ui"
)

func newProfilesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profiles",
		Short: "Manage DNS filtering profiles",
	}
	cmd.AddCommand(newProfilesListCmd())
	cmd.AddCommand(newProfilesGetCmd())
	cmd.AddCommand(newProfilesCreateCmd())
	cmd.AddCommand(newProfilesModifyCmd())
	cmd.AddCommand(newProfilesDeleteCmd())
	cmd.AddCommand(newProfilesRulesCmd())
	cmd.AddCommand(newProfilesFiltersCmd())
	cmd.AddCommand(newProfilesServicesCmd())
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
			_, _ = fmt.Fprintln(tw, "PROFILE_ID\tNAME\tUPDATED")
			for _, p := range profiles {
				_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\n", p.PK, p.Name, p.Updated.Format("2006-01-02"))
			}
			return tw.Flush()
		},
	}
}

func newProfilesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <profile-id>",
		Short: "Get profile details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			profiles, err := client.ListProfiles(cmd.Context())
			if err != nil {
				return err
			}

			profileID := args[0]
			var profile *controld.Profile
			for i := range profiles {
				if profiles[i].PK == profileID {
					profile = &profiles[i]
					break
				}
			}

			if profile == nil {
				return fmt.Errorf("profile not found: %s", profileID)
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, profile)
			}

			tw := outfmt.NewTabWriter(os.Stdout)
			_, _ = fmt.Fprintf(tw, "profile_id\t%s\n", profile.PK)
			_, _ = fmt.Fprintf(tw, "name\t%s\n", profile.Name)
			_, _ = fmt.Fprintf(tw, "updated\t%s\n", profile.Updated.Format("2006-01-02 15:04:05"))
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

func newProfilesModifyCmd() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "modify <profile-id>",
		Short: "Modify an existing profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			u := ui.FromContext(cmd.Context())
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			params := controld.UpdateProfileParams{
				ProfileID: args[0],
			}

			if cmd.Flags().Changed("name") {
				params.Name = &name
			}

			profiles, err := client.UpdateProfile(cmd.Context(), params)
			if err != nil {
				return err
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, profiles)
			}

			if len(profiles) > 0 {
				u.Success(fmt.Sprintf("Modified profile: %s (%s)", profiles[0].Name, profiles[0].PK))
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "New profile name")
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
				_, _ = fmt.Fprintf(os.Stderr, "Delete profile %s? [y/N]: ", args[0])
				var confirm string
				_, _ = fmt.Scanln(&confirm)
				if confirm != "y" && confirm != "Y" {
					_, _ = fmt.Fprintln(os.Stderr, "Cancelled")
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
