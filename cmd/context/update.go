package context

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/saturn-platform/saturn-cli/internal/cli"
	"github.com/saturn-platform/saturn-cli/internal/config"
)

// NewUpdateCommand creates the update command
func NewUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <context_name>",
		Example: `context update myserver --name newname --url https://new.saturn.com --token newtoken`,
		Args:    cli.ExactArgs(1, "<context_name>"),
		Short:   "Update a context's properties (name, URL, token)",
		RunE: func(cmd *cobra.Command, args []string) error {
			oldName := args[0]

			newName, _ := cmd.Flags().GetString("name")
			newURL, _ := cmd.Flags().GetString("url")
			newToken, _ := cmd.Flags().GetString("token")

			if newName == "" && newURL == "" && newToken == "" {
				return fmt.Errorf("at least one of --name, --url, or --token must be provided")
			}

			cfg, err := config.Load()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			instance, err := cfg.GetInstance(oldName)
			if err != nil {
				return fmt.Errorf("context '%s' not found", oldName)
			}

			// If renaming, check if new name already exists
			if newName != "" && newName != oldName {
				if _, dupErr := cfg.GetInstance(newName); dupErr == nil {
					return fmt.Errorf("context with name '%s' already exists", newName)
				}
				instance.Name = newName
			}

			if newURL != "" {
				instance.FQDN = newURL
			}

			if newToken != "" {
				instance.Token = newToken
			}

			if err := cfg.Save(); err != nil {
				return fmt.Errorf("failed to save config: %w", err)
			}

			finalName := oldName
			if newName != "" {
				finalName = newName
			}
			fmt.Printf("Context '%s' updated successfully.\n", finalName)
			return nil
		},
	}

	cmd.Flags().StringP("name", "n", "", "New name for the context")
	cmd.Flags().StringP("url", "u", "", "New URL for the context")
	cmd.Flags().StringP("token", "t", "", "New token for the context")

	return cmd
}
