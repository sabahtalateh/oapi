package seqs

import "gopkg.in/yaml.v3"

func Insert(seq *yaml.Node, i int, ins *yaml.Node) {
	if seq.Kind != yaml.SequenceNode {
		return
	}

	lastI := len(seq.Content) - 1
	if lastI < i {
		seq.Content = append(seq.Content, ins)
		return
	}

	var tail []*yaml.Node
	for j := i; j < len(seq.Content); j++ {
		tail = append(tail, seq.Content[j])
	}

	seq.Content = seq.Content[:i]
	seq.Content = append(seq.Content, ins)
	seq.Content = append(seq.Content, tail...)
}

func Del(seq *yaml.Node, i int) *yaml.Node {
	if seq.Kind != yaml.SequenceNode {
		return nil
	}

	lastI := len(seq.Content) - 1
	if i > lastI {
		return nil
	}

	var tail []*yaml.Node
	for j := i + 1; j < len(seq.Content); j++ {
		tail = append(tail, seq.Content[j])
	}

	del := seq.Content[i]

	seq.Content = seq.Content[:i]
	seq.Content = append(seq.Content, tail...)

	return del
}

func Move(seq *yaml.Node, from, to int) {
	if seq.Kind != yaml.SequenceNode {
		return
	}

	if from == to {
		return
	}

	del := Del(seq, from)
	Insert(seq, to, del)
}
