package parser

import (
	"bytes"
	"errors"

	"github.com/gomarkdown/markdown/ast"
)

var (
	frontmatterDelim = []byte("---")
)

// TODO: doc should be a structure with metadata. We need to report for the user which file and which line caused an error.
func ParseDocument(doc []byte) (frontmatter any, markdown ast.Node, err error) {
	doc = bytes.TrimSpace(doc)
	if len(doc) == 0 {
		return nil, nil, errors.New("document is empty")
	}

	if !bytes.HasPrefix(doc, frontmatterDelim) {
		return nil, ParseMarkdown(doc), nil
	}

	delimLen := len(frontmatterDelim)

	// Find frontmatter start and end
	frontmatterStart := delimLen
	remaining := doc[frontmatterStart:]
	remaining = bytes.TrimLeft(remaining, "\r\n\t ")

	frontmatterEnd := bytes.Index(remaining, frontmatterDelim)
	if frontmatterEnd == -1 {
		return nil, nil, errors.New("closing frontmatter delimiter not found")
	}

	// Calculate actual byte positions
	rawFrontmatter := bytes.TrimSpace(remaining[:frontmatterEnd])
	markdownStart := frontmatterStart + frontmatterEnd + delimLen + 1
	rawMarkdown := bytes.TrimLeft(doc[markdownStart:], "\r\n\t ")

	parsedMd := ParseMarkdown(rawMarkdown)
	parsedFrontmatter, err := ParseFrontmatter(rawFrontmatter)

	return parsedFrontmatter, parsedMd, err
}
