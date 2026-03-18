package server

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/saturn-platform/saturn-cli/internal/cli"
	"github.com/saturn-platform/saturn-cli/internal/output"
	"github.com/saturn-platform/saturn-cli/internal/service"
)

// serverListItem is a simplified view for table display
type serverListItem struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	IP        string `json:"ip" sensitive:"true"`
	Reachable string `json:"reachable"`
	Usable    string `json:"usable"`
}

// NewListCommand creates the list command
func NewListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all servers",
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			client, err := cli.GetAPIClient(cmd)
			if err != nil {
				return fmt.Errorf("failed to get API client: %w", err)
			}

			serverSvc := service.NewServerService(client)
			servers, err := serverSvc.List(ctx)
			if err != nil {
				return fmt.Errorf("failed to list servers: %w", err)
			}

			format, _ := cmd.Flags().GetString("format")
			showSensitive, _ := cmd.Flags().GetBool("show-sensitive")

			if format != "" && format != output.FormatTable {
				formatter, err := output.NewFormatter(format, output.Options{
					ShowSensitive: showSensitive,
				})
				if err != nil {
					return err
				}
				return formatter.Format(servers)
			}

			// Table view: simplified rows with colored reachability
			rows := make([]serverListItem, 0, len(servers))
			for _, s := range servers {
				ip := ""
				if showSensitive {
					ip = s.IP
				}
				rows = append(rows, serverListItem{
					UUID:      s.UUID,
					Name:      s.Name,
					IP:        ip,
					Reachable: colorizeBoolean(s.Settings.IsReachable),
					Usable:    colorizeBoolean(s.Settings.IsUsable),
				})
			}

			formatter, err := output.NewFormatter(output.FormatTable, output.Options{})
			if err != nil {
				return fmt.Errorf("failed to create formatter: %w", err)
			}
			if err := formatter.Format(rows); err != nil {
				return err
			}

			if !showSensitive {
				fmt.Println("\nNote: Use -s to show sensitive information (IP address).")
			}

			return nil
		},
	}
}

// colorizeBoolean returns a colored yes/no string
func colorizeBoolean(v bool) string {
	if v {
		return pterm.FgGreen.Sprint("yes")
	}
	return pterm.FgRed.Sprint("no")
}
