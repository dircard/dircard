package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dircard/dircard/internal/fileio"
	"github.com/dircard/dircard/internal/finder"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a .dircard file in the current directory",
	Long:  `Creates a new notes file in the current directory. Interactively choose a file type by default, or use --skip to create .dircard. Returns an error if the file already exists.`,
	Run: func(cmd *cobra.Command, args []string) {
		path := resolveInitPath(cmd)

		if err := fileio.CreateFile(path); err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			os.Exit(1)
		}

		fmt.Printf("Created %s\n", path)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringP("path", "p", "", "Target directory path")
	initCmd.Flags().BoolP("skip", "k", false, "Skip interactive selection and create .dircard directly")
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
