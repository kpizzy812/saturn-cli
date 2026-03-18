package service

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/saturn-platform/saturn-cli/internal/cli"
	"github.com/saturn-platform/saturn-cli/internal/output"
	"github.com/saturn-platform/saturn-cli/internal/service"
)

// serviceListItem is a simplified view for table display
type serviceListItem struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

// NewListCommand lists all services
func NewListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all services",
		Long:  `List all services in Saturn.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			client, err := cli.GetAPIClient(cmd)
			if err != nil {
				return fmt.Errorf("failed to get API client: %w", err)
			}

			serviceSvc := service.NewService(client)
			services, err := serviceSvc.List(ctx)
			if err != nil {
				return fmt.Errorf("failed to list services: %w", err)
			}

			format, _ := cmd.Flags().GetString("format")

			if format != "" && format != output.FormatTable {
				formatter, err := output.NewFormatter(format, output.Options{})
				if err != nil {
					return err
				}
				return formatter.Format(services)
			}

			// Table view: simplified rows with colored status
			rows := make([]serviceListItem, 0, len(services))
			for _, svc := range services {
				rows = append(rows, serviceListItem{
					UUID:   svc.UUID,
					Name:   svc.Name,
					Status: colorizeServiceStatus(svc.Status),
				})
			}

			formatter, err := output.NewFormatter(output.FormatTable, output.Options{})
			if err != nil {
				return fmt.Errorf("failed to create formatter: %w", err)
			}
			return formatter.Format(rows)
		},
	}
}

// colorizeServiceStatus returns a pterm-colored string for a service status
func colorizeServiceStatus(status string) string {
	switch status {
	case "running":
		return pterm.FgGreen.Sprint(status)
	case "stopped", "exited":
		return pterm.FgGray.Sprint(status)
	case "restarting":
		return pterm.FgYellow.Sprint(status)
	case "error", "failed", "unhealthy":
		return pterm.FgRed.Sprint(status)
	case "starting", "deploying":
		return pterm.FgCyan.Sprint(status)
	default:
		return status
	}
}
