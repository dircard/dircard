/*
Copyright © 2026 yhotta240 <yhotta240@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var version = "dev"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "dircard",
	Short: "Display notes for the current directory",
	Long: `Displays notes and metadata associated with the current directory.

If a .dircard file exists, its contents will be automatically shown upon entering the directory. This helps you maintain context for your projects and tasks.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Version = version
}
