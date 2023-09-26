package maps

import (
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"

	"github.com/sabahtalateh/oapi/internal/node"
)

func MakePath(n *yaml.Node, path ...string) *yaml.Node {
	x := n
	for _, p := range path {
		x = KV(x, p, Present().Get(), Absent(node.Map()).Append())
	}
	return x
}

// IfPresent what to do if key exists
type IfPresent struct {
	isGet       bool                  // return presented map value
	isDel       bool                  // delete
	delFunc     func(*yaml.Node) bool // delete function
	isMove      bool                  // reposition at
	atPos       int                   // position
	isLast      bool                  // move at last position
	replaceWith *yaml.Node            // replace with node
}

func Present() *IfPresent {
	return &IfPresent{}
}

func (p *IfPresent) Get() *IfPresent {
	p.isGet = true
	return p
}

func (p *IfPresent) Del() *IfPresent {
	p.isDel = true
	return p
}

func (p *IfPresent) DelIf(f func(*yaml.Node) bool) *IfPresent {
	p.isDel = true
	p.delFunc = f
	return p
}

func (p *IfPresent) Move(pos int) *IfPresent {
	p.isMove = true
	p.atPos = pos
	return p
}

func (p *IfPresent) Last() *IfPresent {
	p.isLast = true
	return p
}

func (p *IfPresent) Replace(n *yaml.Node) *IfPresent {
	p.replaceWith = n
	return p
}

// IfAbsent what to do if key not exists
type IfAbsent struct {
	isAlpha  bool       // add in alphabetical order
	isPos    bool       // add at
	atPos    int        // position
	isAppend bool       // add at the end
	node     *yaml.Node // node to add
}

func Absent(n *yaml.Node) *IfAbsent {
	return &IfAbsent{node: n}
}

func (a *IfAbsent) Alpha() *IfAbsent {
	a.isAlpha = true
	return a
}

func (a *IfAbsent) Pos(p int) *IfAbsent {
	a.isPos = true
	a.atPos = p
	return a
}

func (a *IfAbsent) Append() *IfAbsent {
	a.isAppend = true
	return a
}

func Keys(n *yaml.Node) []string {
	if n.Kind != yaml.MappingNode {
		return nil
	}

	var keys []string
	for i := 0; i < len(n.Content); i += 2 {
		keys = append(keys, n.Content[i].Value)
	}

	return keys
}

func KV(n *yaml.Node, key string, pres *IfPresent, abs *IfAbsent) *yaml.Node {
	n.Kind = yaml.MappingNode
	n.Tag = "!!map"

	var valN *yaml.Node
	ki, found := keyIndex(n, key)
	if found {
		valN = valueByKeyIndex(n, ki)
	}

	if found {
		if pres != nil {
			if pres.isDel && (pres.delFunc == nil || pres.delFunc(valN)) {
				delKeyByIndex(n, ki)
				return nil
			}
			if pres.isMove {
				moveFromIndexToIndex(n, ki, pres.atPos*2)
				ki, _ = keyIndex(n, key)
				valN = valueByKeyIndex(n, ki)
			}
			if pres.isLast {
				moveFromIndexToIndex(n, ki, len(n.Content))
				ki, _ = keyIndex(n, key)
				valN = valueByKeyIndex(n, ki)
			}
			if pres.replaceWith != nil {
				n.Content[ki+1] = pres.replaceWith
			}
		}
		return valN
	}

	if abs == nil {
		return nil
	}

	var keys []string
	for i := 0; i < len(n.Content); i += 2 {
		keys = append(keys, n.Content[i].Value)
	}

	keyN := &yaml.Node{Kind: yaml.ScalarNode, Value: key}
	valN = abs.node
	insertPos := -1

	if abs.isAlpha {
		if !slices.IsSorted(keys) {
			n.Content = append(n.Content, keyN, valN)
			insertPos = len(n.Content)
		} else {
			for i := 0; i < len(n.Content); i += 2 {
				mKey := n.Content[i].Value
				if key < mKey {
					insertPos = i
					break
				}
			}

			if insertPos == -1 {
				insertPos = len(n.Content)
			}
		}
	}

	if abs.isPos {
		insertPos = abs.atPos * 2
	}

	if insertPos == -1 {
		insertPos = len(n.Content)
	}

	insertKeyValueAtIndex(n, insertPos, keyN, valN)

	return valN
}

func keyIndex(n *yaml.Node, key string) (int, bool) {
	for i := 0; i < len(n.Content); i += 2 {
		kN := n.Content[i]
		if key == kN.Value {
			return i, true
		}
	}

	return 0, false
}

func valueByKeyIndex(n *yaml.Node, keyI int) *yaml.Node {
	lastI := len(n.Content) - 1
	if keyI+1 > lastI {
		n.Content = append(n.Content, &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null"})
	}
	return n.Content[keyI+1]
}

func insertKeyValueAtIndex(n *yaml.Node, keyI int, k *yaml.Node, v *yaml.Node) {
	head := n.Content[0:keyI]

	// tail must be copied unless it will be modified by further appends
	var tail []*yaml.Node
	for j := keyI; j < len(n.Content); j++ {
		tail = append(tail, n.Content[j])
	}

	newContent := append(head, k, v)
	newContent = append(newContent, tail...)
	n.Content = newContent
}

func delKeyByIndex(n *yaml.Node, keyI int) (*yaml.Node, *yaml.Node) {
	if n.Kind != yaml.MappingNode {
		return nil, nil
	}

	lastI := len(n.Content) - 1
	if keyI > lastI {
		return nil, nil
	}

	keyN := n.Content[keyI]
	var valN *yaml.Node
	if keyI+1 > lastI {
		valN = &yaml.Node{Kind: yaml.ScalarNode, Tag: "!!null"}
	} else {
		valN = n.Content[keyI+1]
	}

	newContent := n.Content[:keyI]
	newContent = append(newContent, n.Content[keyI+2:]...)
	n.Content = newContent

	return keyN, valN
}

func moveFromIndexToIndex(n *yaml.Node, from, to int) {
	if from == to {
		return
	}

	remK, remV := delKeyByIndex(n, from)
	if remK == nil {
		return
	}

	insertKeyValueAtIndex(n, to, remK, remV)
}
