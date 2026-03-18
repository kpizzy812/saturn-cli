package database

import (
	"fmt"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/saturn-platform/saturn-cli/internal/cli"
	"github.com/saturn-platform/saturn-cli/internal/output"
	"github.com/saturn-platform/saturn-cli/internal/service"
)

// databaseListItem is a simplified view for table display
type databaseListItem struct {
	UUID   string `json:"uuid"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

// NewListCommand lists all databases
func NewListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all databases",
		Long:  `List all databases in Saturn.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()

			client, err := cli.GetAPIClient(cmd)
			if err != nil {
				return fmt.Errorf("failed to get API client: %w", err)
			}

			dbService := service.NewDatabaseService(client)
			databases, err := dbService.List(ctx)
			if err != nil {
				return fmt.Errorf("failed to list databases: %w", err)
			}

			format, _ := cmd.Flags().GetString("format")

			if format != "" && format != output.FormatTable {
				formatter, err := output.NewFormatter(format, output.Options{})
				if err != nil {
					return err
				}
				return formatter.Format(databases)
			}

			// Table view: simplified rows with colored status
			rows := make([]databaseListItem, 0, len(databases))
			for _, db := range databases {
				rows = append(rows, databaseListItem{
					UUID:   db.UUID,
					Name:   db.Name,
					Type:   db.Type,
					Status: colorizeDatabaseStatus(db.Status),
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

// colorizeDatabaseStatus returns a pterm-colored string for a database status
func colorizeDatabaseStatus(status string) string {
	switch status {
	case "running":
		return pterm.FgGreen.Sprint(status)
	case "stopped", "exited":
		return pterm.FgGray.Sprint(status)
	case "restarting":
		return pterm.FgYellow.Sprint(status)
	case "error", "failed", "unhealthy":
		return pterm.FgRed.Sprint(status)
	case "starting":
		return pterm.FgCyan.Sprint(status)
	default:
		return status
	}
}
