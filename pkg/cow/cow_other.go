//go:build !darwin

package cow

import (
	"fmt"
)

// Supported reports whether copy-on-write cloning is available on this platform.
func Supported() bool { return false }

// CloneDir is unsupported on non-darwin platforms.
func CloneDir(src, dst string) error {
	return fmt.Errorf("copy-on-write clone is only supported on macOS (APFS)")
}
