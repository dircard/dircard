package cmd

import (
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/dircard/dircard/internal/config"
	"github.com/dircard/dircard/internal/finder"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage dircard configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if !isInteractiveTerminal() {
			cfg := config.Load()
			fmt.Printf("Candidate order: %s\n", strings.Join(cfg.CandidateOrder, ", "))
			fmt.Printf("File size limit: %d KB\n", cfg.FileSizeKB)
			fmt.Printf("Line start: %d\n", *cfg.LineStart)
			fmt.Printf("Line count: %d\n", *cfg.LineCount)
			fmt.Printf("Search depth: %d\n", *cfg.Depth)
			return
		}
		runInteractiveConfig()
	},
}

var orderCmd = &cobra.Command{
	Use:   "order [file-name]",
	Short: "Set highest-priority dircard file type",
	Long:  "Set highest-priority dircard file type. If file-name is omitted, choose interactively.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		selected := ""
		if len(args) == 1 {
			selected = args[0]
		} else {
			if !isInteractiveTerminal() {
				fmt.Fprintln(os.Stderr, "error: file-name is required in non-interactive terminal")
				os.Exit(1)
			}
			cfg := config.Load()
			candidates := finder.ReorderCandidates(cfg.CandidateOrder)
			picked, err := selectFilePath(candidates)
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				os.Exit(1)
			}
			selected = picked
		}
		editCandidateOrder(selected)
	},
}

var sizeCmd = &cobra.Command{
	Use:   "size [size]",
	Short: "Set file size limit for dircard",
	Long:  "Set file size limit for dircard. Files larger than this limit will be ignored. If size is omitted, choose interactively.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			size, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, "error: invalid file size limit")
				os.Exit(1)
			}
			saveConfig(config.WithFileSizeKB(size))
			fmt.Printf("File size limit updated: %d KB\n", size)
		} else {
			if !isInteractiveTerminal() {
				fmt.Fprintln(os.Stderr, "error: size is required in non-interactive terminal")
				os.Exit(1)
			}
			cfg := config.Load()
			editNumberSetting("file size limit (KB)", cfg.FileSizeKB, config.WithFileSizeKB)
		}
	},
}

var lineStartCmd = &cobra.Command{
	Use:   "linestart [start]",
	Short: "Set line start for dircard",
	Long:  "Set line start for dircard. Lines before this will be ignored. If start is omitted, choose interactively.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			start, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, "error: invalid line start")
				os.Exit(1)
			}
			saveConfig(config.WithLineStart(start))
			fmt.Printf("Line start updated: %d\n", start)
		} else {
			if !isInteractiveTerminal() {
				fmt.Fprintln(os.Stderr, "error: start is required in non-interactive terminal")
				os.Exit(1)
			}
			cfg := config.Load()
			editNumberSetting("line start", *cfg.LineStart, config.WithLineStart)
		}
	},
}

var lineCountCmd = &cobra.Command{
	Use:   "linecount [count]",
	Short: "Set line count for dircard",
	Long:  "Set line count for dircard. Lines after this will be ignored. If count is omitted, choose interactively.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			count, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, "error: invalid line count")
				os.Exit(1)
			}
			saveConfig(config.WithLineCount(count))
			fmt.Printf("Line count updated: %d\n", count)
		} else {
			if !isInteractiveTerminal() {
				fmt.Fprintln(os.Stderr, "error: count is required in non-interactive terminal")
				os.Exit(1)
			}
			cfg := config.Load()
			editNumberSetting("line count", *cfg.LineCount, config.WithLineCount)
		}
	},
}

var depthCmd = &cobra.Command{
	Use:   "depth [depth]",
	Short: "Set search depth for dircard",
	Long:  "Set search depth for dircard. Files will be searched up to this depth. If depth is omitted, choose interactively.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 1 {
			depth, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Fprintln(os.Stderr, "error: invalid depth")
				os.Exit(1)
			}
			saveConfig(config.WithDepth(depth))
			fmt.Printf("Search depth updated: %d\n", depth)
		} else {
			if !isInteractiveTerminal() {
				fmt.Fprintln(os.Stderr, "error: depth is required in non-interactive terminal")
				os.Exit(1)
			}
			cfg := config.Load()
			editNumberSetting("search depth", *cfg.Depth, config.WithDepth)
		}
	},
}

