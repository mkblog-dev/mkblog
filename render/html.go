package render

import (
	"bytes"
	"html/template"
	"io"
	"strings"

	mkParser "github.com/mkblog-dev/mkblog/parser"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type templateData struct {
	Title   string
	Content template.HTML
	Nav     []*mkParser.NavItem
}

type PageData struct {
	Title string
	Ast   ast.Node
	Doc   []byte
	Nav   []*mkParser.NavItem
}

func RenderHtmlPage(pageData *PageData, output io.Writer, tmpl *template.Template) error {
	stripMdExtensionsFromLinks(pageData.Ast)
	html := renderHtml(pageData.Ast, pageData.Doc)
	data := templateData{
		Title:   pageData.Title,
		Content: template.HTML(html),
		Nav:     pageData.Nav,
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
