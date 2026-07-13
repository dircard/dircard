package renderer

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/mattn/go-runewidth"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/util"
	"golang.org/x/term"
)

const (
	ansiReset             = "\x1b[0m"
	ansiBold              = "\x1b[1m"
	ansiItalic            = "\x1b[3m"
	ansiYellow            = "\x1b[33m"
	ansiBlue              = "\x1b[34m"
	ansiUnderline         = "\x1b[4m"
	ansiDim               = "\x1b[2m"
	codeBlockBackground   = "\x1b[48;5;237m"
	codeBlockForeground   = "\x1b[38;5;252m"
	codeCommentForeground = "\x1b[38;5;114m"
	inlineCodeBackground  = "\x1b[48;5;238m"
	inlineCodeForeground  = "\x1b[38;5;167m"
)

// ANSIRenderer is a renderer that outputs ANSI escape sequences for terminal display
type ANSIRenderer struct {
	renderer.Config
	listDepth  int
	quoteDepth int
}

// RegisterFuncs registers renderer functions
func (r *ANSIRenderer) RegisterFuncs(reg renderer.NodeRendererFuncRegisterer) {
	// Document and structural elements
	reg.Register(ast.KindDocument, r.renderDocument)
	reg.Register(ast.KindHeading, r.renderHeading)
	reg.Register(ast.KindText, r.renderText)
	// Block elements
	reg.Register(ast.KindCodeBlock, r.renderCodeBlock)
	reg.Register(ast.KindFencedCodeBlock, r.renderFencedCodeBlock)
	reg.Register(ast.KindHTMLBlock, r.renderHTMLBlock)
	reg.Register(ast.KindParagraph, r.renderParagraph)
	reg.Register(ast.KindTextBlock, r.renderTextBlock)
	reg.Register(ast.KindList, r.renderList)
	reg.Register(ast.KindListItem, r.renderListItem)
	// Inline elements
	reg.Register(ast.KindCodeSpan, r.renderCodeSpan)
	reg.Register(ast.KindString, r.renderString)
	reg.Register(ast.KindRawHTML, r.renderRawHTML)
	reg.Register(ast.KindEmphasis, r.renderEmphasis)
	reg.Register(ast.KindLink, r.renderLink)
	reg.Register(ast.KindAutoLink, r.renderAutoLink)
	reg.Register(ast.KindImage, r.renderImage)
	// Special blocks
	reg.Register(ast.KindThematicBreak, r.renderThematicBreak)
	reg.Register(ast.KindBlockquote, r.renderBlockquote)
}

func (r *ANSIRenderer) renderDocument(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *ANSIRenderer) renderHeading(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		heading := node.(*ast.Heading)
		w.WriteString(ansiBold)
		w.WriteString(ansiBlue)
		w.WriteString(strings.Repeat("#", heading.Level))
		w.WriteString(" ")
	} else {
		w.WriteString(ansiReset)
		w.WriteString("\n\n")
	}
	return ast.WalkContinue, nil
}

func (r *ANSIRenderer) renderText(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		text := node.(*ast.Text)
		w.Write(text.Segment.Value(source))
		if text.HardLineBreak() || text.SoftLineBreak() {
			w.WriteByte('\n')
		}
	}
	return ast.WalkContinue, nil
}

func (r *ANSIRenderer) renderCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.renderCodeLines(w, source, node)
	} else {
		w.WriteByte('\n')
	}
	return ast.WalkContinue, nil
}

func (r *ANSIRenderer) renderFencedCodeBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.renderCodeLines(w, source, node)
	} else {
		w.WriteByte('\n')
	}
	return ast.WalkContinue, nil
}

func (r *ANSIRenderer) renderParagraph(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.writeQuotePrefix(w)
		return ast.WalkContinue, nil
	}
	if r.isTightListParagraph(node) {
		w.WriteByte('\n')
		return ast.WalkContinue, nil
	}
	w.WriteString("\n\n")
	return ast.WalkContinue, nil
}

