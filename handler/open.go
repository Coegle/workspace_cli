package handler

import (
	"github.com/coegle/workspace_cli/core"
	"github.com/coegle/workspace_cli/pkg/shell"
	"fmt"
	"os"
	"path/filepath"
)

func Open(input core.OpenInput) error {
	cfg := input.Cfg
	wsName := input.GetWorkspaceName()
	wsDir := filepath.Join(cfg.GetWorkspacePath(), wsName)

	if info, err := os.Stat(wsDir); err != nil || !info.IsDir() {
		if err != nil {
			return err
		}
		return fmt.Errorf("workspace directory not found: %s", wsDir)
	}

	fmt.Printf("Opening workspace '%s' in GoLand...\n", wsName)
	shell.OpenGoLand(wsDir)
	return nil
}
