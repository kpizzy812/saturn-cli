package context

import (
	"github.com/spf13/cobra"

	"github.com/saturn-platform/saturn-cli/internal/config"
	"github.com/saturn-platform/saturn-cli/internal/output"
)

// NewListCommand creates the list command
func NewListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all configured contexts",
		RunE: func(cmd *cobra.Command, _ []string) error {
			cfg, err := config.Load()
			if err != nil {
				cfg = config.New()
			}

			format, _ := cmd.Flags().GetString("format")
			showSensitive, _ := cmd.Flags().GetBool("show-sensitive")

			formatter, fmtErr := output.NewFormatter(format, output.Options{
				ShowSensitive: showSensitive,
			})
			if fmtErr != nil {
				return fmtErr
			}

			return formatter.Format(cfg.Instances)
		},
	}
}
