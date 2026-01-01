package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/salmonumbrella/controld-cli/internal/controld"
	"github.com/salmonumbrella/controld-cli/internal/outfmt"
	"github.com/salmonumbrella/controld-cli/internal/ui"
)

func newDevicesCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "devices",
		Short: "Manage DNS resolver devices",
	}
	cmd.AddCommand(newDevicesListCmd())
	cmd.AddCommand(newDevicesGetCmd())
	cmd.AddCommand(newDevicesCreateCmd())
	cmd.AddCommand(newDevicesModifyCmd())
	cmd.AddCommand(newDevicesDeleteCmd())
	cmd.AddCommand(newDevicesTypesCmd())
	return cmd
}

func newDevicesListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all devices",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			devices, err := client.ListDevices(cmd.Context())
			if err != nil {
				return err
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, devices)
			}

			if len(devices) == 0 {
				fmt.Println("No devices found")
				return nil
			}

			tw := outfmt.NewTabWriter(os.Stdout)
			_, _ = fmt.Fprintln(tw, "DEVICE_ID\tNAME\tSTATUS\tPROFILE")
			for _, d := range devices {
				status := "pending"
				switch d.Status {
				case controld.Active:
					status = "active"
				case controld.SoftDisabled:
					status = "soft-disabled"
				case controld.HardDisabled:
					status = "hard-disabled"
				}
				_, _ = fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", d.DeviceID, d.Name, status, d.Profile.Name)
			}
			return tw.Flush()
		},
	}
}

func newDevicesGetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "get <device-id>",
		Short: "Get device details",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			devices, err := client.ListDevices(cmd.Context())
			if err != nil {
				return err
			}

			deviceID := args[0]
			var device *controld.Device
			for i := range devices {
				if devices[i].DeviceID == deviceID {
					device = &devices[i]
					break
				}
			}

			if device == nil {
				return fmt.Errorf("device not found: %s", deviceID)
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, device)
			}

			tw := outfmt.NewTabWriter(os.Stdout)
			_, _ = fmt.Fprintf(tw, "device_id\t%s\n", device.DeviceID)
			_, _ = fmt.Fprintf(tw, "name\t%s\n", device.Name)
			_, _ = fmt.Fprintf(tw, "status\t%d\n", device.Status)
			_, _ = fmt.Fprintf(tw, "profile\t%s\n", device.Profile.Name)
			_, _ = fmt.Fprintf(tw, "doh\t%s\n", device.Resolvers.DoH)
			_, _ = fmt.Fprintf(tw, "dot\t%s\n", device.Resolvers.DoT)
			if device.Icon != nil {
				_, _ = fmt.Fprintf(tw, "icon\t%s\n", *device.Icon)
			}
			return tw.Flush()
		},
	}
}

func newDevicesCreateCmd() *cobra.Command {
	var name string
	var profileID string
	var icon string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new device",
		RunE: func(cmd *cobra.Command, args []string) error {
			u := ui.FromContext(cmd.Context())
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			device, err := client.CreateDevice(cmd.Context(), controld.CreateDeviceParams{
				Name:      name,
				ProfileID: profileID,
				Icon:      controld.IconName(icon),
			})
			if err != nil {
				return err
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, device)
			}

			u.Success(fmt.Sprintf("Created device: %s (%s)", device.Name, device.DeviceID))
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Device name (required)")
	cmd.Flags().StringVar(&profileID, "profile-id", "", "Profile ID (required)")
	cmd.Flags().StringVar(&icon, "icon", "router", "Device icon")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("profile-id")
	return cmd
}

func newDevicesModifyCmd() *cobra.Command {
	var name string
	var profileID string
	var status int

	cmd := &cobra.Command{
		Use:   "modify <device-id>",
		Short: "Modify an existing device",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			u := ui.FromContext(cmd.Context())
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			params := controld.UpdateDeviceParams{
				DeviceID: args[0],
			}

			if cmd.Flags().Changed("name") {
				params.Name = &name
			}
			if cmd.Flags().Changed("profile-id") {
				params.ProfileID = &profileID
			}
			if cmd.Flags().Changed("status") {
				s := controld.DeviceStatus(status)
				params.Status = &s
			}

			device, err := client.UpdateDevice(cmd.Context(), params)
			if err != nil {
				return err
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, device)
			}

			u.Success(fmt.Sprintf("Modified device: %s (%s)", device.Name, device.DeviceID))
			return nil
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "New device name")
	cmd.Flags().StringVar(&profileID, "profile-id", "", "New profile ID")
	cmd.Flags().IntVar(&status, "status", 0, "Device status (0=pending, 1=active, 2=soft-disabled, 3=hard-disabled)")
	return cmd
}

func newDevicesDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete <device-id>",
		Short: "Delete a device",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			u := ui.FromContext(cmd.Context())
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			if !outfmt.GetYes(cmd.Context()) {
				_, _ = fmt.Fprintf(os.Stderr, "Delete device %s? [y/N]: ", args[0])
				var confirm string
				_, _ = fmt.Scanln(&confirm)
				if confirm != "y" && confirm != "Y" {
					_, _ = fmt.Fprintln(os.Stderr, "Cancelled")
					return nil
				}
			}

			_, err = client.DeleteDevice(cmd.Context(), controld.DeleteDeviceParams{
				DeviceID: args[0],
			})
			if err != nil {
				return err
			}

			u.Success(fmt.Sprintf("Deleted device: %s", args[0]))
			return nil
		},
	}
}

func newDevicesTypesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "types",
		Short: "List available device types",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := getClient(cmd.Context())
			if err != nil {
				return err
			}

			types, err := client.ListDeviceType(cmd.Context())
			if err != nil {
				return err
			}

			if outfmt.IsJSON(cmd.Context()) {
				return outfmt.WriteJSON(os.Stdout, types)
			}

			tw := outfmt.NewTabWriter(os.Stdout)
			_, _ = fmt.Fprintln(tw, "CATEGORY\tNAME")
			_, _ = fmt.Fprintf(tw, "os\t%s\n", types.OS.Name)
			_, _ = fmt.Fprintf(tw, "browser\t%s\n", types.Browser.Name)
			_, _ = fmt.Fprintf(tw, "tv\t%s\n", types.TV.Name)
			_, _ = fmt.Fprintf(tw, "router\t%s\n", types.Router.Name)
			return tw.Flush()
		},
	}
}
