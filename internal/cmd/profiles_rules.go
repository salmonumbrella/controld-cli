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
	cmd.AddCommand(newProfilesRulesFoldersCmd())
	cmd.AddCommand(newProfilesRulesListCmd())
	cmd.AddCommand(newProfilesRulesCreateCmd())
	cmd.AddCommand(newProfilesRulesDeleteCmd())
	return cmd
}

func newProfilesRulesFoldersCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "folders <profile-id>",
		Short: "List rule folders for a profile",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			folders, err := client.ListProfileRuleFolders(cmd.Context(), controld.ListProfileRuleFoldersParams{
				ProfileID: args[0],
			})
			if err != nil {
				return err
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, folders)
			}

			if len(folders) == 0 {
				fmt.Println("No rule folders found")
				return nil
			}

			tw := outfmt.NewTabWriter(os.Stdout)
			_, _ = fmt.Fprintln(tw, "FOLDER_ID\tNAME\tRULE_COUNT\tSTATUS")
			for _, f := range folders {
				status := "disabled"
				if f.Action.Status {
					status = "enabled"
				}
				_, _ = fmt.Fprintf(tw, "%d\t%s\t%d\t%s\n", f.PK, f.Group, f.Count, status)
			}
			return tw.Flush()
		},
	}
}

func newProfilesRulesListCmd() *cobra.Command {
	var folderID string

	cmd := &cobra.Command{
		Use:   "list <profile-id>",
		Short: "List custom rules for a profile",
		Long: `List custom rules for a profile.

Rules are organized into folders. Use --folder to specify a folder ID,
or use 'profiles rules folders <profile-id>' to list available folders first.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			profileID := args[0]

			// If no folder specified, first get the list of folders
			if folderID == "" {
				folders, err := client.ListProfileRuleFolders(cmd.Context(), controld.ListProfileRuleFoldersParams{
					ProfileID: profileID,
				})
				if err != nil {
					return err
				}

				if len(folders) == 0 {
					if outfmt.IsJSON(cmd.Context()) {
						return outfmt.WriteJSON(os.Stdout, []controld.CustomRule{})
					}
					fmt.Println("No rule folders found. Create a folder first or add rules via the ControlD dashboard.")
					return nil
				}

				// Use the first folder by default
				folderID = fmt.Sprintf("%d", folders[0].PK)
			}

			rules, err := client.ListProfileCustomRules(cmd.Context(), controld.ListProfileCustomRulesParams{
				ProfileID: profileID,
				FolderID:  folderID,
			})
			if err != nil {
				return err
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, rules)
			}

			if len(rules) == 0 {
				fmt.Printf("No custom rules found in folder %s\n", folderID)
				return nil
			}

			tw := outfmt.NewTabWriter(os.Stdout)
			_, _ = fmt.Fprintln(tw, "DOMAIN\tACTION\tSTATUS\tFOLDER")
			for _, r := range rules {
				action := actionToString(r.Action.Do)
				status := "disabled"
				if r.Action.Status {
					status = "enabled"
				}
				_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", r.PK, action, status, folderID)
			}
			return tw.Flush()
		},
	}

	cmd.Flags().StringVar(&folderID, "folder", "", "Folder ID (use 'folders' subcommand to list available folders)")
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
