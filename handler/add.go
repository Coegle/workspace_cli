package handler

import (
	"github.com/coegle/workspace_cli/core"
	"github.com/coegle/workspace_cli/pkg/git"
	"github.com/coegle/workspace_cli/pkg/shell"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Add(input core.AddInput) error {
	cfg := input.Cfg
	featBranchName := input.FeatBranchName
	serviceRepoNames := input.ServiceRepoNames
	baseBranchName := input.BaseBranchName
	wsName := input.GetWorkspaceName()
	wsDir := filepath.Join(cfg.GetWorkspacePath(), wsName)
	metaFile := filepath.Join(wsDir, ".ws_branch")

	if info, err := os.Stat(wsDir); err == nil && info.IsDir() {
		if b, err := os.ReadFile(metaFile); err == nil {
			existingBranch := strings.TrimSpace(string(b))
			if existingBranch != featBranchName {
				return fmt.Errorf("🚨 ERROR: Workspace directory '%s' already exists and is bound to a different branch ('%s').\n"+
					"         This usually happens when different branch names (e.g., 'feat/a' and 'feat_a') map to the same directory name.\n"+
					"         Please use a different branch name or clean up the existing workspace first.", wsName, existingBranch)
			}
		} else {
			fmt.Printf("Already exists: %s (Warning: missing .ws_branch metadata)\n", wsDir)
		}
	} else if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(wsDir, 0755); err != nil {
			return fmt.Errorf("failed to create workspace directory: %w", err)
		}
		if err := os.WriteFile(metaFile, []byte(featBranchName+"\n"), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to write metadata file: %v\n", err)
		}
	} else {
		return fmt.Errorf("failed to check workspace directory: %w", err)
	}

	fmt.Printf("Creating workspace: %s (Branch: %s)\n", wsName, featBranchName)

	if len(cfg.Symlinks) > 0 {
		fmt.Println("Setting up symlinks...")
		for _, symlink := range cfg.Symlinks {
			src := symlink.Source
			dest := symlink.Destination

			// Resolve ~ in source path
			if strings.HasPrefix(src, "~/") {
				if home, err := os.UserHomeDir(); err == nil {
					src = filepath.Join(home, src[2:])
				}
			}

			destPath := filepath.Join(wsDir, dest)

			// Create parent directory for destination if it doesn't exist
			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				fmt.Fprintf(os.Stderr, "  -> Failed to create directory for symlink %s: %v\n", destPath, err)
				continue
			}

			// Check if source exists
			if _, err := os.Stat(src); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "  -> Warning: Source path for symlink does not exist: %s\n", src)
				continue
			}

			// Remove existing destination if it exists
			if _, err := os.Lstat(destPath); err == nil {
				os.RemoveAll(destPath)
			}

			if err := os.Symlink(src, destPath); err != nil {
				fmt.Fprintf(os.Stderr, "  -> Failed to create symlink %s -> %s: %v\n", src, destPath, err)
			} else {
				fmt.Printf("  -> Created symlink: %s -> %s\n", destPath, src)
			}
		}
	}

	baseRepos := cfg.GetReposPath()

	for _, serviceName := range serviceRepoNames {
		svc, err := shell.ResolveSvc(serviceName, baseRepos)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}

		repoPath := filepath.Join(baseRepos, svc)
		targetPath := filepath.Join(wsDir, svc)

		if info, err := os.Stat(repoPath); err != nil || !info.IsDir() {
			fmt.Fprintf(os.Stderr, "Repo not found: %s\n", svc)
			continue
		}

		currentBaseBranch := baseBranchName
		if currentBaseBranch == "" {
			currentBaseBranch, err = git.GetDefaultBranch(repoPath)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Failed to get default branch for %s: %v\n", svc, err)
				continue
			}
		}

		fmt.Printf("  -> Adding worktree for %s\n", svc)
		out, err := git.AddOrCreateBranchToWorkSpace(repoPath, targetPath, featBranchName, currentBaseBranch)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to add worktree for %s: %s\n", svc, out)
			continue
		}
	}

	return nil
}
