package render

import (
	"html/template"
	"io"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
)

type PageData struct {
	Title   string
	Content template.HTML
}

func RenderHtmlPage(title string, mdAst ast.Node, output io.Writer, tmpl *template.Template) error {
	html := renderHtml(mdAst)
	data := PageData{
		Title:   title,
		Content: template.HTML(html),
	}

	return tmpl.Execute(output, data)
}

func htmlRenderNodeHook(w io.Writer, node ast.Node, entering bool) (ast.WalkStatus, bool) {
	switch n := node.(type) {
	case *ast.Link:
		if entering {
			dest := string(n.Destination)
			if strings.HasSuffix(dest, ".md") {
				n.Destination = []byte(strings.TrimSuffix(dest, ".md"))
			}
		}
	}
	return ast.GoToNext, false
}

func renderHtml(ast ast.Node) []byte {
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags, RenderNodeHook: htmlRenderNodeHook}
	renderer := html.NewRenderer(opts)

	return markdown.Render(ast, renderer)
}
