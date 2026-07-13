/*
Copyright © 2026 yhotta240 <yhotta240@gmail.com>
*/
package cmd

import (
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
	"golang.org/x/mod/module"
	"golang.org/x/mod/semver"
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
	rootCmd.Version = buildVersion()
}

func buildVersion() string {
	if version != "dev" {
		return version
	}

	info, ok := debug.ReadBuildInfo()
	if ok && isTaggedVersion(info.Main.Version) {
		return info.Main.Version
	}

	return version
}

func isTaggedVersion(v string) bool {
	return semver.IsValid(v) && !module.IsPseudoVersion(v)
}
