package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dircard/dircard/internal/shell"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall [bash|zsh|pwsh]",
	Short: "Uninstall dircard shell integration",
	Long: `Uninstall dircard shell integration by removing the hook from your shell configuration files.

Supported shells:
  bash
  zsh
  pwsh

If no shell is specified, the hook will be removed from all supported shell profiles.`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")

		if len(args) == 0 {
			if !force {
				confirmed, err := confirmUninstall("all shells")
				if err != nil {
					if isPromptCanceled(err) {
						os.Exit(130)
					}
					fmt.Fprintln(os.Stderr, "error:", err)
					os.Exit(1)
				}
				if !confirmed {
					fmt.Println("Aborted.")
					return
				}
			}
			teardownBash()
			teardownZsh()
			teardownPwsh()
			return
		}

		shell := strings.ToLower(args[0])
		if !force {
			confirmed, err := confirmUninstall(shell)
			if err != nil {
				if isPromptCanceled(err) {
					os.Exit(130)
				}
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
			if !confirmed {
				fmt.Println("Aborted.")
				return
			}
		}

		switch shell {
		case "bash":
			teardownBash()
		case "zsh":
			teardownZsh()
		case "pwsh":
			teardownPwsh()
		default:
			fmt.Fprintf(os.Stderr, "error: unsupported shell: %s\n", shell)
			fmt.Fprintln(os.Stderr, "Supported shells: bash, zsh, pwsh")
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
	uninstallCmd.Flags().BoolP("force", "f", false, "Uninstall without confirmation")
}

func confirmUninstall(target string) (bool, error) {
	prompt := promptui.Select{
		Label: fmt.Sprintf("Remove dircard hook from %s?", target),
		Items: []string{"Yes", "No"},
		Size:  2,
	}
	_, result, err := prompt.Run()
	if err != nil {
		return false, err
	}
	return result == "Yes", nil
}

func teardownBash() {
	home, _ := os.UserHomeDir()
	rcPath := filepath.Join(home, ".bashrc")
	if removed, err := shell.RemoveHookFromFile(rcPath); err != nil {
		fmt.Fprintln(os.Stderr, "error: failed to update .bashrc:", err)
		os.Exit(1)
	} else if removed {
		fmt.Printf("dircard hook removed from %s\n", rcPath)
		fmt.Println("To apply changes: source ~/.bashrc")
	} else {
		fmt.Printf("dircard hook not found in %s\n", rcPath)
	}
}

func teardownZsh() {
	home, _ := os.UserHomeDir()
	rcPath := filepath.Join(home, ".zshrc")
	if removed, err := shell.RemoveHookFromFile(rcPath); err != nil {
		fmt.Fprintln(os.Stderr, "error: failed to update .zshrc:", err)
		os.Exit(1)
	} else if removed {
		fmt.Printf("dircard hook removed from %s\n", rcPath)
		fmt.Println("To apply changes: source ~/.zshrc")
	} else {
		fmt.Printf("dircard hook not found in %s\n", rcPath)
	}
}

func teardownPwsh() {
	profiles := shell.GetPowerShellProfiles()
	if len(profiles) == 0 {
		fmt.Fprintln(os.Stderr, "error: failed to get PowerShell profile path")
		os.Exit(1)
	}

	for _, p := range profiles {
		if removed, err := shell.RemoveHookFromFile(p); err != nil {
			fmt.Fprintln(os.Stderr, "error: failed to update PowerShell profile:", err)
			os.Exit(1)
		} else if removed {
			fmt.Printf("dircard hook removed from %s\n", p)
			fmt.Println("To apply changes: restart your PowerShell session")
		} else {
			fmt.Printf("dircard hook not found in %s\n", p)
		}
	}
}
