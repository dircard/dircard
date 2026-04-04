package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dircard/dircard/internal/finder"
	"github.com/dircard/dircard/internal/marker"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a .dircard file in the current directory",
	Long:  `Creates a .dircard file in the current directory. Interactively choose a file type by default, or use --skip to skip the prompt.`,
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		path := resolveInitPath(cmd)

		action, err := marker.CreateOrUpdate(path, force)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}

		switch action {
		case marker.ActionCreated:
			fmt.Printf("Created %s\n", path)
		case marker.ActionAppended:
			fmt.Printf("Appended dircard markers to %s\n", path)
		case marker.ActionUpdated:
			fmt.Printf("Updated %s\n", path)
		}
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("path", "p", "", "Target directory path")
	initCmd.Flags().BoolP("skip", "k", false, "Skip interactive selection and create .dircard directly")
	initCmd.Flags().BoolP("force", "f", false, "Overwrite existing file (or append markers if target is README)")
}

// Path resolution helpers

func resolveTargetDir(cmd *cobra.Command) string {
	targetDir, _ := cmd.Flags().GetString("path")
	if targetDir == "" {
		return "."
	}

	if info, err := os.Stat(targetDir); err == nil {
		if info.IsDir() {
			return targetDir
		}
		fmt.Fprintln(os.Stderr, "error: --path must be a directory path")
		os.Exit(1)
	}

	if filepath.Ext(filepath.Base(targetDir)) != "" {
		fmt.Fprintln(os.Stderr, "error: --path must be a directory path")
		os.Exit(1)
	}

	return targetDir
}

func resolveInitPath(cmd *cobra.Command) string {
	targetDir := resolveTargetDir(cmd)
	skip, _ := cmd.Flags().GetBool("skip")

	if skip {
		return filepath.Join(targetDir, ".dircard")
	}

	if !isInteractiveTerminal() {
		fmt.Fprintln(os.Stderr, "error: interactive terminal required. Use --skip to skip selection.")
		os.Exit(1)
	}

	selected, err := selectFilePath(finder.Candidates)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}

	if targetDir == "." {
		return selected
	}
	return filepath.Join(targetDir, selected)
}

func isInteractiveTerminal() bool {
	fi, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (fi.Mode() & os.ModeCharDevice) != 0
}

func selectFilePath(candidates []finder.FileCandidate) (string, error) {
	names := make([]string, len(candidates))
	for i, c := range candidates {
		names[i] = c.Name
	}
	sel := promptui.Select{
		Label: "Select the .dircard file type to create",
		Items: names,
	}
	_, selected, err := sel.Run()
	return selected, err
}
