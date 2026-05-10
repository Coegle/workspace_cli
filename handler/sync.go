package handler

import (
	"github.com/coegle/workspace_cli/core"
	"github.com/coegle/workspace_cli/pkg/git"
	"github.com/coegle/workspace_cli/pkg/shell"
	"fmt"
	"os"
	"path/filepath"
)

func Sync(input core.SyncInput) {
	cfg := input.Cfg
	reposDir := cfg.GetReposPath()
	targetBranch := input.BaseBranchName

	for _, svcInput := range input.ServiceRepoNames {
		svc, err := shell.ResolveSvc(svcInput, reposDir)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error resolving service %s: %v\n", svcInput, err)
			continue
		}

		repoPath := filepath.Join(reposDir, svc)

		// Ensure it's a directory (and likely a git repo)
		if info, err := os.Stat(repoPath); err != nil || !info.IsDir() {
			_, _ = fmt.Fprintf(os.Stderr, "Service repository not found or invalid: %s\n", repoPath)
			continue
		}
		if targetBranch == "" {
			targetBranch, err = git.GetDefaultBranch(repoPath)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Failed to get default branch for %s: %v\n", svc, err)
				continue
			}
		}

		fmt.Printf("Syncing %s base branch '%s'...\n", svc, targetBranch)

		currentBranch, err := git.GetCurrentBranch(repoPath)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to get current branch for %s: %v\n", svc, err)
			continue
		}

		var out []byte
		var syncErr error

		if currentBranch == targetBranch {
			// The target branch is currently checked out, so we pull with fast-forward only
			out, syncErr = git.PullFastForward(repoPath, targetBranch)
		} else {
			// The target branch is not currently checked out, we can safely fetch and update it
			out, syncErr = git.FetchBranch(repoPath, targetBranch)
		}

		if syncErr != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to sync %s: %v\nOutput:\n%s\n", svc, syncErr, out)
		} else {
			fmt.Printf("Successfully synced %s.\n", svc)
			if len(out) > 0 {
				// Optionally print the output if it's not empty, it might contain useful info like what was updated
				fmt.Print(string(out))
			}
		}
	}
}
