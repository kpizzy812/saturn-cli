package deployment

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	"github.com/saturn-platform/saturn-cli/internal/output"
	"github.com/saturn-platform/saturn-cli/internal/service"
)

const (
	// ExitCodeSuccess means all deployments finished successfully
	ExitCodeSuccess = 0
	// ExitCodeFailed means one or more deployments failed/cancelled/timed-out
	ExitCodeFailed = 1
	// ExitCodeWaitTimeout means the --timeout was exceeded while waiting
	ExitCodeWaitTimeout = 2
)

// AddWaitFlags adds --wait, --timeout, and --poll-interval flags to a deploy command
func AddWaitFlags(cmd *cobra.Command) {
	cmd.Flags().BoolP("wait", "w", false, "Wait for deployment to complete before exiting")
	cmd.Flags().Int("timeout", 600, "Timeout in seconds when using --wait (default 600)")
	cmd.Flags().Int("poll-interval", 3, "Poll interval in seconds when using --wait (default 3)")
}

// HandleWait checks if --wait was set and blocks until all deployments complete.
// Returns nil if --wait was not set. Returns an error if deployments failed or timed out.
func HandleWait(cmd *cobra.Command, deploySvc *service.DeploymentService, deploymentUUIDs []string) error {
	wait, _ := cmd.Flags().GetBool("wait")
	if !wait || len(deploymentUUIDs) == 0 {
		return nil
	}

	timeoutSec, _ := cmd.Flags().GetInt("timeout")
	pollSec, _ := cmd.Flags().GetInt("poll-interval")

	ctx, cancel := context.WithTimeout(cmd.Context(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	pollInterval := time.Duration(pollSec) * time.Second

	sp := output.NewSpinner(fmt.Sprintf("Waiting for %d deployment(s) (timeout: %ds)...", len(deploymentUUIDs), timeoutSec))

	// Track last printed status per UUID to avoid spamming
	lastStatus := make(map[string]string)

	onStatus := func(uuid, status string) {
		if lastStatus[uuid] != status {
			lastStatus[uuid] = status
			sp.UpdateText(fmt.Sprintf("Waiting for %d deployment(s): [%s] %s", len(deploymentUUIDs), uuid, status))
		}
	}

	results, err := deploySvc.WaitForMultiple(ctx, deploymentUUIDs, pollInterval, onStatus)

	// Determine overall success
	allSuccess := true
	for _, res := range results {
		if !res.Finished {
			allSuccess = false
		}
	}

	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			sp.Fail(fmt.Sprintf("Wait timeout exceeded (%ds)", timeoutSec))
			return fmt.Errorf("wait timeout exceeded (%ds), exit code %d", timeoutSec, ExitCodeWaitTimeout)
		}
		sp.Fail(fmt.Sprintf("Error waiting for deployments: %v", err))
		return fmt.Errorf("error waiting for deployments: %w", err)
	}

	if !allSuccess {
		// Print per-deployment status for failures
		for _, res := range results {
			if !res.Finished {
				fmt.Fprintf(cmd.ErrOrStderr(), "  [%s] %s\n", res.DeploymentUUID, res.Status)
			}
		}
		sp.Fail("One or more deployments did not finish successfully")
		return fmt.Errorf("one or more deployments did not finish successfully")
	}

	sp.Success(fmt.Sprintf("All %d deployment(s) completed successfully", len(deploymentUUIDs)))
	return nil
}

// CollectDeploymentUUIDs extracts deployment UUIDs from a DeployResponse
func CollectDeploymentUUIDs(result *service.DeployResponse) []string {
	uuids := make([]string, 0, len(result.Deployments))
	for _, dep := range result.Deployments {
		if dep.DeploymentUUID != "" {
			uuids = append(uuids, dep.DeploymentUUID)
		}
	}
	return uuids
}
