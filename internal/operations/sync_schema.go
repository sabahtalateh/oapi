package operations

import (
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"

	"github.com/sabahtalateh/oapi/internal/node"
	"github.com/sabahtalateh/oapi/internal/node/maps"
)

func (s Schema) sync(n *yaml.Node) {
	// {schema}
	schemaN := maps.KV(n, s.Name, maps.Present().Get(), maps.Absent(node.Map()).Alpha())
	switch st := s.Schema.(type) {
	case Object:
		delArrayFields(schemaN)
		syncObject(schemaN, st)
	case Array:
		delObjectFields(schemaN)
		syncArray(schemaN, st)
	}
}

func syncObject(n *yaml.Node, obj Object) {
	withProps := len(obj.Properties) != 0
	withRefs := len(obj.Refs) != 0

	// if no props and no refs then insert empty object
	if !withProps && !withRefs {
		delObjectFields(n)

		// {schema} -> type: object
		typeN := node.Scalar("object")
		maps.KV(n, "type", nil, maps.Absent(typeN).Pos(0))

		return
	}

	// if only props then insert type: object, ..
	if withProps && !withRefs {
		n = objectValue(n, false)

		keyIdx := 0

		// {schema} -> type: object
		typeN := node.Scalar("object")
		maps.KV(n, "type", maps.Present().Move(0).Replace(typeN), maps.Absent(typeN).Pos(0))
		keyIdx++

		syncProperties(keyIdx, n, obj.Properties, obj.Required)

		// remove allOf
		maps.KV(n, "allOf", maps.Present().Del(), nil)

		return
	}

	// if refs exists then create sequence [$ref: .., $ref: ..]
	allOfContent := refNodes(obj.Refs)

	// if pros exists then add object [- $ref: .., - $ref: .., - type: object, ..]
	if withProps {
		objN := objectValue(n, true)

		keyIdx := 0

		// {schema} -> allOf -> [.., - type: object, ..]
		typeN := node.Scalar("object")
		maps.KV(objN, "type", maps.Present().Move(0).Replace(typeN), maps.Absent(typeN).Pos(0))
		keyIdx++

		syncProperties(keyIdx, objN, obj.Properties, obj.Required)
		allOfContent = append(allOfContent, objN)
	}

	// {schema} -> allOf
	allOfN := &yaml.Node{Kind: yaml.SequenceNode, Tag: "!!seq", Content: allOfContent}
	maps.KV(n, "allOf", maps.Present().Move(0).Replace(allOfN), maps.Absent(allOfN).Pos(0))

	// remove object fields
	maps.KV(n, "type", maps.Present().Del(), nil)
	maps.KV(n, "required", maps.Present().Del(), nil)
	maps.KV(n, "properties", maps.Present().Del(), nil)
	maps.KV(n, "items", maps.Present().Del(), nil)
}

// TODO add required !!
func syncProperties(keyIdx int, objN *yaml.Node, props []Property, required []string) int {
	if len(props) == 0 {
		maps.KV(objN, "required", maps.Present().Del(), nil)
		maps.KV(objN, "properties", maps.Present().Del(), nil)
		return keyIdx
	}

	if len(required) == 0 {
		maps.KV(objN, "required", maps.Present().Del(), nil)
	}

	if len(required) != 0 {
		// {schema} -> required
		reqsN := node.SequenceOf(required...)
		maps.KV(objN, "required", maps.Present().Move(keyIdx).Replace(reqsN), maps.Absent(reqsN).Pos(keyIdx))
		keyIdx += 1
	}

	// {schema} -> properties
	propsN := maps.KV(objN, "properties", maps.Present().Move(keyIdx), maps.Absent(node.Map()).Pos(keyIdx))

	removeAbsentProps(propsN, props)

	for i, prop := range props {
		// {schema} -> properties -> {property}
		propN := maps.KV(propsN, prop.Name, maps.Present().Move(i), maps.Absent(node.Map()).Pos(i))

		switch v := prop.Val.(type) {
		case Type:
			// {schema} -> properties -> {property} -> type: ..
			typeN := node.Scalar(v.Type)
			maps.KV(propN, "type", maps.Present().Move(0).Replace(typeN), maps.Absent(typeN).Pos(0))
			if v.Format != "" {
				// {schema} -> properties -> {property} -> format: ..
				maps.KV(propN, "format", maps.Present().Move(1), maps.Absent(node.Scalar(v.Format)).Pos(1))
			}

			maps.KV(propN, "items", maps.Present().Del(), nil)
			maps.KV(propN, "$ref", maps.Present().Del(), nil)
		case Ref:
			// {schema} -> properties -> {property} -> $ref: ..
			refN := node.Scalar(v.Ref)
			maps.KV(propN, "$ref", maps.Present().Move(0).Replace(refN), maps.Absent(refN).Pos(0))

			maps.KV(propN, "type", maps.Present().Del(), nil)
			maps.KV(propN, "format", maps.Present().Del(), nil)
			maps.KV(propN, "items", maps.Present().Del(), nil)
		case Array:
			syncArray(propN, v)
		}
	}

	return keyIdx + 1
}

