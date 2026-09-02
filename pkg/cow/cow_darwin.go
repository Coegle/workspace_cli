//go:build darwin

package cow

import (
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

// Supported reports whether copy-on-write cloning is available on this platform.
func Supported() bool { return true }

// CloneDir clones src into dst using APFS clonefile (copy-on-write).
//
// It only shares data blocks when src and dst live on the same APFS volume;
// clonefile fails with EXDEV across volumes. The clone is treated as a pure
// cache/warm-up: on any failure we return an error and the caller is expected
// to skip (never fall back to a full byte copy, which would defeat the purpose).
//
// dst must not already exist; the parent directory of dst is created if needed.
func CloneDir(src, dst string) error {
	info, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("donor not accessible: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("donor is not a directory: %s", src)
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return fmt.Errorf("failed to create parent dir for %s: %w", dst, err)
	}

	// clonefile requires the destination not to exist.
	if _, err := os.Lstat(dst); err == nil {
		if err := os.RemoveAll(dst); err != nil {
			return fmt.Errorf("failed to remove existing destination %s: %w", dst, err)
		}
	}

	// CLONE_NOFOLLOW: clone symlinks themselves rather than their targets,
	// which preserves the directory tree faithfully.
	if err := unix.Clonefile(src, dst, unix.CLONE_NOFOLLOW); err != nil {
		return fmt.Errorf("clonefile %s -> %s failed: %w", src, dst, err)
	}
	return nil
}
