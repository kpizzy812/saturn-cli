package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/saturn-platform/saturn-cli/internal/cli"
	"github.com/saturn-platform/saturn-cli/internal/config"
	"github.com/saturn-platform/saturn-cli/internal/output"
)

// NewGetCommand creates the get command
func NewGetCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "get <context_name>",
		Example: `context get myserver`,
		Args:    cli.ExactArgs(1, "<context_name>"),
		Short:   "Get details of a specific context",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			instance, err := cfg.GetInstance(name)
			if err != nil {
				return fmt.Errorf("Context '%s' not found", name)
			}

			format, _ := cmd.Flags().GetString("format")
			showSensitive, _ := cmd.Flags().GetBool("show-sensitive")

			formatter, fmtErr := output.NewFormatter(format, output.Options{
				ShowSensitive: showSensitive,
			})
			if fmtErr != nil {
				return fmtErr
			}

			return formatter.Format([]config.Instance{*instance})
		},
	}
}
