package finder

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReorderCandidates_RespectsOrder(t *testing.T) {
	order := []string{".dircard", ".dircard.md"}
	got := ReorderCandidates(order)
	names := Names(got)

	if len(names) < 2 {
		t.Fatalf("unexpected candidates length: %d", len(names))
	}

	if names[0] != ".dircard" {
		t.Fatalf("expected first candidate to be .dircard, got %s", names[0])
	}
	if names[1] != ".dircard.md" {
		t.Fatalf("expected second candidate to be .dircard.md, got %s", names[1])
	}
}

func TestFindFilePath_CurrentDirOnlyAndRequireSection(t *testing.T) {
	tmp := t.TempDir()

	// create nested directory structure
	sub := filepath.Join(tmp, "sub", "dir")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}

	// create a .dircard in root
	rootDircard := filepath.Join(tmp, ".dircard")
	if err := os.WriteFile(rootDircard, []byte("root content"), 0o644); err != nil {
		t.Fatalf("failed to write root .dircard: %v", err)
	}

	// README.md in current dir with 'dircard' should be found first
	readme := filepath.Join(sub, "README.md")
	if err := os.WriteFile(readme, []byte("This file contains DirCard section"), 0o644); err != nil {
		t.Fatalf("failed to write README.md: %v", err)
	}

	p, err := FindFilePath(sub, 10, nil)
	if err != nil {
		t.Fatalf("expected to find README.md in current dir, got error: %v", err)
	}
	if filepath.Clean(p) != filepath.Clean(readme) {
		t.Fatalf("expected %s, got %s", readme, p)
	}

	// If README.md lacks the section, it should be skipped and root .dircard should be found
	if err := os.WriteFile(readme, []byte("no matching section here"), 0o644); err != nil {
		t.Fatalf("failed to overwrite README.md: %v", err)
	}

	p2, err := FindFilePath(sub, 10, nil)
	if err != nil {
		t.Fatalf("expected to find root .dircard, got error: %v", err)
	}
	if filepath.Clean(p2) != filepath.Clean(rootDircard) {
		t.Fatalf("expected %s, got %s", rootDircard, p2)
	}
}

func TestReorderCandidates_EmptyOrderReturnsDefault(t *testing.T) {
	got := ReorderCandidates(nil)
	if len(got) != len(Candidates) {
		t.Fatalf("expected default candidates length %d, got %d", len(Candidates), len(got))
	}
}

func TestReorderCandidates_UnknownNameIgnored(t *testing.T) {
	order := []string{"UNKNOWN", ".dircard"}
	got := ReorderCandidates(order)
	names := Names(got)
	// .dircard should still come before .dircard.md when requested
	foundIdx := -1
	for i, n := range names {
		if n == ".dircard" {
			foundIdx = i
			break
		}
	}
	if foundIdx == -1 {
		t.Fatalf(".dircard not found in reordered names")
	}
}

func TestFindFilePath_DepthLimitAndStopOnGit(t *testing.T) {
	tmp := t.TempDir()
	// /tmp/a/b/c
	deep := filepath.Join(tmp, "a", "b", "c")
	if err := os.MkdirAll(deep, 0o755); err != nil {
		t.Fatalf("failed to create deep: %v", err)
	}

	// put target in tmp (above a)
	target := filepath.Join(tmp, ".dircard")
	if err := os.WriteFile(target, []byte("root"), 0o644); err != nil {
		t.Fatalf("failed to write target: %v", err)
	}

	// depthLimit 1 should not find it from deep
	if _, err := FindFilePath(deep, 1, nil); err == nil {
		t.Fatalf("expected not to find file with small depth limit")
	}

	// depthLimit large should find
	if p, err := FindFilePath(deep, 10, nil); err != nil || filepath.Clean(p) != filepath.Clean(target) {
		t.Fatalf("expected to find %s, got %v %v", target, p, err)
	}

	// Create .git in a/b so search stops at that level
	gitDir := filepath.Join(tmp, "a", "b", ".git")
	if err := os.MkdirAll(gitDir, 0o755); err != nil {
		t.Fatalf("failed to create git dir: %v", err)
	}

	// From deep, even with large depth, should not find target above .git
	if _, err := FindFilePath(deep, 20, nil); err == nil {
		t.Fatalf("expected not to find file above .git directory")
	}
}

func TestFindFilePath_CustomCandidates(t *testing.T) {
	tmp := t.TempDir()
	sub := filepath.Join(tmp, "sub")
	if err := os.MkdirAll(sub, 0o755); err != nil {
		t.Fatalf("failed to create sub: %v", err)
	}

	custom := filepath.Join(tmp, "MYCARD")
	if err := os.WriteFile(custom, []byte("custom"), 0o644); err != nil {
		t.Fatalf("failed to write custom: %v", err)
	}

	candidates := []FileCandidate{{Name: "MYCARD"}}
	p, err := FindFilePath(sub, 10, candidates)
	if err != nil {
		t.Fatalf("expected to find custom candidate, got error: %v", err)
	}
	if filepath.Clean(p) != filepath.Clean(custom) {
		t.Fatalf("expected %s, got %s", custom, p)
	}
}
