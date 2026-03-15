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

func appendHook(base, hook string) string {
	trimmed := strings.TrimRight(base, "\r\n")
	if trimmed == "" {
		return hook
	}
	return trimmed + "\n\n" + hook
}

func stripHookBlock(content string) string {
	re := regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(HookComment) + `\r?\n[^\r\n]*\r?\n?`)
	return re.ReplaceAllString(content, "")
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
