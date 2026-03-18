package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/saturn-platform/saturn-cli/internal/cli"
	"github.com/saturn-platform/saturn-cli/internal/config"
)

// NewDeleteCommand creates the delete command
func NewDeleteCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <context_name>",
		Example: `context delete myserver`,
		Args:    cli.ExactArgs(1, "<context_name>"),
		Short:   "Delete a context",
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Check if it was the default before removing
			instance, err := cfg.GetInstance(name)
			if err != nil {
				return fmt.Errorf("context '%s' not found", name)
			}
			wasDefault := instance.Default

			if err := cfg.RemoveInstance(name); err != nil {
				return fmt.Errorf("failed to remove context: %w", err)
			}

			if err := cfg.Save(); err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}

			switch {
			case wasDefault && len(cfg.Instances) > 0:
				fmt.Printf("Context '%s' deleted. '%s' is now the default context.\n", name, cfg.Instances[0].Name)
			case len(cfg.Instances) == 0:
				fmt.Printf("Context '%s' deleted. No contexts remaining.\n", name)
			default:
				fmt.Printf("Context '%s' deleted.\n", name)
			}

			return nil
		},
	}
}
