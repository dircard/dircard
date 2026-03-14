package fileio

import (
	"bufio"
	"os"
	"strings"
)

func ReadLines(path string, start, count int) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	lines := make([]string, 0, count)
	lineNum := 0

	for scanner.Scan() {
		if lineNum >= start {
			lines = append(lines, scanner.Text())
			if count > 0 && len(lines) >= count {
				break
			}
		}
		lineNum++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func ReadAllLines(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	s := strings.TrimPrefix(string(data), "\uFEFF")
	return strings.Split(strings.TrimRight(s, "\n"), "\n"), nil
}

func ReadFileLines(path string, full bool, start, count int) ([]string, error) {
	if full {
		return ReadAllLines(path)
	}
	return ReadLines(path, start, count)
}
