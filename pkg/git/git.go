package git

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

func BranchExist(branchName string, repoDir string) bool {
	cmd := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+branchName)
	cmd.Dir = repoDir
	err := cmd.Run()
	return err == nil
}

func AddOrCreateBranchToWorkSpace(repoDir string, targetPath string, branchName string, baseBranchName string) ([]byte, error) {
	if BranchExist(branchName, repoDir) {
		cmd := exec.Command("git", "worktree", "add", targetPath, branchName)
		cmd.Dir = repoDir
		return cmd.CombinedOutput()
	}

	cmd := exec.Command("git", "worktree", "add", "-b", branchName, targetPath, baseBranchName)
	cmd.Dir = repoDir
	out, err := cmd.CombinedOutput()
	if err == nil {
		return out, nil
	}

	// If creating branch from local baseBranchName fails, try creating from remote baseBranchName directly
	if !strings.Contains(baseBranchName, "/") {
		remoteBranch := "origin/" + baseBranchName
		cmdRemote := exec.Command("git", "worktree", "add", "-b", branchName, targetPath, remoteBranch)
		cmdRemote.Dir = repoDir
		if outRemote, errRemote := cmdRemote.CombinedOutput(); errRemote == nil {
			return outRemote, nil
		}
	}

	return out, err
}

func RemoveFromWorkSpace(repoDir string, targetPath string, force bool) ([]byte, error) {
	cmd := exec.Command("git", "worktree", "remove", targetPath)
	if force {
		cmd = exec.Command("git", "worktree", "remove", "-f", targetPath)
	}
	cmd.Dir = repoDir
	return cmd.CombinedOutput()
}

var ErrUncommittedChanges = errors.New("uncommitted changes exist in the repository")

// CheckUncommittedChanges returns ErrUncommittedChanges if there are uncommitted changes
func CheckUncommittedChanges(repoDir string) error {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	if string(out) != "" {
		return ErrUncommittedChanges
	}
	return nil
}

// GetCurrentBranch returns the name of the currently checked out branch
func GetCurrentBranch(repoDir string) (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = repoDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %v, output: %s", err, string(out))
	}
	return strings.TrimSpace(string(out)), nil
}

// FetchBranch updates a local branch from origin without checking it out
func FetchBranch(repoDir, branchName string) ([]byte, error) {
	cmd := exec.Command("git", "fetch", "origin", fmt.Sprintf("%s:%s", branchName, branchName))
	cmd.Dir = repoDir
	return cmd.CombinedOutput()
}

// PullFastForward updates the current branch using fast-forward only
func PullFastForward(repoDir, branchName string) ([]byte, error) {
	cmd := exec.Command("git", "pull", "origin", branchName, "--ff-only")
	cmd.Dir = repoDir
	return cmd.CombinedOutput()
}

// GetDefaultBranch returns the default branch name of the remote origin (e.g., master or main)
func GetDefaultBranch(repoDir string) (string, error) {
	cmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD")
	cmd.Dir = repoDir
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get default branch: %v, output: %s", err, string(out))
	}

	// Output looks like: refs/remotes/origin/main
	// We want to extract just the branch name "main"
	ref := strings.TrimSpace(string(out))
	parts := strings.Split(ref, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1], nil
	}

	return "master", nil // fallback
}
