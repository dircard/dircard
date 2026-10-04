package fileio

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCreateFile(t *testing.T) {
	for _, name := range []string{".dircard", ".dircard.md", "README.md", "README"} {
		t.Run(name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "nested", name)
			if err := CreateFile(path); err != nil {
				t.Fatal(err)
			}
			lines, err := ReadFileLines(path, true, 0, 0)
			if err != nil {
				t.Fatal(err)
			}
			if got := strings.Join(lines, "\n"); got != "# "+name+"\n\nAdd notes for this directory." {
				t.Fatalf("unexpected initial content: %q", got)
			}

			const existing = "# Existing notes\nKeep this content.\n"
			if err := os.WriteFile(path, []byte(existing), 0644); err != nil {
				t.Fatal(err)
			}
			if err := CreateFile(path); !errors.Is(err, fs.ErrExist) {
				t.Fatalf("expected file-exists error, got %v", err)
			}
			data, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != existing {
				t.Fatalf("existing content changed: %q", data)
			}
		})
	}
}
