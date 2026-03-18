package application

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/saturn-platform/saturn-cli/internal/cli"
	"github.com/saturn-platform/saturn-cli/internal/models"
	"github.com/saturn-platform/saturn-cli/internal/output"
	"github.com/saturn-platform/saturn-cli/internal/service"
)

func NewListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all applications",
		Long:  `List all applications in Saturn.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			client, err := cli.GetAPIClient(cmd)
			if err != nil {
				return fmt.Errorf("failed to get API client: %w", err)
			}

			appSvc := service.NewApplicationService(client)
			apps, err := appSvc.List(ctx)
			if err != nil {
				return fmt.Errorf("failed to list applications: %w", err)
			}

			format, _ := cmd.Flags().GetString("format")
			showSensitive, _ := cmd.Flags().GetBool("show-sensitive")

			// For JSON/pretty formats, return the full application structure
			if format != output.FormatTable {
				formatter, err := output.NewFormatter(format, output.Options{
					ShowSensitive: showSensitive,
				})
				if err != nil {
					return err
				}
				return formatter.Format(apps)
			}

			// For table format, convert to simplified rows with colored status
			var rows []models.ApplicationListItem
			for _, app := range apps {
				rows = append(rows, models.ApplicationListItem{
					UUID:        app.UUID,
					Name:        app.Name,
					Description: app.Description,
					Status:      colorizeStatus(app.Status),
					GitBranch:   app.GitBranch,
					FQDN:        app.FQDN,
				})
			}

			formatter, err := output.NewFormatter(format, output.Options{
				ShowSensitive: showSensitive,
			})
			if err != nil {
				return err
			}

			return formatter.Format(rows)
		},
	}
}

// colorizeStatus returns a pterm-colored string for the given deployment status
func colorizeStatus(status string) string {
	switch status {
	case "running":
		return pterm.FgGreen.Sprint(status)
	case "stopped", "exited":
		return pterm.FgGray.Sprint(status)
	case "restarting":
		return pterm.FgYellow.Sprint(status)
	case "error", "failed", "unhealthy":
		return pterm.FgRed.Sprint(status)
	case "starting", "building", "deploying":
		return pterm.FgCyan.Sprint(status)
	default:
		return status
	}
}
