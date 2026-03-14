package finder

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type FileCandidate struct {
	Name           string
	CurrentDirOnly bool
	RequireSection bool
}

var Candidates = []FileCandidate{
	{Name: ".dircard"},
	{Name: ".dircard.md"},
	{Name: "README", CurrentDirOnly: true, RequireSection: true},
	{Name: "README.md", CurrentDirOnly: true, RequireSection: true},
	{Name: ".dircard.local", CurrentDirOnly: true},
}

func FindFilePath(startDir string, depthLimit int) (string, error) {
	dir := startDir
	depth := 0

	for {
		if depthLimit > 0 && depth > depthLimit {
			break
		}

		for _, c := range Candidates {
			if c.CurrentDirOnly && depth > 0 {
				continue
			}
			p := filepath.Join(dir, c.Name)
			if !fileExists(p) {
				continue
			}
			if c.RequireSection {
				ok, _ := hasDircardSection(p)
				if !ok {
					continue
				}
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

func hasDircardSection(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(strings.ToLower(scanner.Text()), "dircard") {
			return true, nil
		}
	}
	return false, scanner.Err()
}
