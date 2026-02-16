package shared

import (
	"bytes"
	"strings"
	"unicode"

	"github.com/tsawler/prose/v3"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// ExtractProse parses markdown and returns only prose paragraph text,
// stripping headings, code blocks, lists, block quotes, tables, and
// HTML blocks. Each returned string is one markdown paragraph.
func ExtractProse(markdown string) []string {
	source := []byte(markdown)
	reader := text.NewReader(source)
	parser := goldmark.DefaultParser()
	doc := parser.Parse(reader)

	var paragraphs []string
	for node := doc.FirstChild(); node != nil; node = node.NextSibling() {
		if node.Kind() != ast.KindParagraph {
			continue
		}
		var buf bytes.Buffer
		extractText(node, source, &buf)
		if trimmed := strings.TrimSpace(buf.String()); trimmed != "" {
			paragraphs = append(paragraphs, trimmed)
		}
	}
	return paragraphs
}

// extractText recursively collects text content from an AST node.
func extractText(node ast.Node, source []byte, buf *bytes.Buffer) {
	if node.Type() == ast.TypeInline {
		if t, ok := node.(*ast.Text); ok {
			buf.Write(t.Segment.Value(source))
			if t.SoftLineBreak() || t.HardLineBreak() {
				buf.WriteByte(' ')
			}
		}
	}
	for child := node.FirstChild(); child != nil; child = child.NextSibling() {
		extractText(child, source, buf)
	}
}

// ExtractProseText parses markdown and returns all prose paragraphs
// joined into a single string, separated by spaces.
func ExtractProseText(markdown string) string {
	return strings.Join(ExtractProse(markdown), " ")
}

// SplitParagraphs extracts prose paragraphs from markdown text.
func SplitParagraphs(markdown string) []string {
	return ExtractProse(markdown)
}

// SplitSentences splits text into sentences using prose's sentence tokenizer.
// Returns an empty slice if the text is empty or cannot be parsed.
func SplitSentences(text string) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return []string{}
	}
	doc, err := prose.NewDocument(text, prose.WithExtraction(false), prose.WithTagging(false))
	if err != nil {
		return []string{}
	}
	sents := doc.Sentences()
	result := make([]string, 0, len(sents))
	for _, s := range sents {
		if trimmed := strings.TrimSpace(s.Text); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// ListWords splits text into words on whitespace.
func ListWords(text string) []string {
	return strings.Fields(text)
}

// CountWords returns the number of whitespace-separated words in text.
func CountWords(text string) int {
	return len(strings.Fields(text))
}

// StripPunctuation removes leading and trailing non-letter characters from a word.
func StripPunctuation(word string) string {
	runes := []rune(word)
	start, end := 0, len(runes)
	for start < end && !unicode.IsLetter(runes[start]) {
		start++
	}
	for end > start && !unicode.IsLetter(runes[end-1]) {
		end--
	}
	if start >= end {
		return ""
	}
	return string(runes[start:end])
}

// Tokenize returns POS-tagged tokens using prose's NLP pipeline.
// Returns an empty slice if the text is empty or cannot be parsed.
func Tokenize(text string) []Token {
	text = strings.TrimSpace(text)
	if text == "" {
		return []Token{}
	}
	doc, err := prose.NewDocument(text)
	if err != nil {
		return []Token{}
	}
	toks := doc.Tokens()
	result := make([]Token, len(toks))
	for i, t := range toks {
		result[i] = Token{
			Text:  t.Text,
			Tag:   t.Tag,
			Label: t.Label,
		}
	}
	return result
}
