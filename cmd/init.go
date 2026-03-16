package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/saturn-platform/saturn-cli/internal/cli"
	"github.com/saturn-platform/saturn-cli/internal/models"
	"github.com/saturn-platform/saturn-cli/internal/service"
)

// NewInitCommand creates the init command
func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize Saturn project in the current directory",
		Long: `Interactive wizard to link the current directory to a Saturn resource.

Detects your git remote, matches it against Saturn resources, and generates
a .saturn.yml config file. After init, you can deploy with:

  saturn deploy smart

If no matching resource is found, you can select one interactively.`,
		Example: `  saturn init
  saturn init --yes`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return runInit(cmd)
		},
	}

	cmd.Flags().BoolP("yes", "y", false, "Skip confirmation prompts")

	return cmd
}

func runInit(cmd *cobra.Command) error {
	ctx := cmd.Context()
	reader := bufio.NewReader(os.Stdin)

	// Check if .saturn.yml already exists
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	existingCfg, _ := service.LoadConfig(dir)
	if existingCfg != nil {
		yes, _ := cmd.Flags().GetBool("yes")
		if !yes {
			fmt.Print(".saturn.yml already exists. Overwrite? [y/N] ")
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToLower(answer))
			if !strings.HasPrefix(answer, "y") {
				fmt.Println("Aborted")
				return nil
			}
		}
	}

	// Get API client (auto-login if needed)
	client, err := cli.GetAPIClient(cmd)
	if err != nil {
		return fmt.Errorf("failed to connect to Saturn: %w", err)
	}

	resourceSvc := service.NewResourceService(client)

	// Try to auto-detect by git remote
	fmt.Println("Detecting git remote...")
	localURL, gitErr := service.GetGitRemoteURL(ctx)

	var matched []models.Resource

	if gitErr == nil {
		localNorm := service.NormalizeGitURL(localURL)
		fmt.Printf("Found git remote: %s\n", localURL)

		resources, err := resourceSvc.List(ctx)
		if err != nil {
			return fmt.Errorf("failed to list resources: %w", err)
		}

		for _, r := range resources {
			if r.GitRepository == nil {
				continue
			}
			if service.NormalizeGitURL(*r.GitRepository) == localNorm {
				matched = append(matched, r)
			}
		}
	} else {
		fmt.Println("Not a git repository or no remote configured")
	}

	var cfg *models.SmartConfig

	if len(matched) > 0 {
		// Auto-detected resources
		fmt.Printf("\nFound %d Saturn resource(s) matching this repo:\n", len(matched))
		for i, r := range matched {
			status := r.Status
			if status == "" {
				status = "unknown"
			}
			baseDir := "/"
			if r.BaseDirectory != nil && *r.BaseDirectory != "" {
				baseDir = *r.BaseDirectory
			}
			fmt.Printf("  %d. %s (type: %s, status: %s, base: %s)\n", i+1, r.Name, r.Type, status, baseDir)
		}

		cfg = service.GenerateConfig(matched)
	} else {
		// No auto-detect match — let user pick from all resources
		fmt.Println("\nNo matching resources found by git remote.")
		fmt.Println("Fetching all resources...")

		resources, err := resourceSvc.List(ctx)
		if err != nil {
			return fmt.Errorf("failed to list resources: %w", err)
		}

		if len(resources) == 0 {
			return fmt.Errorf("no resources found in Saturn. Create an application first via the web UI")
		}

		fmt.Printf("\nAvailable resources (%d):\n", len(resources))
		for i, r := range resources {
			status := r.Status
			if status == "" {
				status = "unknown"
			}
			fmt.Printf("  %d. %s (type: %s, status: %s)\n", i+1, r.Name, r.Type, status)
		}

		fmt.Print("\nSelect resource(s) to link (comma-separated numbers, e.g. 1,3): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "" {
			fmt.Println("Aborted")
			return nil
		}

		var selected []models.Resource
		for _, part := range strings.Split(input, ",") {
			part = strings.TrimSpace(part)
			idx, err := strconv.Atoi(part)
			if err != nil || idx < 1 || idx > len(resources) {
				return fmt.Errorf("invalid selection: %s", part)
			}
			selected = append(selected, resources[idx-1])
		}

		cfg = service.GenerateConfig(selected)
	}

	if cfg == nil || len(cfg.Components) == 0 {
		return fmt.Errorf("no components to configure")
	}

	// Show what will be written
	fmt.Printf("\nGenerated .saturn.yml (%d component(s)):\n", len(cfg.Components))
	for name, comp := range cfg.Components {
		fmt.Printf("  %s: path=%q resource=%q\n", name, comp.Path, comp.Resource)
	}

	// Confirm
	yes, _ := cmd.Flags().GetBool("yes")
	if !yes {
		fmt.Print("\nWrite .saturn.yml? [Y/n] ")
		answer, _ := reader.ReadString('\n')
		answer = strings.TrimSpace(strings.ToLower(answer))
		if answer != "" && !strings.HasPrefix(answer, "y") {
			fmt.Println("Aborted")
			return nil
		}
	}

	// Write config
	if err := service.WriteConfig(dir, cfg); err != nil {
		return fmt.Errorf("failed to write .saturn.yml: %w", err)
	}

	fmt.Println("\n.saturn.yml created successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("  saturn deploy smart            # Deploy changed components")
	fmt.Println("  saturn deploy smart --dry-run   # Preview deploy plan")
	fmt.Println("  saturn deploy name <app-name>   # Deploy specific app")

	return nil
}