func saveConfig(opts ...config.ConfigOption) {
	if err := config.ApplyOptions(opts...); err != nil {
		fmt.Fprintln(os.Stderr, "error: failed to save config:", err)
		os.Exit(1)
	}
}

func promptNumber(prompt string, defaultValue int) (int, error) {
	promptStr := promptui.Prompt{
		Label:   prompt,
		Default: fmt.Sprintf("%d", defaultValue),
		Validate: func(input string) error {
			if input == "" {
				return nil
			}
			_, err := strconv.Atoi(input)
			if err != nil {
				return fmt.Errorf("must be a number")
			}
			return nil
		},
	}
	result, err := promptStr.Run()
	if err != nil {
		return 0, err
	}
	if result == "" {
		return defaultValue, nil
	}
	return strconv.Atoi(result)
}

func runInteractiveConfig() {
	for {
		cfg := config.Load()

		items := []string{
			fmt.Sprintf("Candidate order: %s", strings.Join(cfg.CandidateOrder, ", ")),
			fmt.Sprintf("File size limit: %d KB", cfg.FileSizeKB),
			fmt.Sprintf("Line start: %d", *cfg.LineStart),
			fmt.Sprintf("Line count: %d", *cfg.LineCount),
			fmt.Sprintf("Search depth: %d", *cfg.Depth),
			"Exit",
		}

		sel := promptui.Select{
			Label: "Select configuration to edit",
			Items: items,
		}

		index, _, err := sel.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error:", err)
			return
		}

		switch index {
		case 0: // Candidate order
			selected, err := selectFilePath(finder.ReorderCandidates(cfg.CandidateOrder))
			if err != nil {
				fmt.Fprintln(os.Stderr, "error:", err)
				continue
			}
			editCandidateOrder(selected)

		case 1: // File size limit
			editNumberSetting("file size limit (KB)", cfg.FileSizeKB, config.WithFileSizeKB)

		case 2: // Line start
			editNumberSetting("line start", *cfg.LineStart, config.WithLineStart)

		case 3: // Line count
			editNumberSetting("line count", *cfg.LineCount, config.WithLineCount)

		case 4: // Search depth
			editNumberSetting("search depth", *cfg.Depth, config.WithDepth)

		case 5: // Exit
			return
		}
	}
}

func editNumberSetting(label string, currentValue int, optionFunc func(int) config.ConfigOption) {
	prompt := fmt.Sprintf("Enter %s [current: %d]", label, currentValue)
	value, err := promptNumber(prompt, currentValue)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		return
	}
	saveConfig(optionFunc(value))
	fmt.Printf("%s updated: %d\n", label, value)
}

func editCandidateOrder(selected string) {
	candidateNames := finder.Names(finder.Candidates)
	if !slices.Contains(candidateNames, selected) {
		fmt.Fprintf(os.Stderr, "error: invalid file name. Must be one of %s\n", strings.Join(candidateNames, ", "))
		os.Exit(1)
	}

	cfg := config.Load()
	ordered := finder.ReorderCandidates(cfg.CandidateOrder)
	newOrder := []string{selected}
	for _, c := range ordered {
		if c.Name != selected {
			newOrder = append(newOrder, c.Name)
		}
	}
	saveConfig(config.WithCandidateOrder(newOrder))
	fmt.Printf("Priority updated: %s\n", selected)
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(orderCmd)
	configCmd.AddCommand(sizeCmd)
	configCmd.AddCommand(lineStartCmd)
	configCmd.AddCommand(lineCountCmd)
	configCmd.AddCommand(depthCmd)
}
