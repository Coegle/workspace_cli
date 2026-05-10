package handler

import (
	"github.com/coegle/workspace_cli/core"
	"github.com/coegle/workspace_cli/pkg/git"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Destroy(input core.DestroyInput) {
	cfg := input.Cfg
	wsName := input.GetWorkspaceName()
	wsDir := filepath.Join(cfg.GetWorkspacePath(), wsName)

	if info, err := os.Stat(wsDir); err != nil || !info.IsDir() {
		_, _ = fmt.Fprintf(os.Stderr, "Workspace directory not found: %s\n", wsDir)
		return
	}

	entries, err := os.ReadDir(wsDir)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to read workspace directory: %v\n", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() {
			svc := entry.Name()
			targetPath := filepath.Join(wsDir, svc)

			// Check if it is a git repository
			if _, err := os.Stat(filepath.Join(targetPath, ".git")); os.IsNotExist(err) {
				fmt.Printf("Skipping git remove for %s: Not a git repository\n", svc)
				continue
			}

			fmt.Printf("Removing worktree for %s...\n", svc)

			if input.Force {
				out, e := git.RemoveFromWorkSpace(targetPath, targetPath, true)
				if e != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Failed to force remove worktree for %s: %s\n", svc, out)
				} else {
					fmt.Printf("Successfully removed %s\n", svc)
				}
				continue
			}

			e := git.CheckUncommittedChanges(targetPath)
			if errors.Is(e, git.ErrUncommittedChanges) {
				fmt.Printf("Warning: There are uncommitted changes in %s.\n", svc)
				fmt.Printf("Are you sure you want to force remove this worktree? (y/n) ")
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
			} else if e != nil {
				_, _ = fmt.Fprintf(os.Stderr, "Failed to check uncommitted changes for %s: %v\n", svc, e)
			} else {
				out, err := git.RemoveFromWorkSpace(targetPath, targetPath, false)
				if err != nil {
					_, _ = fmt.Fprintf(os.Stderr, "Failed to remove worktree for %s: %s\n", svc, out)
				} else {
					fmt.Printf("Successfully removed %s\n", svc)
				}
			}
		}
	}

	// Remove the whole directory
	if err = os.RemoveAll(wsDir); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to delete workspace directory: %v\n", err)
	} else {
		fmt.Printf("Workspace %s has been completely destroyed.\n", wsName)
	}
}
