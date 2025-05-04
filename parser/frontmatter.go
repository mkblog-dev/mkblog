package parser

import (
	"gopkg.in/yaml.v2"
)

func ParseFrontmatter(doc []byte) (frontmatter any, err error) {
	var result any
	if err := yaml.Unmarshal(doc, &result); err != nil {
		return nil, err
	}
	return result, nil
}
