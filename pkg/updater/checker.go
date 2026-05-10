package updater

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/coegle/workspace_cli/core"

	"github.com/Masterminds/semver/v3"
	"github.com/creativeprojects/go-selfupdate"
	"github.com/fatih/color"
)

type updateCache struct {
	LastCheck     time.Time `json:"last_check"`
	LatestVersion string    `json:"latest_version"`
}

func getCachePath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, core.ConfigPath, core.UpdateCacheFile)
}

// CheckAndNotice silently checks for updates (once a day) and prompts the user if a new version is available.
func CheckAndNotice() {
	currentVer := core.Version

	cachePath := getCachePath()
	var cache updateCache

	data, err := os.ReadFile(cachePath)
	if err == nil {
		_ = json.Unmarshal(data, &cache)
	}

	// Print notice and prompt if a newer version is found in cache
	if cache.LatestVersion != "" {
		curr, err1 := semver.NewVersion(currentVer)
		latest, err2 := semver.NewVersion(cache.LatestVersion)
		if err1 == nil && err2 == nil && latest.GreaterThan(curr) {
			notice := color.New(color.FgYellow).Sprintf("\n🚀 A new version of ws is available: %s -> %s\n", core.Version, cache.LatestVersion)
			_, _ = fmt.Fprint(os.Stderr, notice)

			// Check if running in an interactive terminal before prompting
			if fileInfo, _ := os.Stdout.Stat(); (fileInfo.Mode() & os.ModeCharDevice) != 0 {
				_, _ = fmt.Fprint(os.Stderr, "Do you want to update now? [y/N]: ")

				reader := bufio.NewReader(os.Stdin)
				input, _ := reader.ReadString('\n')
				input = strings.ToLower(strings.TrimSpace(strings.ToLower(input)))

				if input == "y" {
					_, _ = fmt.Println("Updating...")
					release, e := Update(core.Version)
					if e != nil {
						_, _ = fmt.Fprintf(os.Stderr, "Update failed: %v\n", e)
						return
					}
					_, _ = fmt.Printf("Successfully updated to version %s\n", release.Version())
					if release.ReleaseNotes != "" {
						_, _ = fmt.Println("Release Notes:\n", release.ReleaseNotes)
					}
				}
			} else {
				// Non-interactive mode, just print instructions
				_, _ = fmt.Fprintln(os.Stderr, "Run `ws update` to update manually.")
			}
		}
	}
	UpdateCache(cache, cachePath)
}

func UpdateCache(cache updateCache, cachePath string) {
	// Check GitHub if cache is older than 24 hours
	if time.Since(cache.LastCheck) > 24*time.Hour {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		up, e := selfupdate.NewUpdater(selfupdate.Config{})
		if e == nil {
			if latest, found, _ := up.DetectLatest(ctx, selfupdate.ParseSlug(RepoSlug)); found {
				cache.LatestVersion = latest.Version()
			}
		}
		cache.LastCheck = time.Now()

		_ = os.MkdirAll(filepath.Dir(cachePath), 0755)
		data, _ := json.Marshal(cache)
		_ = os.WriteFile(cachePath, data, 0644)
	}
}
