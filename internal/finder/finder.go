package finder

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
)

type FileCandidate struct {
	Name           string
	CurrentDirOnly bool
}

var Candidates = []FileCandidate{
	{Name: ".dircard.md"},
	{Name: ".dircard"},
	{Name: "README.md", CurrentDirOnly: true},
	{Name: "README", CurrentDirOnly: true},
}

func ReorderCandidates(order []string) []FileCandidate {
	if len(order) == 0 {
		return Candidates
	}

	priority := make(map[string]int, len(order))
	for i, name := range order {
		priority[name] = i
	}

	result := append([]FileCandidate(nil), Candidates...)

	sort.SliceStable(result, func(i, j int) bool {
		pi, iok := priority[result[i].Name]
		pj, jok := priority[result[j].Name]

		if iok && jok {
			return pi < pj
		}
		if iok {
			return true
		}
		if jok {
			return false
		}
		return false
	})

	return result
}

func FindFilePath(startDir string, depthLimit int, candidates []FileCandidate) (string, error) {
	if len(candidates) == 0 {
		candidates = Candidates
	}
	dir := startDir
	depth := 0

	for {
		if depthLimit > 0 && depth > depthLimit {
			break
		}

		for _, c := range candidates {
			if c.CurrentDirOnly && depth > 0 {
				continue
			}
			p := filepath.Join(dir, c.Name)
			if !fileExists(p) {
				continue
			}
			return p, nil
		}

		// Stop searching when .git or .svn is found.
		if fileExists(filepath.Join(dir, ".git")) || fileExists(filepath.Join(dir, ".svn")) {
			break
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}

		dir = parent
		depth++
	}

	return "", errors.New("no dircard found")
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func Names(candidates []FileCandidate) []string {
	names := make([]string, len(candidates))
	for i, c := range candidates {
		names[i] = c.Name
	}
	return names
}
