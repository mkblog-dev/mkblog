package parser

import (
	"bytes"
	"errors"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"go.abhg.dev/goldmark/frontmatter"
)

type MdFrontmatter struct {
	Title string   `yaml:"title"`
	Tags  []string `yaml:"tags"`
	Desc  string   `yaml:"description"`
}

// TODO: doc should be a structure with metadata. We need to report for the user which file and which line caused an error.
func ParseDocument(doc []byte) (mdAst ast.Node, mdFrontmatter MdFrontmatter, err error) {
	var parsedFrontmatter MdFrontmatter
	doc = bytes.TrimSpace(doc)
	if len(doc) == 0 {
		return nil, parsedFrontmatter, errors.New("document is empty")
	}

	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithExtensions(
			&frontmatter.Extender{},
		),
	)

	ctx := parser.NewContext()
	reader := text.NewReader(doc)
	parsedMd := md.Parser().Parse(reader, parser.WithContext(ctx))

	d := frontmatter.Get(ctx)
	if d == nil {
		return parsedMd, parsedFrontmatter, nil
	}

	if err := d.Decode(&parsedFrontmatter); err != nil {
		return parsedMd, parsedFrontmatter, err
	}

	return parsedMd, parsedFrontmatter, err
}