func syncArray(n *yaml.Node, a Array) {
	typeN := node.Scalar("array")
	maps.KV(n, "type", maps.Present().Move(0).Replace(typeN), maps.Absent(typeN).Pos(0))
	itemsN := maps.KV(n, "items", maps.Present().Move(1), maps.Absent(node.Map()).Pos(1))

	switch ii := a.Items.(type) {
	case Type:
		typeN := node.Scalar(ii.Type)
		maps.KV(itemsN, "type", maps.Present().Move(0).Replace(typeN), maps.Absent(typeN).Pos(0))
		if ii.Format != "" {
			maps.KV(itemsN, "format", maps.Present().Move(1), maps.Absent(node.Scalar(ii.Format)).Pos(1))
		}
		maps.KV(itemsN, "$ref", maps.Present().Del(), nil)
	case Ref:
		refN := node.Scalar(ii.Ref)
		maps.KV(itemsN, "$ref", maps.Present().Move(0).Replace(refN), maps.Absent(refN).Pos(0))

		maps.KV(itemsN, "type", maps.Present().Del(), nil)
		maps.KV(itemsN, "format", maps.Present().Del(), nil)
	}

	maps.KV(n, "format", maps.Present().Del(), nil)
	maps.KV(n, "$ref", maps.Present().Del(), nil)
}

func objectValue(objN *yaml.Node, withinAllOf bool) *yaml.Node {
	// if object not contains allOf then return this object
	allOfN := maps.KV(objN, "allOf", maps.Present().Get(), nil)
	if allOfN == nil {
		if !withinAllOf {
			return objN
		}

		// to prevent cycles
		objNCopy := node.Map(objN.Content...)
		objN.Content = nil

		allOfN = node.Sequence()
		allOfN.Content = append(allOfN.Content, objNCopy)

		maps.KV(objN, "allOf", nil, maps.Absent(allOfN).Pos(0))
	}

	if allOfN.Kind != yaml.SequenceNode {
		return objN
	}

	// starting from here working with allOf
	var (
		objEl *yaml.Node
	)

	for _, el := range allOfN.Content {
		// return first found element with `type` property
		typeN := maps.KV(el, "type", maps.Present().Get(), nil)
		if typeN != nil {
			objEl = el
			break
		}
	}

	// returns object element from allOf
	if withinAllOf {
		if objEl == nil {
			objEl = node.Map()
			allOfN.Content = append(allOfN.Content, objEl)
		}
		return objEl
	}

	// replace node content with element from allOf
	if objEl == nil {
		objEl = node.Map()
	}

	// TODO remove allOf and keep other original content
	objN.Content = objEl.Content
	return objN
}

func refNodes(rr []string) []*yaml.Node {
	var out []*yaml.Node

	for _, ref := range rr {
		m := node.Map()
		maps.KV(m, "$ref", nil, maps.Absent(node.Scalar(ref)).Pos(0))
		out = append(out, m)
	}

	return out
}

func removeAbsentProps(propsN *yaml.Node, props []Property) {
	var presented []string
	for _, prop := range props {
		presented = append(presented, prop.Name)
	}

	var newContent []*yaml.Node
	lastI := len(propsN.Content) - 1
	for i := 0; i < len(propsN.Content); i += 2 {
		key := propsN.Content[i].Value
		if slices.Contains(presented, key) {
			newContent = append(newContent, propsN.Content[i])
			if lastI >= i+1 {
				newContent = append(newContent, propsN.Content[i+1])
			}
		}
	}
	propsN.Content = newContent
}

func delObjectFields(n *yaml.Node) {
	maps.KV(n, "type", maps.Present().DelIf(func(y *yaml.Node) bool { return y.Value == "object" }), nil)
	maps.KV(n, "required", maps.Present().Del(), nil)
	maps.KV(n, "properties", maps.Present().Del(), nil)
	maps.KV(n, "allOf", maps.Present().Del(), nil)
}

func delArrayFields(n *yaml.Node) {
	maps.KV(n, "type", maps.Present().DelIf(func(y *yaml.Node) bool { return y.Value == "array" }), nil)
	maps.KV(n, "items", maps.Present().Del(), nil)
}

func b2s(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
