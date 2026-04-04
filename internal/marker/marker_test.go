package marker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsReadme(t *testing.T) {
	cases := []struct {
		path string
		want bool
	}{
		{"README", true},
		{"README.md", true},
		{".dircard", false},
		{".dircard.md", false},
		{filepath.Join("some", "dir", "README"), true},
	}
	for _, c := range cases {
		if got := IsReadme(c.path); got != c.want {
			t.Errorf("IsReadme(%q) = %v, want %v", c.path, got, c.want)
		}
	}
}

func TestContains(t *testing.T) {
	withMarkers := "hello\n" + Start + "\nfoo\n" + End + "\nworld\n"
	if !Contains(withMarkers) {
		t.Error("Contains: expected true for content with markers")
	}
	if Contains("no markers here") {
		t.Error("Contains: expected false for content without markers")
	}
}

func TestInitialContent(t *testing.T) {
	cases := []struct {
		path       string
		wantPrefix string
		wantMarker bool
	}{
		{"README", "# README", true},
		{"README.md", "# README", true},
		{".dircard", "# .dircard", false},
		{".dircard.md", "# .dircard.md", false},
	}
	for _, c := range cases {
		got := string(InitialContent(c.path))
		if len(got) < len(c.wantPrefix) || got[:len(c.wantPrefix)] != c.wantPrefix {
			t.Errorf("InitialContent(%q) prefix: got %q, want prefix %q", c.path, got, c.wantPrefix)
		}
		if Contains(got) != c.wantMarker {
			t.Errorf("InitialContent(%q) markers: got %v, want %v", c.path, Contains(got), c.wantMarker)
		}
	}
}

func TestAppendIfNeeded(t *testing.T) {
	dir := t.TempDir()
	p := filepath.Join(dir, "README")
	if err := os.WriteFile(p, []byte("existing content\n"), 0644); err != nil {
		t.Fatalf("setup: %v", err)
	}

	if err := AppendIfNeeded(p); err != nil {
		t.Fatalf("first append: %v", err)
	}
	data, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !Contains(string(data)) {
		t.Fatalf("markers not found after append: %q", data)
	}

	// idempotency
	if err := AppendIfNeeded(p); err != nil {
		t.Fatalf("second append: %v", err)
	}
	data2, _ := os.ReadFile(p)
	if string(data) != string(data2) {
		t.Fatalf("second append changed file: before=%q after=%q", data, data2)
	}
}

func TestCreateOrUpdate(t *testing.T) {
	t.Run("creates new file", func(t *testing.T) {
		dir := t.TempDir()
		p := filepath.Join(dir, ".dircard")
		action, err := CreateOrUpdate(p, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if action != ActionCreated {
			t.Errorf("action = %q, want %q", action, ActionCreated)
		}
		if _, err := os.Stat(p); err != nil {
			t.Errorf("file not created: %v", err)
		}
	})

	t.Run("errors on existing file without force", func(t *testing.T) {
		dir := t.TempDir()
		p := filepath.Join(dir, ".dircard")
		_ = os.WriteFile(p, []byte("x"), 0644)
		_, err := CreateOrUpdate(p, false)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("overwrites existing non-README with force", func(t *testing.T) {
		dir := t.TempDir()
		p := filepath.Join(dir, ".dircard")
		_ = os.WriteFile(p, []byte("old"), 0644)
		action, err := CreateOrUpdate(p, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if action != ActionUpdated {
			t.Errorf("action = %q, want %q", action, ActionUpdated)
		}
	})

	t.Run("appends to existing README with force", func(t *testing.T) {
		dir := t.TempDir()
		p := filepath.Join(dir, "README")
		_ = os.WriteFile(p, []byte("# My Project\n"), 0644)
		action, err := CreateOrUpdate(p, true)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if action != ActionAppended {
			t.Errorf("action = %q, want %q", action, ActionAppended)
		}
		data, _ := os.ReadFile(p)
		if !Contains(string(data)) {
			t.Errorf("markers not found after append: %q", data)
		}
	})

	t.Run("creates nested path", func(t *testing.T) {
		dir := t.TempDir()
		p := filepath.Join(dir, "subdir", ".dircard")
		action, err := CreateOrUpdate(p, false)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if action != ActionCreated {
			t.Errorf("action = %q, want %q", action, ActionCreated)
		}
	})
}
