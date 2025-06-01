package parser

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mkblog-dev/mkblog/utils"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Nav []*NavItem `yaml:"nav"`
}

type NavItem struct {
	Title string `yaml:"title"`
	Href  string `yaml:"href"`
}

func LoadConfig(inputDir string) (*Config, error) {
	pathToBlog, err := utils.RelPathFromCwd(inputDir)

	if err != nil {
		return nil, err
	}

	paths := []string{
		filepath.Join(pathToBlog, "mkblog.yaml"),
		filepath.Join(pathToBlog, "mkblog.yml"),
	}

	for _, path := range paths {
		content, err := os.ReadFile(path)

		if err == nil {
			var root yaml.Node
			if err := yaml.Unmarshal(content, &root); err != nil {
				return nil, err
			}

			var cfg Config
			if err := root.Decode(&cfg); err != nil {
				return nil, err
			}

			for _, doc := range root.Content {
				for i := 0; i < len(doc.Content); i += 2 {
					keyNode := doc.Content[i]
					if keyNode.Value != "nav" {
						continue
					}
					navList := doc.Content[i+1]
					for _, navItem := range navList.Content {
						if navItem.Kind != yaml.MappingNode {
							continue
						}
						var title, href string
						for j := 0; j < len(navItem.Content); j += 2 {
							k := navItem.Content[j]
							v := navItem.Content[j+1]
							switch k.Value {
							case "title":
								title = v.Value
							case "href":
								href = v.Value
							}
						}
						if title == "" {
							return nil,
								fmt.Errorf("invalid nav item at line %d, column %d: missing 'title'", navItem.Line, navItem.Column)
						}
						if href == "" {
							return nil,
								fmt.Errorf("invalid nav item at line %d, column %d: missing 'href'", navItem.Line, navItem.Column)
						}
					}
				}
			}
			return &cfg, nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err // unexpected error
		}
	}

	return nil, errors.New("no mkblog.yaml or mkblog.yml file found")
}
