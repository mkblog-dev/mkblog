package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/urfave/cli/v3"
)

func mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

func Build(inputDir string, outputDir string) error {
	// we always work with paths relative to CWD
	pathToBlog, err := RelPathFromCwd(inputDir)
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

func main() {
	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "d",
				Usage:    "Input directory",
				Required: true,
			},
			&cli.StringFlag{
				Name:     "o",
				Usage:    "Output directory",
				Required: true,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			inputDir := cmd.String("d")
			outputDir := cmd.String("o")

			if inputDir == "" || outputDir == "" {
				log.Println("Both -d (input) and -o (output) must be specified.")
				os.Exit(1)
			}

			err := Build(inputDir, outputDir)
			if err != nil {
				log.Fatalf("build failed: %v", err)
			}

			log.Println("done")
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func HandleGet(w http.ResponseWriter, r *http.Request) {
	// md := []byte(mds)
	md, err := os.ReadFile("test.md")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
		return
	}
	html := mdToHTML(md)

	// fmt.Printf("--- Markdown:\n%s\n\n--- HTML:\n%s\n", md, htmli
	w.Write([]byte(fmt.Sprintf(`<head><link rel="stylesheet" href="https://cdn.simplecss.org/simple.min.css"></head><body>%s</body>`, html)))
}

func RelPathFromCwd(path string) (string, error) {
	cwd, err := GetCwd()
	if err != nil {
		return "", err
	}

	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	return filepath.Rel(cwd, abs)
}

func GetCwd() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	exPath := filepath.Dir(ex)
	return exPath, nil
}
