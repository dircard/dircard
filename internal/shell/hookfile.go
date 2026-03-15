package shell

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const HookComment = "# ---- dircard hook ----"

var hookBlockRe = regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(HookComment) + `\r?\n[^\r\n]*\r?\n?`)

func AppendHookToFile(path, hook string, force bool) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}

	current := string(data)
	hasHook := strings.Contains(current, HookComment)
	if hasHook && !force {
		return false, nil
	}

	updated := current
	if hasHook && force {
		updated = appendHook(stripHookBlock(current), hook)
	} else if !hasHook {
		updated = appendHook(current, hook)
	}

	if err := os.WriteFile(path, []byte(updated), 0644); err != nil {
		return false, err
	}
	return true, nil
}

func RemoveHookFromFile(path string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	current := string(data)
	if !strings.Contains(current, HookComment) {
		return false, nil
	}

	updated := strings.TrimRight(stripHookBlock(current), "\r\n")
	if updated != "" {
		updated += "\n"
	}

	if err := os.WriteFile(path, []byte(updated), 0644); err != nil {
		return false, err
	}
	return true, nil
}

func appendHook(base, hook string) string {
	trimmed := strings.TrimRight(base, "\r\n")
	if trimmed == "" {
		return hook
	}
	return trimmed + "\n\n" + hook
}

func stripHookBlock(content string) string {
	return hookBlockRe.ReplaceAllString(content, "")
}

func getDocumentsDir() string {
	if runtime.GOOS == "windows" {
		if _, err := exec.LookPath("powershell"); err == nil {
			out, err := exec.Command(
				"powershell",
				"-NoProfile",
				"-Command",
				"[Console]::OutputEncoding=[System.Text.Encoding]::UTF8; [System.Environment]::GetFolderPath('MyDocuments')",
			).Output()
			if err == nil {
				p := strings.TrimSpace(string(out))
				if p != "" {
					return p
				}
			}
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	return filepath.Join(home, "Documents")
}

func GetPowerShellProfiles() []string {
	if runtime.GOOS == "windows" {
		docs := getDocumentsDir()
		if docs == "" {
			return nil
		}

		return []string{
			filepath.Join(docs, "WindowsPowerShell", "Microsoft.PowerShell_profile.ps1"),
			filepath.Join(docs, "PowerShell", "Microsoft.PowerShell_profile.ps1"),
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	return []string{
		filepath.Join(home, ".config", "powershell", "Microsoft.PowerShell_profile.ps1"),
	}
}
