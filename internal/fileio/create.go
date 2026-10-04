package fileio

import (
	"fmt"
	"os"
	"path/filepath"
)

// CreateFile creates a notes file without modifying existing files.
func CreateFile(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}

	_, writeErr := fmt.Fprintf(f, "# %s\n\nAdd notes for this directory.\n", filepath.Base(path))
	closeErr := f.Close()
	if writeErr != nil {
		return writeErr
	}
	return closeErr
}
