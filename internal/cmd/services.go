package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/controld-cli/internal/controld"
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
				params.Category = category
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
			_, _ = fmt.Fprintln(tw, "SERVICE_ID\tNAME\tCATEGORY")
			for _, s := range services {
				_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\n", s.PK, s.Name, s.Category)
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
			_, _ = fmt.Fprintln(tw, "CATEGORY_ID\tNAME\tCOUNT")
			for _, c := range categories {
				_, _ = fmt.Fprintf(tw, "%s\t%s\t%d\n", c.PK, c.Name, c.Count)
			}
			return tw.Flush()
		},
	}
}
