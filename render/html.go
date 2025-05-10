package render

import (
	"bytes"
	"html/template"
	"io"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type PageData struct {
	Title   string
	Content template.HTML
}

func RenderHtmlPage(title string, mdAst ast.Node, doc []byte, output io.Writer, tmpl *template.Template) error {
	stripMdExtensionsFromLinks(mdAst)
	html := renderHtml(mdAst, doc)
	data := PageData{
		Title:   title,
		Content: template.HTML(html),
	}

	return tmpl.Execute(output, data)
}

func stripMdExtensionsFromLinks(root ast.Node) {
	ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if link, ok := n.(*ast.Link); ok && entering {
			dest := string(link.Destination)
			if strings.HasSuffix(dest, ".md") {
				// Strip .md extension
				newDest := strings.TrimSuffix(dest, ".md")
				link.Destination = []byte(newDest)
			}
		}
		return ast.WalkContinue, nil
	})
}

func renderHtml(ast ast.Node, doc []byte) []byte {
	var buf bytes.Buffer
	md := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)
	if err := md.Renderer().Render(&buf, doc, ast); err != nil {
		// TODO: figure out when it can panic and why
		panic(err)
	}
	return buf.Bytes()
}
