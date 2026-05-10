package shell

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func ResolveSvc(search, searchDir string) (string, error) {
	if info, err := os.Stat(searchDir); err != nil || !info.IsDir() {
		return "", fmt.Errorf("error: directory %s does not exist", searchDir)
	}

	exactPath := filepath.Join(searchDir, search)
	if info, err := os.Stat(exactPath); err == nil && info.IsDir() {
		return search, nil
	}

	entries, err := os.ReadDir(searchDir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}

	var matches []string
	for _, entry := range entries {
		if entry.IsDir() && strings.Contains(entry.Name(), search) {
			matches = append(matches, entry.Name())
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("error: no service matching '%s' found in %s", search, searchDir)
	}

	if len(matches) == 1 {
		return matches[0], nil
	}

	// Check if running in a terminal
	if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) == 0 {
		return "", fmt.Errorf("error: multiple services match '%s' (%v), but running in non-interactive mode", search, matches)
	}

	_, _ = fmt.Fprintf(os.Stderr, "Multiple services match '%s'. Please select one:\n", search)
	for i, match := range matches {
		_, _ = fmt.Fprintf(os.Stderr, "%d) %s\n", i+1, match)
	}
	_, _ = fmt.Fprintf(os.Stderr, "%d) Cancel\n", len(matches)+1)

	reader := bufio.NewReader(os.Stdin)
	for {
		_, _ = fmt.Fprintf(os.Stderr, "Enter number: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)
		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(matches)+1 {
			_, _ = fmt.Fprintf(os.Stderr, "Invalid selection.\n")
			continue
		}

		if choice == len(matches)+1 {
			return "", fmt.Errorf("cancelled")
		}

		return matches[choice-1], nil
	}
}

func OpenGoLand(wsDir string) {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}

	appPath1 := "/Applications/GoLand.app"
	appPath2 := filepath.Join(home, "Applications/GoLand.app")

	if info, e := os.Stat(appPath1); e == nil && info.IsDir() {
		_ = exec.Command("open", "-na", "GoLand.app", "--args", wsDir).Start()
		return
	}
	if info, e := os.Stat(appPath2); e == nil && info.IsDir() {
		_ = exec.Command("open", "-na", "GoLand.app", "--args", wsDir).Start()
		return
	}

	if _, err = exec.LookPath("goland"); err == nil {
		_ = exec.Command("goland", wsDir).Start()
		return
	}
}
