package handler

import (
	"coegle/workspace_cli/core"
	"coegle/workspace_cli/pkg/git"
	"coegle/workspace_cli/pkg/shell"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Remove(input core.RemoveInput) {
	cfg := input.Cfg
	wsName := input.GetWorkspaceName()
	wsDir := filepath.Join(cfg.GetWorkspacePath(), wsName)

	for _, svcInput := range input.ServiceRepoNames {
		svc, err := shell.ResolveSvc(svcInput, wsDir)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
			continue
		}

		targetPath := filepath.Join(wsDir, svc)
		fmt.Printf("Removing %s from %s\n", svc, wsName)

		if info, err := os.Stat(targetPath); err == nil && info.IsDir() {
			// Check if it is a git repository
			if _, err := os.Stat(filepath.Join(targetPath, ".git")); os.IsNotExist(err) {
				fmt.Printf("Skipping %s: Not a git repository (or missing .git dir)\n", svc)
				continue
			}

			err := git.CheckUncommittedChanges(targetPath)
			if err == git.ErrUncommittedChanges {
				fmt.Printf("Warning: There are uncommitted changes in %s.\n", svc)
				fmt.Printf("Are you sure you want to force remove the worktree? (y/n) ")
				var confirm string
				_, _ = fmt.Scanln(&confirm)
				if strings.ToLower(confirm) != "y" {
					fmt.Printf("Skipped %s\n", svc)
					continue
				}

				out, err := git.RemoveFromWorkSpace(targetPath, targetPath, true)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Failed to force remove worktree for %s: %s\n", svc, out)
				} else {
					fmt.Printf("Successfully removed %s\n", svc)
				}
			} else if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Failed to check uncommitted changes for %s: %v\n", svc, err)
			} else {
				out, err := git.RemoveFromWorkSpace(targetPath, targetPath, false)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Failed to remove worktree for %s: %s\n", svc, out)
				} else {
					fmt.Printf("Successfully removed %s\n", svc)
				}
			}
		} else {
			fmt.Printf("Not found: %s\n", targetPath)
		}
	}

	// Check if the entire workspace directory contains only metadata files (e.g., .ws_branch), if so, remove it
	entries, err := os.ReadDir(wsDir)
	if err == nil {
		isEmpty := true
		for _, entry := range entries {
			if entry.Name() != ".ws_branch" {
				isEmpty = false
				break
			}
		}
		if isEmpty {
			if err := os.RemoveAll(wsDir); err == nil {
				fmt.Printf("Workspace %s is now empty, removed directory.\n", wsName)
			}
		}
	}
}
