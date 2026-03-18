package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/saturn-platform/saturn-cli/internal/cli"
	"github.com/saturn-platform/saturn-cli/internal/output"
	"github.com/saturn-platform/saturn-cli/internal/service"
)

// NewPullCommand creates the pull command that downloads env vars
func NewPullCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "pull <uuid>",
		Short: "Download environment variables to a local .env file",
		Long:  `Pull environment variables from a Saturn application and save them to a local file (default: .env.local).`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			styler := output.DefaultStyler()
			uuid := args[0]

			client, err := cli.GetAPIClient(cmd)
			if err != nil {
				return fmt.Errorf("failed to get API client: %w", err)
			}

			outputFile, _ := cmd.Flags().GetString("output")
			force, _ := cmd.Flags().GetBool("force")

			// Check if output file already exists
			if !force {
				if _, err := os.Stat(outputFile); err == nil {
					return fmt.Errorf("file %s already exists, use --force to overwrite", outputFile)
				}
			}

			sp := output.NewSpinner("Pulling environment variables...")

			appSvc := service.NewApplicationService(client)
			envVars, err := appSvc.ListEnvs(ctx, uuid)
			if err != nil {
				sp.Fail("Failed to pull env vars")
				return fmt.Errorf("failed to get environment variables: %w", err)
			}

			// Build .env content
			var lines []string
			for _, env := range envVars {
				if env.Key != "" {
					// Quote values that contain spaces/special chars
					val := env.Value
					if strings.ContainsAny(val, " \t\n\"'#") {
						val = fmt.Sprintf("%q", val)
					}
					lines = append(lines, fmt.Sprintf("%s=%s", env.Key, val))
				}
			}

			content := strings.Join(lines, "\n") + "\n"
			if err := os.WriteFile(outputFile, []byte(content), 0600); err != nil {
				sp.Fail("Failed to write file")
				return fmt.Errorf("failed to write %s: %w", outputFile, err)
			}

			sp.Success(fmt.Sprintf("Pulled %d variables to %s", len(envVars), outputFile))
			styler.Warning("Remember to add " + outputFile + " to .gitignore!")

			return nil
		},
	}

	cmd.Flags().StringP("output", "o", ".env.local", "Output file path")
	cmd.Flags().BoolP("force", "f", false, "Overwrite existing file")
	return cmd
}
