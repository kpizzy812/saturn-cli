package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/saturn-platform/saturn-cli/internal/cli"
	"github.com/saturn-platform/saturn-cli/internal/config"
	"github.com/saturn-platform/saturn-cli/internal/output"
	"github.com/saturn-platform/saturn-cli/internal/service"
)

// NewStatusCommand creates the status command
func NewStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show current context and resource summary",
		Long:  `Display the active Saturn context, authenticated team, and a summary of resources.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			styler := output.DefaultStyler()

			// Load config
			cfg, err := config.Load()
			if err != nil || len(cfg.Instances) == 0 {
				styler.Warning("Not logged in. Run: saturn login")
				return nil
			}

			instance, err := cfg.GetDefault()
			if err != nil {
				styler.Warning("No default context set. Run: saturn context use <name>")
				return nil
			}

			styler.Bold("Saturn Status")
			fmt.Println()
			fmt.Printf("  Context:  %s (%s)\n", instance.Name, instance.FQDN)

			// Try to fetch live data (best-effort — errors are silently ignored)
			apiClient, err := cli.GetAPIClient(cmd)
			if err == nil {
				ctx := cmd.Context()

				// Team info — use current team endpoint for the active token
				teamSvc := service.NewTeamService(apiClient)
				if team, teamErr := teamSvc.Current(ctx); teamErr == nil {
					fmt.Printf("  Team:     %s\n", team.Name)
				} else {
					// Fallback: list teams and use the first one
					if teams, listErr := teamSvc.List(ctx); listErr == nil && len(teams) > 0 {
						fmt.Printf("  Team:     %s\n", teams[0].Name)
					}
				}

				// Server count
				serverSvc := service.NewServerService(apiClient)
				if servers, serversErr := serverSvc.List(ctx); serversErr == nil {
					fmt.Printf("  Servers:  %d connected\n", len(servers))
				}

				// Application count
				appSvc := service.NewApplicationService(apiClient)
				if apps, appsErr := appSvc.List(ctx); appsErr == nil {
					fmt.Printf("  Apps:     %d total\n", len(apps))
				}
			}

			return nil
		},
	}
}
