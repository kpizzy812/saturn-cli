package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/saturn-platform/saturn-cli/internal/cli"
	"github.com/saturn-platform/saturn-cli/internal/config"
)

// NewUseCommand creates the use command
func NewUseCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "use <context_name>",
		Example: `context use myserver`,
		Args:    cli.ExactArgs(1, "<context_name>"),
		Short:   "Switch to a different context (set as default)",
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			if err := cfg.SetDefault(name); err != nil {
				return fmt.Errorf("Context '%s' not found", name)
			}

			if err := cfg.Save(); err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}

			fmt.Printf("Switched to context '%s'.\n", name)
			return nil
		},
	}
}
