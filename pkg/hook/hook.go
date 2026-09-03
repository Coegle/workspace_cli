// Package hook runs user-defined lifecycle scripts by convention.
//
// Scripts live under $HOME/<ConfigPath>/hooks/ (e.g. ~/.ws/hooks/post-add) and
// are opt-in: if the script doesn't exist it's silently skipped. Context is
// passed to the script purely through environment variables (WS_* keys), so the
// script can do whatever it wants without ws needing to know about it.
package hook

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/coegle/workspace_cli/core"
)

// Run executes a "post-*" convention hook (e.g. "post-add"). A missing hook is
// skipped silently and any failure is reported to stderr but never aborts ws,
// since the main operation has already completed by the time it runs.
func Run(hookName string, env map[string]string, workDir string) {
	if err := run(hookName, env, workDir); err != nil {
		fmt.Fprintf(os.Stderr, "  -> Hook %s failed: %v\n", hookName, err)
	}
}

// RunPre executes a "pre-*" convention hook (e.g. "pre-add") and returns its
// error so the caller can abort the operation when the hook fails. A missing
// hook returns nil (nothing to gate).
func RunPre(hookName string, env map[string]string, workDir string) error {
	if err := run(hookName, env, workDir); err != nil {
		return fmt.Errorf("hook %s failed: %w", hookName, err)
	}
	return nil
}

// run locates the convention-based hook script named hookName under
// $HOME/<ConfigPath>/hooks/, injects env on top of the current process
// environment, and runs it with workDir as its working directory (when it
// exists). Returns nil if the hook is undefined or exits successfully.
func run(hookName string, env map[string]string, workDir string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}
	hookPath := filepath.Join(home, core.ConfigPath, "hooks", hookName)

	info, err := os.Stat(hookPath)
	if err != nil || info.IsDir() {
		// Hook not defined -> nothing to do.
		return nil
	}
	if info.Mode()&0111 == 0 {
		fmt.Fprintf(os.Stderr, "  -> Hook %s exists but is not executable, skipping (run: chmod +x %s)\n", hookName, hookPath)
		return nil
	}

	fmt.Printf("  -> Running hook: %s\n", hookName)

	cmd := exec.Command(hookPath)
	// pre-* hooks may run before workDir exists yet; only set it when valid,
	// otherwise exec would fail with "no such file or directory".
	if info, err := os.Stat(workDir); err == nil && info.IsDir() {
		cmd.Dir = workDir
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	// Inherit the full environment first, then layer WS_* context on top so the
	// script still sees PATH/HOME and can locate git, code, etc.
	cmd.Env = os.Environ()
	for k, v := range env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}

	return cmd.Run()
}
