package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/saturn-platform/saturn-cli/internal/cli"
	"github.com/saturn-platform/saturn-cli/internal/config"
)

// NewSetDefaultCommand creates the set-default command
func NewSetDefaultCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "set-default <context_name>",
		Example: `context set-default myserver`,
		Args:    cli.ExactArgs(1, "<context_name>"),
		Short:   "Set a context as the default",
		RunE: func(cmd *cobra.Command, args []string) error {
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

			// Show the list after updating
			return NewListCommand().RunE(cmd, args)
		},
	}
}
