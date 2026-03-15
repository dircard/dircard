package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dircard/dircard/internal/shell"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

const bash = shell.HookComment + `
eval "$(dircard hook bash)"`

const zsh = shell.HookComment + `
eval "$(dircard hook zsh)"`

const pwsh = shell.HookComment + `
iex (dircard hook pwsh | Out-String)`

var installCmd = &cobra.Command{
	Use:   "install [bash|zsh|pwsh]",
	Short: "Install dircard shell integration",
	Long: `Install dircard shell integration for your environment.

Supported shells:
  bash
  zsh
  pwsh`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		shell, err := resolveInstallShell(args)
		if err != nil {
			if isPromptCanceled(err) {
				os.Exit(130)
			}
			fmt.Fprintln(os.Stderr, "error: failed to select shell:", err)
			os.Exit(1)
		}

		force, _ := cmd.Flags().GetBool("force")
		switch shell {
		case "bash":
			setupBash(force)
		case "zsh":
			setupZsh(force)
		case "pwsh":
			setupPowerShell(force)
		default:
			fmt.Fprintf(os.Stderr, "error: unsupported shell: %s\n", shell)
			fmt.Fprintln(os.Stderr, "Supported shells: bash, zsh, pwsh")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(installCmd)
	installCmd.Flags().BoolP("force", "f", false, "Overwrite existing hook if it exists")
}

func setupBash(force bool) {
	home, _ := os.UserHomeDir()
	rcPath := filepath.Join(home, ".bashrc")
	if added, err := shell.AppendHookToFile(rcPath, bash, force); err != nil {
		fmt.Fprintln(os.Stderr, "error: failed to update .bashrc:", err)
		os.Exit(1)
	} else if added {
		fmt.Printf("dircard hook added to %s\n", rcPath)
		fmt.Println("To apply changes: source ~/.bashrc")
	} else {
		fmt.Printf("dircard hook already exists in %s\n", rcPath)
		fmt.Println("Please use --force to update the existing hook if you want to apply changes.")
	}
}

func setupZsh(force bool) {
	home, _ := os.UserHomeDir()
	rcPath := filepath.Join(home, ".zshrc")
	if added, err := shell.AppendHookToFile(rcPath, zsh, force); err != nil {
		fmt.Fprintln(os.Stderr, "error: failed to update .zshrc:", err)
		os.Exit(1)
	} else if added {
		fmt.Printf("dircard hook added to %s\n", rcPath)
		fmt.Println("To apply changes: source ~/.zshrc")
	} else {
		fmt.Printf("dircard hook already exists in %s\n", rcPath)
		fmt.Println("Please use --force to update the existing hook if you want to apply changes.")
	}
}

func setupPowerShell(force bool) {
	profiles := shell.GetPowerShellProfiles()
	if len(profiles) == 0 {
		fmt.Fprintln(os.Stderr, "error: failed to get PowerShell profile path")
		os.Exit(1)
	}

	for _, p := range profiles {
		os.MkdirAll(filepath.Dir(p), 0755)
		if added, err := shell.AppendHookToFile(p, pwsh, force); err != nil {
			fmt.Fprintln(os.Stderr, "error: failed to update PowerShell profile:", err)
			os.Exit(1)
		} else if added {
			fmt.Printf("dircard hook added to %s\n", p)
			fmt.Println("To apply changes: restart your PowerShell session")
		} else {
			fmt.Printf("dircard hook already exists in %s\n", p)
			fmt.Println("Please use --force to update the existing hook if you want to apply changes.")
		}
	}
}

func selectShell() (string, error) {
	prompt := promptui.Select{
		Label: "Select shell",
		Items: []string{"zsh", "bash", "pwsh"},
		Size:  3,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return result, nil
}

func isPromptCanceled(err error) bool {
	return errors.Is(err, promptui.ErrInterrupt) || errors.Is(err, promptui.ErrEOF)
}

func resolveInstallShell(args []string) (string, error) {
	if len(args) > 0 {
		return strings.ToLower(args[0]), nil
	}

	selectedShell, err := selectShell()
	if err != nil {
		return "", err
	}

	return strings.ToLower(selectedShell), nil
}
