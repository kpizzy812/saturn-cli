package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/saturn-platform/saturn-cli/internal/cli"
	"github.com/saturn-platform/saturn-cli/internal/output"
	"github.com/saturn-platform/saturn-cli/internal/service"
)

// NewDevCommand creates the dev command for local development
func NewDevCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dev",
		Short: "Run local dev server with Saturn environment variables",
		Long: `Pull environment variables from Saturn and start your local dev server.

If .env.local already exists it will be reused unless --refresh is set.
The dev command is detected automatically or can be overridden with --command.

Examples:
  saturn dev --app <uuid>
  saturn dev --app <uuid> --command "go run ./..."
  saturn dev --app <uuid> --refresh`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			ctx := cmd.Context()
			styler := output.DefaultStyler()

			appUUID, _ := cmd.Flags().GetString("app")
			envFile, _ := cmd.Flags().GetString("env-file")
			devCommand, _ := cmd.Flags().GetString("command")
			refresh, _ := cmd.Flags().GetBool("refresh")

			// Pull env vars if needed
			_, statErr := os.Stat(envFile)
			needPull := os.IsNotExist(statErr) || refresh

			if needPull {
				if appUUID == "" {
					return fmt.Errorf("--app <uuid> is required to pull environment variables")
				}

				client, err := cli.GetAPIClient(cmd)
				if err != nil {
					return fmt.Errorf("failed to get API client: %w", err)
				}

				sp := output.NewSpinner("Pulling environment variables from Saturn...")

				appSvc := service.NewApplicationService(client)
				envVars, err := appSvc.ListEnvs(ctx, appUUID)
				if err != nil {
					sp.Fail("Failed to pull env vars")
					return fmt.Errorf("failed to get environment variables: %w", err)
				}

				var lines []string
				for _, ev := range envVars {
					if ev.Key != "" {
						val := ev.Value
						if strings.ContainsAny(val, " \t\n\"'#") {
							val = fmt.Sprintf("%q", val)
						}
						lines = append(lines, fmt.Sprintf("%s=%s", ev.Key, val))
					}
				}

				content := strings.Join(lines, "\n") + "\n"
				if err := os.WriteFile(envFile, []byte(content), 0600); err != nil {
					sp.Fail("Failed to write env file")
					return fmt.Errorf("failed to write %s: %w", envFile, err)
				}

				sp.Success(fmt.Sprintf("Pulled %d variables to %s", len(envVars), envFile))
				styler.Warning("Remember to add " + envFile + " to .gitignore!")
			} else {
				styler.Info(fmt.Sprintf("Using existing %s (use --refresh to update)", envFile))
			}

			// Detect dev command if not provided
			if devCommand == "" {
				devCommand = detectDevCommand()
			}

			if devCommand == "" {
				return fmt.Errorf("could not detect dev command — use --command to specify one")
			}

			styler.Info(fmt.Sprintf("Starting: %s", devCommand))

			// Build command with env vars loaded from envFile
			shell, shellFlag := detectShell()
			cmdArgs := []string{shellFlag, fmt.Sprintf("set -a; source %s; set +a; %s", envFile, devCommand)}
			c := exec.CommandContext(ctx, shell, cmdArgs...)
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr

			// Forward signals
			sigCh := make(chan os.Signal, 1)
			signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				sig := <-sigCh
				if c.Process != nil {
					_ = c.Process.Signal(sig)
				}
			}()
			defer signal.Stop(sigCh)

			if err := c.Run(); err != nil {
				// Exit code 130 (Ctrl+C) is expected, don't surface as error
				var exitErr *exec.ExitError
				if ok := isExitError(err, &exitErr); ok && exitErr.ExitCode() == 130 {
					return nil
				}
				return err
			}

			return nil
		},
	}

	cmd.Flags().String("app", "", "Application UUID to pull env vars from")
	cmd.Flags().String("env-file", ".env.local", "Local env file to load")
	cmd.Flags().String("command", "", "Dev command to run (auto-detected if not set)")
	cmd.Flags().Bool("refresh", false, "Force re-pull env vars even if env file exists")

	return cmd
}

// detectDevCommand auto-detects the dev command from common project files
func detectDevCommand() string {
	// Check package.json for scripts.dev
	if _, err := os.Stat("package.json"); err == nil {
		data, err := os.ReadFile("package.json")
		if err == nil && strings.Contains(string(data), `"dev"`) {
			return "npm run dev"
		}
		return "npm start"
	}
	// Go project
	if _, err := os.Stat("go.mod"); err == nil {
		return "go run ./..."
	}
	// Python
	if _, err := os.Stat("manage.py"); err == nil {
		return "python manage.py runserver"
	}
	if _, err := os.Stat("pyproject.toml"); err == nil {
		return "python -m uvicorn main:app --reload"
	}
	return ""
}

// detectShell returns the shell binary and the flag to execute a command string
func detectShell() (string, string) {
	if shell := os.Getenv("SHELL"); shell != "" {
		return shell, "-c"
	}
	// Fallback to sh
	return "sh", "-c"
}

// isExitError checks if err is an *exec.ExitError and sets the pointer
func isExitError(err error, target **exec.ExitError) bool {
	var exitErr *exec.ExitError
	if ok := func() bool {
		exitErr, _ = err.(*exec.ExitError)
		return exitErr != nil
	}(); ok {
		*target = exitErr
		return true
	}
	return false
}
