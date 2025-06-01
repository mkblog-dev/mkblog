package mkblog

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mkblog-dev/mkblog/parser"
	"github.com/mkblog-dev/mkblog/render"
	"github.com/mkblog-dev/mkblog/utils"
)

//go:embed templates
var templates embed.FS

func Build(inputDir string, outputDir string, cfg *parser.Config) error {
	layoutTmpl, err := template.New("layout.tmpl").ParseFS(templates, "templates/*.tmpl")

	if err != nil {
		return err
	}

	pathToBlog, err := utils.RelPathFromCwd(inputDir)
	if err != nil {
		return err
	}
	pathToBlogLen := len(pathToBlog)

	err = filepath.Walk(pathToBlog, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if info.IsDir() {
			return nil
		}

		pathInsideBlog := path[pathToBlogLen:]
		ext := filepath.Ext(pathInsideBlog)
		switch ext {
		case ".md":
			buffer, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			mdAst, mdFrontmatter, err := parser.ParseDocument(buffer)
			if err != nil {
				return err
			}

			err = os.MkdirAll(filepath.Join(outputDir, filepath.Dir(pathInsideBlog)), 0755)
			if err != nil {
				return err
			}

			var f *os.File
			f, err = os.Create(filepath.Join(outputDir, strings.Replace(pathInsideBlog, ".md", ".html", 1)))
			if err != nil {
				return err
			}
			defer f.Close()

			// TODO: should be the first header in the document if empty
			title := mdFrontmatter.Title
			if title == "" {
				title = "Title"
			}

			nav := []*parser.NavItem{}
			if cfg != nil {
				for _, item := range cfg.Nav {
					normalized := *item // shallow copy
					normalized.Href = strings.TrimSuffix(normalized.Href, ".md")
					nav = append(nav, &normalized)
				}
			}

			pageData := render.PageData{
				Title: title,
				Ast:   mdAst,
				Doc:   buffer,
				Nav:   nav,
			}
			render.RenderHtmlPage(&pageData, f, layoutTmpl)
		default:
			// TODO: not to copy dot folders and files.
			// It can be a toggle feature, whether should files to be copied or not by default.
			utils.CopyFile(path, filepath.Join(outputDir, pathInsideBlog))
		}

		return nil
	})

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", inputDir, err)
		return err
	}

	return nil
}
