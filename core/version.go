package core

import (
	"runtime/debug"
)

// Version is the current version of the application.
// It will be overwritten by goreleaser's ldflags during release builds.
var Version = "v0.0.1"

func init() {
	// If Version is still the default, try to read it from go module build info
	if Version == "v0.0.1" {
		if info, ok := debug.ReadBuildInfo(); ok {
			// When installed via `go install package@version`, Main.Version contains the version
			if info.Main.Version != "" && info.Main.Version != "(devel)" {
				Version = info.Main.Version
			}
		}
	}
}
