package renderer

import (
	"strings"
	"testing"

	"github.com/mattn/go-runewidth"
)

func TestRenderMarkdownKeepsBlankLineBetweenParagraphs(t *testing.T) {
	got := ParseMarkdown("first\n\nsecond")
	want := "first\n\nsecond"

	if got != want {
		t.Fatalf("ParseMarkdown() = %q, want %q", got, want)
	}
}

func TestRenderMarkdownRendersInlineCodeWithDedicatedStyle(t *testing.T) {
	got := ParseMarkdown("Use `value` here")
	want := "Use " + inlineCodeBackground + inlineCodeForeground + " value " + ansiReset + " here"

	if got != want {
		t.Fatalf("ParseMarkdown() = %q, want %q", got, want)
	}
}

func TestRenderMarkdownRestoresListStyleAfterInlineCode(t *testing.T) {
	got := ParseMarkdown("- before `value` after")
	want := "・ " + ansiYellow + "before " + inlineCodeBackground + inlineCodeForeground +
		" value " + ansiReset + ansiYellow + " after" + ansiReset

	if got != want {
		t.Fatalf("ParseMarkdown() = %q, want %q", got, want)
	}
}

func TestRenderMarkdownResetsCodeBlockStyleBeforeBlankLine(t *testing.T) {
	originalFunc := getTerminalWidthFunc
	getTerminalWidthFunc = func() int { return 80 }
	defer func() { getTerminalWidthFunc = originalFunc }()

	got := ParseMarkdown("```go\nfmt.Println(\"x\")\n```\n\nnext")
	want := codeBlockLine("fmt.Println(\"x\")") + "\n\nnext"

	if got != want {
		t.Fatalf("ParseMarkdown() = %q, want %q", got, want)
	}
}

func TestRenderMarkdownPadsCodeBlockBackgroundToTerminalWidth(t *testing.T) {
	originalFunc := getTerminalWidthFunc
	getTerminalWidthFunc = func() int { return 80 }
	defer func() { getTerminalWidthFunc = originalFunc }()

	got := ParseMarkdown("```go\na\nlong\n```")
	want := codeBlockLine("a") + "\n" + codeBlockLine("long")

	if got != want {
		t.Fatalf("ParseMarkdown() = %q, want %q", got, want)
	}
}

func TestRenderMarkdownPadsWideCodeBlockCommentsToDisplayWidth(t *testing.T) {
	originalFunc := getTerminalWidthFunc
	getTerminalWidthFunc = func() int { return 20 }
	defer func() { getTerminalWidthFunc = originalFunc }()

	got := ParseMarkdown("```sh\nlong #コメント\n```")
	want := renderedCodeBlockLineWithWidth(
		"long "+codeCommentForeground+"#コメント",
		"long #コメント",
		20,
	)

	if got != want {
		t.Fatalf("ParseMarkdown() = %q, want %q", got, want)
	}
}

func TestRenderMarkdownRendersCodeBlockCommentsInGreen(t *testing.T) {
	originalFunc := getTerminalWidthFunc
	getTerminalWidthFunc = func() int { return 80 }
	defer func() { getTerminalWidthFunc = originalFunc }()

	got := ParseMarkdown("```go\nvalue := 1 // keep comment\n/* block */ value\n```")
	want := renderedCodeBlockLine("value := 1 "+codeCommentForeground+"// keep comment", "value := 1 // keep comment") +
		"\n" + renderedCodeBlockLine(codeCommentForeground+"/* block */"+codeBlockForeground+" value", "/* block */ value")

	if got != want {
		t.Fatalf("ParseMarkdown() = %q, want %q", got, want)
	}
}

func TestRenderMarkdownRemovesCarriageReturnsFromCodeBlockLines(t *testing.T) {
	originalFunc := getTerminalWidthFunc
	getTerminalWidthFunc = func() int { return 80 }
	defer func() { getTerminalWidthFunc = originalFunc }()

	got := ParseMarkdown("```go\r\nfmt.Println(\"x\")\r\n```\r\n")
	want := codeBlockLine("fmt.Println(\"x\")")

	if got != want {
		t.Fatalf("ParseMarkdown() = %q, want %q", got, want)
	}
}

func TestRenderMarkdownIncrementsOrderedListNumbers(t *testing.T) {
	got := ParseMarkdown("1. one\n2. two")
	want := "1. " + ansiYellow + "one" + ansiReset + "\n" +
		"2. " + ansiYellow + "two" + ansiReset

	if got != want {
		t.Fatalf("ParseMarkdown() = %q, want %q", got, want)
	}
}

func codeBlockLine(value string) string {
	return renderedCodeBlockLine(value, value)
}

func renderedCodeBlockLine(rendered, visible string) string {
	return renderedCodeBlockLineWithWidth(rendered, visible, 80)
}

func renderedCodeBlockLineWithWidth(rendered, visible string, terminalWidth int) string {
	currentWidth := 2 + runewidth.StringWidth(visible)
	padding := terminalWidth - currentWidth
	if padding < 0 {
		padding = 0
	}
	return codeBlockBackground + codeBlockForeground + "  " + rendered +
		strings.Repeat(" ", padding) + ansiReset
}
