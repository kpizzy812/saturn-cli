package version

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"sort"
	"time"

	compareVersion "github.com/hashicorp/go-version"
)

// version is injected by GoReleaser at build time via ldflags.
// If not set (local build), falls back to Go module version or "dev".
var version string

// GitHubAPIURL is the URL for fetching CLI version tags (exported for testing)
var GitHubAPIURL = "https://api.github.com/repos/kpizzy812/saturn-cli/git/refs/tags"

func GetVersion() string {
	if version != "" {
		return version
	}
	// Fallback: go install sets module version in build info
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}
	return "dev"
}

// Tag represents a git tag for version checking
type Tag struct {
	Ref string `json:"ref"`
}

// CheckLatestVersionOfCli checks for CLI updates on every command.
// Errors are handled silently - the function returns without printing anything
// if the GitHub API call fails.
func CheckLatestVersionOfCli(_ bool) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", GitHubAPIURL, nil)
	if err != nil {
		return "", nil // Silent fail
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", nil // Silent fail
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", nil // Silent fail
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", nil // Silent fail
	}

	var tags []Tag
	if err := json.Unmarshal(body, &tags); err != nil {
		return "", nil // Silent fail
	}

	if len(tags) == 0 {
		return "", nil // Silent fail
	}

	versionsRaw := make([]string, 0, len(tags))
	for _, tag := range tags {
		if len(tag.Ref) > 10 {
			versionStr := tag.Ref[10:]
			versionsRaw = append(versionsRaw, versionStr)
		}
	}

	if len(versionsRaw) == 0 {
		return "", nil // Silent fail
	}

	versions := make([]*compareVersion.Version, 0, len(versionsRaw))
	for _, raw := range versionsRaw {
		v, err := compareVersion.NewVersion(raw)
		if err != nil {
			continue // Skip invalid versions
		}
		versions = append(versions, v)
	}

	if len(versions) == 0 {
		return "", nil // Silent fail
	}

	sort.Sort(compareVersion.Collection(versions))
	latestVersion := versions[len(versions)-1]

	// Compare versions properly using semantic versioning
	currentVersion, err := compareVersion.NewVersion(GetVersion())
	if err != nil {
		return "", nil // Silent fail
	}

	if latestVersion.GreaterThan(currentVersion) {
		fmt.Printf("A new version (%s) is available. Update with: saturn update\n", latestVersion.String())
	}
	return latestVersion.String(), nil
}
