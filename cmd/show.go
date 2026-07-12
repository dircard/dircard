/*
Copyright © 2026 yhotta240 <yhotta240@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dircard/dircard/internal/config"
	"github.com/dircard/dircard/internal/fileio"
	"github.com/dircard/dircard/internal/finder"
	"github.com/dircard/dircard/internal/renderer"
	"github.com/spf13/cobra"
)

type Dircard struct {
	ID        int64  `json:"id"`
	CreatedAt string `json:"created_at"`
	Path      string `json:"path"`
	Depth     int    `json:"depth"`
	Content   string `json:"content"`
	StartLine int    `json:"start_line"`
	LineCount int    `json:"line_count"`
}

var fileSize int
var lineStart int
var lineCount int
var searchDepth int
var section string

var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the contents of the .dircard file for the current directory",
	Long:  `Show the contents of the .dircard file for the current directory. By default, it shows the last 10 lines, but you can specify a different number of lines or a starting line number using the --lines and --start flags.`,
	Run: func(cmd *cobra.Command, args []string) {
		quiet, _ := cmd.Flags().GetBool("quiet")
		if quiet {
			return
		}

		cwd, err := os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: failed to get current directory:", err)
			return
		}

		cfg := config.Load()
		sortedCandidates := finder.ReorderCandidates(cfg.CandidateOrder)
		dircardFile, err := finder.FindFilePath(cwd, searchDepth, sortedCandidates)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: .dircard file not found in current or parent directories")
			return
		}

		path, _ := cmd.Flags().GetBool("path")
		if path {
			fmt.Printf("%s\n", dircardFile)
			return
		}

		info, err := os.Stat(dircardFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: failed to get file info:", err)
			return
		}
		if info.Size() > int64(fileSize*1024) {
			fmt.Fprintf(os.Stderr, "error: file is too large (max %dKB): %s\n", fileSize, dircardFile)
			return
		}

		full, _ := cmd.Flags().GetBool("full")
		lines, err := fileio.ReadFileLines(dircardFile, full, lineStart, lineCount)
		if err != nil {
			fmt.Fprintln(os.Stderr, "error: failed to read file:", err)
			return
		}

		jsonOutput, _ := cmd.Flags().GetBool("json")
		if jsonOutput {
			b, err := createJSONOutput(dircardFile, lines)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error: failed to create JSON output:", err)
				return
			}
			fmt.Println(b)
			return
		}

		result := renderer.ParseMarkdown(strings.Join(lines, "\n"))
		fmt.Println(result)

		if !full && len(lines) == lineCount {
			fmt.Println("... (run `dircard show --full` to see more)")
		}
	},
}

func init() {
	cfg := config.Load()
	rootCmd.AddCommand(showCmd)
	showCmd.Flags().BoolP("quiet", "q", false, "Suppress all output (quiet mode)")
	showCmd.Flags().BoolP("path", "p", false, "Show the path of the .dircard file")
	showCmd.Flags().BoolP("full", "f", false, "Show the full contents of the .dircard file")
	showCmd.Flags().IntVarP(&searchDepth, "depth", "d", *cfg.Depth, "Limit the depth of parent directory search (0=only current directory)")
	showCmd.Flags().StringVarP(&section, "section", "e", "", "Show only the specified section")
	showCmd.Flags().IntVarP(&fileSize, "size", "z", cfg.FileSizeKB, "Maximum file size to show in KB")
	showCmd.Flags().IntVarP(&lineStart, "start", "s", *cfg.LineStart, "Line number to start showing from (0=beginning)")
	showCmd.Flags().IntVarP(&lineCount, "lines", "n", *cfg.LineCount, "Number of lines to show (0=all)")
	showCmd.Flags().BoolP("json", "j", false, "Output in JSON format")
}

func createJSONOutput(dircardFile string, lines []string) (string, error) {
	today := time.Now().Format("2006-01-02 15:04:05")
	content := strings.Join(lines, "\n")
	data := Dircard{
		ID:        time.Now().UnixNano(),
		CreatedAt: today,
		Path:      dircardFile,
		Depth:     searchDepth,
		Content:   content,
		StartLine: lineStart,
		LineCount: lineCount,
	}

	b, err := json.Marshal(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: failed to marshal JSON:", err)
		return "", err
	}

	return string(b), nil
}