func (r *ANSIRenderer) renderTextBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *ANSIRenderer) renderList(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.listDepth++
	} else {
		r.listDepth--
		if r.listDepth == 0 {
			w.WriteByte('\n')
		}
	}
	return ast.WalkContinue, nil
}

func (r *ANSIRenderer) renderListItem(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString(strings.Repeat("  ", max(r.listDepth-1, 0)))
		listItem := node.(*ast.ListItem)
		parent := listItem.Parent()
		if parent != nil {
			if list, ok := parent.(*ast.List); ok {
				if list.IsOrdered() {
					w.WriteString(fmt.Sprintf("%d. ", r.listItemNumber(listItem, list)))
				} else {
					w.WriteString("・ ")
				}
			}
		}
		w.WriteString(ansiYellow)
	} else {
		w.WriteString(ansiReset)
		if _, ok := node.FirstChild().(*ast.TextBlock); ok {
			w.WriteByte('\n')
		}
	}
	return ast.WalkContinue, nil
}

func (r *ANSIRenderer) renderCodeSpan(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString(inlineCodeBackground)
		w.WriteString(inlineCodeForeground)
		w.WriteString(" ") // Add a space
		for c := node.FirstChild(); c != nil; c = c.NextSibling() {
			text, ok := c.(*ast.Text)
			if !ok {
				continue
			}
			value := text.Segment.Value(source)
			value = bytes.ReplaceAll(value, []byte("\n"), []byte(" "))
			w.Write(value)
		}
		w.WriteString(" ") // Add a space
		w.WriteString(ansiReset)
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

func (r *ANSIRenderer) renderString(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		str := node.(*ast.String)
		w.Write(str.Value)
	}
	return ast.WalkContinue, nil
}

func (r *ANSIRenderer) renderRawHTML(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	return ast.WalkContinue, nil
}

func (r *ANSIRenderer) renderHTMLBlock(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

func (r *ANSIRenderer) renderEmphasis(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	emphasis := node.(*ast.Emphasis)
	if entering {
		if emphasis.Level == 1 {
			w.WriteString(ansiItalic)
		} else {
			w.WriteString(ansiBold)
		}
	} else {
		w.WriteString(ansiReset)
	}
	return ast.WalkContinue, nil
}

func (r *ANSIRenderer) renderLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString(ansiUnderline)
		w.WriteString(ansiBlue)
	} else {
		w.WriteString(ansiReset)
		link := node.(*ast.Link)
		w.WriteString(" (")
		w.Write(link.Destination)
		w.WriteString(")")
	}
	return ast.WalkContinue, nil
}

func (r *ANSIRenderer) renderAutoLink(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		autoLink := node.(*ast.AutoLink)
		w.WriteString(ansiUnderline)
		w.WriteString(ansiBlue)
		w.Write(autoLink.Label(source))
		w.WriteString(ansiReset)
		return ast.WalkSkipChildren, nil
	}
	return ast.WalkContinue, nil
}

func (r *ANSIRenderer) renderImage(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString(ansiDim)
		w.WriteString("[Image: ")
	} else {
		image := node.(*ast.Image)
		w.WriteString("] (")
		w.Write(image.Destination)
		w.WriteString(")")
		w.WriteString(ansiReset)
	}
	return ast.WalkContinue, nil
}

func (r *ANSIRenderer) renderThematicBreak(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		w.WriteString(ansiDim)
		w.WriteString("────────")
		w.WriteString(ansiReset)
		w.WriteByte('\n')
	}
	return ast.WalkContinue, nil
}

func (r *ANSIRenderer) renderBlockquote(w util.BufWriter, source []byte, node ast.Node, entering bool) (ast.WalkStatus, error) {
	if entering {
		r.quoteDepth++
	} else {
		r.quoteDepth--
		if r.quoteDepth == 0 {
			w.WriteByte('\n')
		}
	}
	return ast.WalkContinue, nil
}

var getTerminalWidthFunc = getTerminalWidth

func getTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 80
	}
	return width
}

