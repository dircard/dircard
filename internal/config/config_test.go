package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveLoad_UsesUserConfigDirEnv(t *testing.T) {
	tmp := t.TempDir()
	// Prefer XDG_CONFIG_HOME for predictable behavior on CI
	if err := os.Setenv("XDG_CONFIG_HOME", tmp); err != nil {
		t.Fatalf("failed to set XDG_CONFIG_HOME: %v", err)
	}
	defer os.Unsetenv("XDG_CONFIG_HOME")

	cfg := &Config{CandidateOrder: []string{".dircard", "README.md"}}
	if err := cfg.Save(); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	// Ensure file exists
	path, err := configFilePath()
	if err != nil {
		t.Fatalf("configFilePath failed: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected config file to exist at %s: %v", path, err)
	}

	loaded := Load()
	if len(loaded.CandidateOrder) != 2 || loaded.CandidateOrder[0] != ".dircard" {
		t.Fatalf("loaded config mismatch: %#v", loaded)
	}

	// cleanup created file
	os.Remove(path)
	os.RemoveAll(filepath.Dir(path))
}
