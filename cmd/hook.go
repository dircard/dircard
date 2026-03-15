package cmd

import (
	"fmt"
	"os"
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
esac
`

const zshHook = `_dircard_hook() {
  command dircard show 2>/dev/null
}
if ! (( ${chpwd_functions[(I)_dircard_hook]} )); then
  chpwd_functions+=(_dircard_hook)
fi
_dircard_hook
`

const pwshHook = `$global:_dircardLastDir = ""
function global:prompt {
    [Console]::OutputEncoding = [System.Text.Encoding]::UTF8
    $cwd = (Get-Location).Path
    if ($cwd -ne $global:_dircardLastDir) {
        dircard show 2>$null | Out-Host
        $global:_dircardLastDir = $cwd
    }
    "PS $cwd> "
}
`

var hookCmd = &cobra.Command{
	Use:   "hook [bash|zsh|pwsh]",
	Short: "Output the dircard shell hook",
	Long:  `Output the dircard shell hook for the specified shell. This command is intended to be called from your shell configuration files (e.g. .bashrc, .zshrc, profile.ps1) to enable dircard integration.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shell := strings.ToLower(args[0])

		switch shell {
		case "bash":
			fmt.Print(bashHook)
		case "zsh":
			fmt.Print(zshHook)
		case "pwsh":
			fmt.Print(pwshHook)
		default:
			fmt.Fprintf(os.Stderr, "Unsupported shell: %s\n", shell)
			fmt.Fprintln(os.Stderr, "Supported shells: bash, zsh, pwsh")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(hookCmd)
	// Hidden because this command is called from shell configuration files (e.g. .bashrc, .zshrc, profile.ps1),
	// not invoked directly by users.
	hookCmd.Hidden = true
}
