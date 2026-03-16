package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/saturn-platform/saturn-cli/internal/cli"
	"github.com/saturn-platform/saturn-cli/internal/config"
)

// NewAddCommand creates the add command
func NewAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "add <context_name> <url> <token>",
		Example: `context add myserver https://saturn.example.com your-api-token`,
		Args:    cli.ExactArgs(3, "<context_name> <url> <token>"),
		Short:   "Add a new context",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			host := args[1]
			token := args[2]

			force, _ := cmd.Flags().GetBool("force")
			setDefault, _ := cmd.Flags().GetBool("default")

			cfg, err := config.Load()
			if err != nil {
				cfg = config.New()
			}

			// Check if instance already exists
			existing, _ := cfg.GetInstance(name)
			if existing != nil {
				if !force {
					fmt.Printf("%s already exists.\n", name)
					fmt.Println("\nNote: Use --force to force overwrite.")
					return nil
				}

				// Force overwrite: update token and FQDN
				existing.Token = token
				existing.FQDN = host

				if setDefault {
					if err := cfg.SetDefault(name); err != nil {
						return fmt.Errorf("failed to set default: %w", err)
					}
					fmt.Printf("%s already exists. Force overwriting. Setting it as default.\n", name)
				} else {
					fmt.Printf("%s already exists. Force overwriting.\n", name)
				}

				if err := cfg.Save(); err != nil {
					return fmt.Errorf("failed to write config: %w", err)
				}
				return nil
			}

			// Add new instance
			newInstance := config.Instance{
				Name:    name,
				FQDN:    host,
				Token:   token,
				Default: setDefault,
			}

			if err := cfg.AddInstance(newInstance); err != nil {
				return fmt.Errorf("failed to add instance: %w", err)
			}

			if setDefault {
				fmt.Printf("Context '%s' added and set as default.\n", name)
			} else {
				fmt.Printf("Context '%s' added successfully.\n", name)
			}

			if err := cfg.Save(); err != nil {
				return fmt.Errorf("failed to write config: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().BoolP("default", "d", false, "Set as default context")
	cmd.Flags().BoolP("force", "f", false, "Force overwrite if context already exists")

	return cmd
}
