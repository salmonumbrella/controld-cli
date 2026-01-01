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
			_, _ = fmt.Fprintln(tw, "POP\tCITY\tCOUNTRY\tAPI\tDNS\tPROXY")
			for _, n := range network {
				_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\t%d\t%d\t%d\n",
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
			_, _ = fmt.Fprintf(tw, "ip\t%s\n", ip.IP)
			_, _ = fmt.Fprintf(tw, "type\t%s\n", ip.Type)
			_, _ = fmt.Fprintf(tw, "country\t%s\n", ip.Country)
			_, _ = fmt.Fprintf(tw, "org\t%s\n", ip.Org)
			_, _ = fmt.Fprintf(tw, "pop\t%s\n", ip.Pop)
			return tw.Flush()
		},
	}
}
