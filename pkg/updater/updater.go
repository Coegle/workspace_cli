package updater

import (
	"context"

	"github.com/creativeprojects/go-selfupdate"
)

const RepoSlug = "coegle/workspace_cli"

// Update checks for updates and updates the binary if a newer version is found.
// It returns the latest Release info and any error encountered.
func Update(currentVersion string) (*selfupdate.Release, error) {
	// go-selfupdate will automatically use GITHUB_TOKEN if it's in the environment
	config := selfupdate.Config{
		Filters: []string{"^ws$|^ws.exe$"}, // Tell selfupdate to look for 'ws' (Linux/macOS) or 'ws.exe' (Windows) inside the archive
	}

	updater, err := selfupdate.NewUpdater(config)
	if err != nil {
		return nil, err
	}

	latest, err := updater.UpdateSelf(context.Background(), currentVersion, selfupdate.ParseSlug(RepoSlug))
	if err != nil {
		return nil, err
	}

	return latest, nil
}
