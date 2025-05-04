package parser

import (
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/parser"
)

func ParseMarkdown(doc []byte) (ast ast.Node) {
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	normalizedMd := parser.NormalizeNewlines(doc)
	return p.Parse(normalizedMd)
}
