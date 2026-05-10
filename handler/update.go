package handler

import (
	"fmt"
	"os"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/coegle/workspace_cli/core"
	"github.com/coegle/workspace_cli/pkg/updater"
)

func Update() error {
	fmt.Println("Checking for updates...")

	latest, err := updater.Update(core.Version)
	if err != nil {
		if strings.Contains(err.Error(), "rate limit exceeded") {
			_, _ = fmt.Fprintf(os.Stderr, "Update failed: GitHub API rate limit exceeded.\n")
			_, _ = fmt.Fprintf(os.Stderr, "Hint: You can increase the rate limit by setting the GITHUB_TOKEN environment variable.\n")
			_, _ = fmt.Fprintf(os.Stderr, "Example: export GITHUB_TOKEN=your_personal_access_token\n")
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "Update failed: %v\n", err)
		}
		return err
	}

	curr, err1 := semver.NewVersion(core.Version)
	latestVer, err2 := semver.NewVersion(latest.Version())

	// If both can be parsed as semver and are equal (semver ignores 'v' prefix automatically)
	if err1 == nil && err2 == nil && curr.Equal(latestVer) {
		fmt.Println("You are already using the latest version:", core.Version)
	} else if latest.Version() == core.Version {
		// Fallback to exact string match if semver parsing fails
		fmt.Println("You are already using the latest version:", core.Version)
	} else {
		fmt.Printf("Successfully updated to version %s\n", latest.Version())
		if latest.ReleaseNotes != "" {
			fmt.Println("Release Notes:\n", latest.ReleaseNotes)
		}
	}
	return nil
}
