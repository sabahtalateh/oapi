package node

import "gopkg.in/yaml.v3"

func Scalar(v any) *yaml.Node {
	switch vv := v.(type) {
	case string:
		return &yaml.Node{Kind: yaml.ScalarNode, Value: vv}
	default:
		return &yaml.Node{Kind: yaml.ScalarNode}
	}
}

func Map(content ...*yaml.Node) *yaml.Node {
	return &yaml.Node{Kind: yaml.MappingNode, Tag: "!!map", Content: content}
}

func Sequence() *yaml.Node {
	return &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq"}
}

func SequenceOf[T any](vv ...T) *yaml.Node {
	n := Sequence()

	for _, v := range vv {
		n.Content = append(n.Content, Scalar(v))
	}

	return n
}
