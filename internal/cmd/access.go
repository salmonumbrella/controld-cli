package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/controld-cli/internal/controld"
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
			_, _ = fmt.Fprintln(tw, "IP\tCOUNTRY\tCITY\tISP")
			for _, ip := range ips {
				_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", ip.IP, ip.Country, ip.City, ip.ISP)
			}
			return tw.Flush()
		},
	}
}

func parseIPs(args []string) ([]net.IP, error) {
	var ips []net.IP
	for _, ipStr := range args {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			return nil, fmt.Errorf("invalid IP: %s", ipStr)
		}
		ips = append(ips, ip)
	}
	return ips, nil
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

			ips, err := parseIPs(args[1:])
			if err != nil {
				return err
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

			ips, err := parseIPs(args[1:])
			if err != nil {
				return err
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
