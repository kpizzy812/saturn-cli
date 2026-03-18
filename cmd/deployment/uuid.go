package deployment

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/saturn-platform/saturn-cli/internal/cli"
	"github.com/saturn-platform/saturn-cli/internal/output"
	"github.com/saturn-platform/saturn-cli/internal/service"
)

// ResultDisplay represents a deploy result for table display
type ResultDisplay struct {
	Message        string `json:"message"`
	DeploymentUUID string `json:"deployment_uuid"`
}

// NewUUIDCommand deploys a resource by UUID
func NewUUIDCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uuid <uuid>",
		Short: "Deploy by uuid",
		Args:  cli.ExactArgs(1, "<uuid>"),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			uuid := args[0]

			client, err := cli.GetAPIClient(cmd)
			if err != nil {
				return fmt.Errorf("failed to get API client: %w", err)
			}

			force, _ := cmd.Flags().GetBool("force")
			deploySvc := service.NewDeploymentService(client)

			sp := output.NewSpinner(fmt.Sprintf("Deploying %s...", uuid))
			result, err := deploySvc.Deploy(ctx, uuid, force)
			if err != nil {
				sp.Fail("Deployment failed")
				return fmt.Errorf("failed to deploy resource: %w", err)
			}

			if len(result.Deployments) == 0 {
				sp.Warning("Deployed (no deployment UUID returned)")
			} else {
				sp.Success(fmt.Sprintf("Deployed: %s", result.Deployments[0].DeploymentUUID))
			}

			format, _ := cmd.Flags().GetString("format")
			formatter, err := output.NewFormatter(format, output.Options{})
			if err != nil {
				return err
			}

			// For table format, convert deployment info array to display format
			if format == output.FormatTable {
				displays := make([]ResultDisplay, len(result.Deployments))
				for i, dep := range result.Deployments {
					displays[i] = ResultDisplay{
						Message:        dep.Message,
						DeploymentUUID: dep.DeploymentUUID,
					}
				}
				if err := formatter.Format(displays); err != nil {
					return err
				}
			} else {
				if err := formatter.Format(result); err != nil {
					return err
				}
			}

			// Show application URL if available (FQDN is *string)
			styler := output.DefaultStyler()
			appSvc := service.NewApplicationService(client)
			app, appErr := appSvc.Get(ctx, uuid)
			if appErr == nil && app.FQDN != nil && *app.FQDN != "" {
				styler.Info("Application URL:")
				for _, fqdn := range strings.Split(*app.FQDN, ",") {
					styler.URL("  " + strings.TrimSpace(fqdn))
				}
			}

			// Handle --wait flag
			return HandleWait(cmd, deploySvc, CollectDeploymentUUIDs(result))
		},
	}

	cmd.Flags().Bool("force", false, "Force deployment")
	AddWaitFlags(cmd)
	return cmd
}
