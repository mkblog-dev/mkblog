package mkblog

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	normalizedMd := parser.NormalizeNewlines(md)
	doc := p.Parse(normalizedMd)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags, RenderNodeHook: htmlRenderNodeHook}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func Build(inputDir string, outputDir string) error {
	// we always work with paths relative to CWD
	pathToBlog, err := relPathFromCwd(inputDir)
	if err != nil {
		return err
	}
	pathToBlogLen := len(pathToBlog)

	err = filepath.Walk(pathToBlog, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		pathInsideBlog := path[pathToBlogLen:]
		ext := filepath.Ext(pathInsideBlog)
		switch ext {
		case ".md":
			buffer, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			html := mdToHTML(buffer)

			err = os.MkdirAll(filepath.Join(outputDir, filepath.Dir(pathInsideBlog)), 0755)
			if err != nil {
				return err
			}
			err = os.WriteFile(filepath.Join(outputDir, strings.Replace(pathInsideBlog, ".md", ".html", 1)), html, 0644)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", inputDir, err)
		return err
	}

	return nil
}

func relPathFromCwd(path string) (string, error) {
	cwd, err := getCwd()
	if err != nil {
		return "", err
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return filepath.Rel(cwd, abs)
}

func getCwd() (string, error) {
	return os.Getwd()
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
