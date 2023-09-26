package internal

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const PlaceHolder = "<< YOUR TEXT >>"

// TODO usage comments
const emptyOApi = `
openapi: 3.0.3
info:
  title: ` + PlaceHolder + `
  version: 0.0.1
paths:
components: 
  schemas:
`

type Node struct {
	N       *yaml.Node
	Content *yaml.Node
	Path    string
}

func ReadOutputNode(file string) (Node, error) {
	bb, err := os.ReadFile(file)

	if err != nil && !os.IsNotExist(err) {
		return Node{}, err
	}

	if os.IsNotExist(err) || len(bb) == 0 {
		bb = []byte(emptyOApi)
	}

	n := new(yaml.Node)
	err = yaml.Unmarshal(bb, n)
	if err != nil {
		return Node{}, err
	}

	return Node{N: n, Content: n.Content[0], Path: file}, nil
}

func (o Node) Write(indent int) error {
	if err := os.MkdirAll(filepath.Dir(o.Path), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(o.Path)
	if err != nil {
		return err
	}

	enc := yaml.NewEncoder(f)
	enc.SetIndent(indent)

	fixNulls(o.Content)
	return enc.Encode(o.N)
}

// unless node will be rendered as {} which can affect further generation
func fixNulls(n *yaml.Node) {
	for _, node := range n.Content {
		if node.Kind == yaml.MappingNode && len(node.Content) == 0 {
			node.Kind = yaml.ScalarNode
			node.Tag = "!!null"
		}
		fixNulls(node)
	}
}
