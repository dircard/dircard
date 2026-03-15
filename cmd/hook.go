package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

const bashHook = `_dircard_last=""
_dircard_hook() {
  if [ "$PWD" != "$_dircard_last" ]; then
    command dircard show 2>/dev/null
    _dircard_last="$PWD"
  fi
}
case "$PROMPT_COMMAND" in
  *_dircard_hook*) ;;
  *) PROMPT_COMMAND="_dircard_hook${PROMPT_COMMAND:+;$PROMPT_COMMAND}" ;;
esac`

const zshHook = `_dircard_hook() {
  command dircard show 2>/dev/null
}
if ! (( ${chpwd_functions[(I)_dircard_hook]} )); then
  chpwd_functions+=(_dircard_hook)
fi
_dircard_hook`

const pwshHook = `$global:_dircardLastDir = ""
function global:prompt {
    [Console]::OutputEncoding = [System.Text.Encoding]::UTF8
    $cwd = (Get-Location).Path
    if ($cwd -ne $global:_dircardLastDir) {
        dircard show 2>$null | Out-Host
        $global:_dircardLastDir = $cwd
    }
    "PS $cwd> "
}`

var hookCmd = &cobra.Command{
	Use:   "hook [bash|zsh|pwsh]",
	Short: "Output the dircard shell hook",
	Long: `Output the dircard shell hook for the specified shell.

If no shell is specified, it will try to detect the current shell.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shell := ""

		if len(args) > 0 {
			shell = strings.ToLower(args[0])
		} else {
			shell = detectShell()
		}

		switch shell {
		case "bash":
			fmt.Print(bashHook)
		case "zsh":
			fmt.Print(zshHook)
		case "powershell", "pwsh":
			fmt.Print(pwshHook)
		default:
			fmt.Fprintf(os.Stderr, "error: unsupported shell: %s\n", shell)
			fmt.Fprintln(os.Stderr, "Supported shells: bash, zsh, pwsh")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(hookCmd)
	hookCmd.Hidden = true
}

func detectShell() string {
	// Git Bash / MSYS
	if os.Getenv("MSYSTEM") != "" {
		return "bash"
	}

	// PowerShell
	if os.Getenv("PSModulePath") != "" {
		return "pwsh"
	}

	// Unix-like shells
	if runtime.GOOS != "windows" {
		if shell := os.Getenv("SHELL"); shell != "" {
			name := strings.ToLower(filepath.Base(shell))

			switch name {
			case "bash", "zsh":
				return name
			}
		}
	}

	// Windows fallback
	if runtime.GOOS == "windows" {
		return "pwsh"
	}

	return ""
}