func (r *ANSIRenderer) renderCodeLines(w util.BufWriter, source []byte, node ast.Node) {
	lines := node.Lines()
	terminalWidth := getTerminalWidthFunc()
	inBlockComment := false
	for i := 0; i < lines.Len(); i++ {
		line := lines.At(i)
		value := line.Value(source)
		value = trimLineBreak(value)
		w.WriteString(codeBlockBackground)
		w.WriteString(codeBlockForeground)
		w.WriteString("  ")
		r.renderCodeLine(w, value, &inBlockComment)
		currentWidth := 2 + runewidth.StringWidth(string(value))
		padding := terminalWidth - currentWidth
		if padding > 0 {
			w.WriteString(strings.Repeat(" ", padding))
		}
		w.WriteString(ansiReset)
		w.WriteByte('\n')
	}
}

func (r *ANSIRenderer) renderCodeLine(w util.BufWriter, value []byte, inBlockComment *bool) {
	for len(value) > 0 {
		if *inBlockComment {
			end := bytes.Index(value, []byte("*/"))
			w.WriteString(codeCommentForeground)
			if end == -1 {
				w.Write(value)
				return
			}
			end += len("*/")
			w.Write(value[:end])
			value = value[end:]
			*inBlockComment = false
			w.WriteString(codeBlockForeground)
			continue
		}

		commentStart, blockComment := findCodeCommentStart(value)
		if commentStart == -1 {
			w.Write(value)
			return
		}

		w.Write(value[:commentStart])
		value = value[commentStart:]
		w.WriteString(codeCommentForeground)
		if !blockComment {
			w.Write(value)
			return
		}

		end := bytes.Index(value, []byte("*/"))
		if end == -1 {
			w.Write(value)
			*inBlockComment = true
			return
		}
		end += len("*/")
		w.Write(value[:end])
		value = value[end:]
		w.WriteString(codeBlockForeground)
	}
}

func findCodeCommentStart(value []byte) (int, bool) {
	commentStart := -1
	blockComment := false
	for _, marker := range []struct {
		value        []byte
		blockComment bool
	}{
		{[]byte("//"), false},
		{[]byte("#"), false},
		{[]byte("<!--"), false},
		{[]byte("/*"), true},
	} {
		index := bytes.Index(value, marker.value)
		if index == -1 {
			continue
		}
		if commentStart == -1 || index < commentStart {
			commentStart = index
			blockComment = marker.blockComment
		}
	}
	return commentStart, blockComment
}

func trimLineBreak(value []byte) []byte {
	value = bytes.TrimSuffix(value, []byte("\n"))
	return bytes.TrimSuffix(value, []byte("\r"))
}

func (r *ANSIRenderer) writeQuotePrefix(w util.BufWriter) {
	for i := 0; i < r.quoteDepth; i++ {
		w.WriteString(ansiDim)
		w.WriteString("│ ")
		w.WriteString(ansiReset)
	}
}

func (r *ANSIRenderer) isTightListParagraph(node ast.Node) bool {
	parent := node.Parent()
	if parent == nil {
		return false
	}
	listItem, ok := parent.(*ast.ListItem)
	if !ok {
		return false
	}
	if listItem.Parent() == nil {
		return false
	}
	list, ok := listItem.Parent().(*ast.List)
	return ok && list.IsTight
}

func (r *ANSIRenderer) listItemNumber(item *ast.ListItem, list *ast.List) int {
	number := list.Start
	for sibling := item.PreviousSibling(); sibling != nil; sibling = sibling.PreviousSibling() {
		if _, ok := sibling.(*ast.ListItem); ok {
			number++
		}
	}
	return number
}

// ParseMarkdown converts Markdown to ANSI escape sequences for terminal display
func ParseMarkdown(content string) string {
	md := goldmark.New(
		goldmark.WithRenderer(renderer.NewRenderer(renderer.WithNodeRenderers(util.Prioritized(&ANSIRenderer{}, 100)))),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(content), &buf); err != nil {
		return content
	}

	return trimTrailingNewlines(buf.String())
}

func trimTrailingNewlines(result string) string {
	if result == "" {
		return result
	}
	return strings.TrimRight(result, "\n") + "\n"
}
