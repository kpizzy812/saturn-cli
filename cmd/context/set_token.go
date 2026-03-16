package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/saturn-platform/saturn-cli/internal/cli"
	"github.com/saturn-platform/saturn-cli/internal/config"
)

// NewSetTokenCommand creates the set-token command
func NewSetTokenCommand() *cobra.Command {
	return &cobra.Command{
		Use:     "set-token <context_name> <token>",
		Example: `context set-token myserver your-new-api-token`,
		Args:    cli.ExactArgs(2, "<context_name> <token>"),
		Short:   "Update the API token for a context",
		RunE: func(_ *cobra.Command, args []string) error {
			name := args[0]
			token := args[1]

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			if err := cfg.UpdateInstanceToken(name, token); err != nil {
				return fmt.Errorf("failed to update token for context '%s': %w", name, err)
			}

			if err := cfg.Save(); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			fmt.Printf("Token updated for context '%s'.\n", name)
			return nil
		},
	}
}
