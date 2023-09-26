package operations

import (
	"regexp"
	"sort"

	xMaps "golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"

	"github.com/sabahtalateh/oapi/internal"
	"github.com/sabahtalateh/oapi/internal/node"
	"github.com/sabahtalateh/oapi/internal/node/maps"
	"github.com/sabahtalateh/oapi/internal/node/seqs"
)

// turns /hello/{id}/{name} into /hello/{}/{}
var normPathRE = regexp.MustCompile(`{.*?}`)

func (p Path) sync(n *yaml.Node) {
	pN := pathN(n, p)

	// {method}
	methodN := maps.KV(pN, p.method, maps.Present().Get(), maps.Absent(node.Map()).Alpha())

	keyIdx := 0

	// {method} -> summary
	maps.KV(methodN, "summary", maps.Present().Move(keyIdx), maps.Absent(node.Scalar(internal.PlaceHolder)).Pos(keyIdx))
	keyIdx++

	keyIdx = syncParams(keyIdx, methodN, p.params)

	keyIdx = syncRequests(keyIdx, methodN, p.requests)
	_ = syncResponses(keyIdx, methodN, p.responses)
}

func syncParams(keyIdx int, n *yaml.Node, pp []param) int {
	// if no path parameters
	if len(pp) == 0 {
		// find parameters in yaml
		paramsN := maps.KV(n, "parameters", maps.Present().Move(keyIdx), nil)
		if paramsN == nil {
			return keyIdx
		}

		// make it sequence if not
		if paramsN.Kind != yaml.SequenceNode {
			paramsN.Kind = yaml.SequenceNode
			paramsN.Tag = "!!seq"
			paramsN.Content = nil
		}

		// remove all path parameters from yaml node
		removePathParams(paramsN, nil)

		// if no parameters on yaml left after removing path parameters
		// then remove entire parameters node
		if len(paramsN.Content) == 0 {
			maps.KV(n, "parameters", maps.Present().Del(), nil)
			return keyIdx
		}
		return keyIdx + 1
	}

	// if path parameters exists
	// find parameters node
	paramsN := maps.KV(n, "parameters", maps.Present().Move(keyIdx), maps.Absent(node.Map()).Pos(keyIdx))

	// make it sequence if not
	if paramsN.Kind != yaml.SequenceNode {
		paramsN.Kind = yaml.SequenceNode
		paramsN.Tag = "!!seq"
		paramsN.Content = nil
	}

	// for each path parameter
	var inPaths []string
	for i, p := range pp {
		if p.in == "path" {
			inPaths = append(inPaths, p.name)
		}

		// find corresponding yaml node
		pI, pN := findParam(paramsN, p.name)

		// if yaml node not found then create it
		// and insert at it's position
		if pN == nil {
			seqs.Insert(paramsN, i, syncParamN(nil, p))
			continue
		}

		// if yaml node found then sync it's content
		pN = syncParamN(pN, p)
		// and move onto corresponding position
		seqs.Move(paramsN, pI, i)
	}

	// remove path parameters except ones in path
	removePathParams(paramsN, inPaths)

	return keyIdx + 1
}

func syncRequests(keyIdx int, n *yaml.Node, requests []Request) int {
	// if no requests then remove them from yaml and return
	if len(requests) == 0 {
		maps.KV(n, "requestBody", maps.Present().Del(), nil)
		return keyIdx
	}

	// else create request body node
	// requestBody
	reqBodyN := maps.KV(n, "requestBody", maps.Present().Move(keyIdx), maps.Absent(node.Map()).Pos(keyIdx))

	// requestBody -> description
	maps.KV(reqBodyN, "description", maps.Present().Move(0), maps.Absent(node.Scalar(internal.PlaceHolder)).Pos(0))

	// requestBody -> required: true
	trueN := node.Scalar("true")
	maps.KV(reqBodyN, "required", maps.Present().Move(1).Replace(trueN), maps.Absent(trueN).Pos(1))

	// requestBody -> content
	contentN := maps.KV(reqBodyN, "content", maps.Present().Move(2), maps.Absent(node.Map()).Pos(2))

	byContentType := map[string][]Request{}
	for _, r := range requests {
		byContentType[r.contentType] = append(byContentType[r.contentType], r)
	}

	contentTypes := xMaps.Keys(byContentType)
	sort.Strings(contentTypes)

	// remove non-existing requestBody -> content -> {content type}
	for _, k := range maps.Keys(contentN) {
		if !slices.Contains(contentTypes, k) {
			maps.KV(contentN, k, maps.Present().Del(), nil)
		}
	}

	for i, contentType := range contentTypes {
		// requestBody -> content -> {content type}
		contentTypeN := maps.KV(contentN, contentType, maps.Present().Move(i), maps.Absent(node.Map()).Pos(i))

		// requestBody -> content -> {content type} -> schema
		schemaN := maps.KV(contentTypeN, "schema", maps.Present().Move(0), maps.Absent(node.Map()).Pos(0))
		reqs := byContentType[contentType]

		if len(reqs) == 0 {
			// will never happen because of if len(requests) == 0. just for completeness

			// remove requestBody -> content -> {content type} -> schema -> oneOf: ..
			maps.KV(schemaN, "oneOf", maps.Present().Del(), nil)

			// requestBody -> content -> {content type} -> schema -> $ref:
			reqN := node.Scalar("")
			maps.KV(schemaN, "$ref", maps.Present().Replace(reqN), maps.Absent(reqN).Pos(0))
		} else if len(reqs) == 1 {
			// remove requestBody -> content -> {content type} -> schema -> oneOf: ..
			maps.KV(schemaN, "oneOf", maps.Present().Del(), nil)

			// requestBody -> content -> {content type} -> schema -> $ref: {path to schema}
			reqN := node.Scalar(reqs[0].ref)
			maps.KV(schemaN, "$ref", maps.Present().Move(0).Replace(reqN), maps.Absent(reqN).Pos(0))
		} else {
			// remove requestBody -> content -> {content type} -> schema -> $ref: ..
			maps.KV(schemaN, "$ref", maps.Present().Del(), nil)

			oneOfN := node.Sequence()
			var refs []string
			for _, req := range reqs {
				refs = append(refs, req.ref)
			}
			sort.Strings(refs)

			for _, ref := range refs {
				m := node.Map()
				maps.KV(m, "$ref", nil, maps.Absent(node.Scalar(ref)).Pos(0))
				oneOfN.Content = append(oneOfN.Content, m)
			}

			// requestBody -> content -> {content type} -> schema -> oneOf: [$ref: .., $ref: .., ..]
			maps.KV(schemaN, "oneOf", maps.Present().Move(0).Replace(oneOfN), maps.Absent(oneOfN).Pos(0))
		}
	}

	return keyIdx + 1
}

func syncResponses(keyIdx int, n *yaml.Node, responses []Response) int {
	// responses
	responsesN := maps.KV(n, "responses", maps.Present().Move(keyIdx), maps.Absent(node.Map()).Pos(keyIdx))
	println(responsesN)

	byCode := map[string][]Response{}
	for _, r := range responses {
		byCode[r.response] = append(byCode[r.response], r)
	}

	codes := xMaps.Keys(byCode)
	sort.Strings(codes)

	// remove non-existing responses -> {response code}
	for _, k := range maps.Keys(responsesN) {
		if !slices.Contains(codes, k) {
			maps.KV(responsesN, k, maps.Present().Del(), nil)
		}
	}

	for i, code := range codes {
		// responses -> {response code}
		codeN := maps.KV(responsesN, code, maps.Present().Move(i), maps.Absent(node.Map()).Pos(i))

		// responses -> {response code} -> description
		maps.KV(codeN, "description", maps.Present().Move(0), maps.Absent(node.Scalar(internal.PlaceHolder)).Pos(0))

		// responses -> {response code} -> content
		contentN := maps.KV(codeN, "content", maps.Present().Move(1), maps.Absent(node.Map()).Pos(1))
		byContentType := map[string][]Response{}

		codeResponses := byCode[code]
		for _, r := range codeResponses {
			byContentType[r.contentType] = append(byContentType[r.contentType], r)
		}

		contentTypes := xMaps.Keys(byContentType)
		sort.Strings(contentTypes)

		for i, contentType := range contentTypes {
			// remove non-existing responses -> {response code} -> content -> {content type}
			for _, k := range maps.Keys(contentN) {
				if !slices.Contains(contentTypes, k) {
					maps.KV(contentN, k, maps.Present().Del(), nil)
				}
			}

			// responses -> {response code} -> content -> {content type}
			contentTypeN := maps.KV(contentN, contentType, maps.Present().Move(i), maps.Absent(node.Map()).Pos(i))

			// responses -> {response code} -> content -> {content type} -> schema
			schemaN := maps.KV(contentTypeN, "schema", maps.Present().Move(0), maps.Absent(node.Map()).Pos(0))
			println(schemaN)

			resps := byContentType[contentType]

			if len(resps) == 0 {
				// will never happen. just for completeness

				// remove responses -> {response code} -> content -> {content type} -> schema -> oneOf: ..
				maps.KV(schemaN, "oneOf", maps.Present().Del(), nil)

				// responses -> {response code} -> content -> {content type} -> schema -> $ref:
				respN := node.Scalar("")
				maps.KV(schemaN, "$ref", maps.Present().Replace(respN), maps.Absent(respN).Pos(0))
			} else if len(resps) == 1 {
				// remove responses -> {response code} -> content -> {content type} -> schema -> oneOf: ..
				maps.KV(schemaN, "oneOf", maps.Present().Del(), nil)

				// responses -> {response code} -> content -> {content type} -> schema -> $ref: ..
				respN := node.Scalar(resps[0].ref)
				maps.KV(schemaN, "$ref", maps.Present().Replace(respN), maps.Absent(respN).Pos(0))
			} else {
				// remove responses -> {response code} -> content -> {content type} -> schema -> $ref
				maps.KV(schemaN, "$ref", maps.Present().Del(), nil)

				oneOfN := node.Sequence()
				var refs []string
				for _, resp := range resps {
					refs = append(refs, resp.ref)
				}
				sort.Strings(refs)

				for _, ref := range refs {
					m := node.Map()
					maps.KV(m, "$ref", nil, maps.Absent(node.Scalar(ref)).Pos(0))
					oneOfN.Content = append(oneOfN.Content, m)
				}

				// remove responses -> {response code} -> content -> {content type} -> schema -> oneOf: [$ref: .., $ref: .., ..]
				maps.KV(schemaN, "oneOf", maps.Present().Move(0).Replace(oneOfN), maps.Absent(oneOfN).Pos(0))
			}
		}
	}

	return keyIdx + 1
}

func syncParamN(n *yaml.Node, p param) *yaml.Node {
	if n == nil {
		n = node.Map()
	}

	// name
	nameN := node.Scalar(p.name)
	maps.KV(n, "name", maps.Present().Move(0).Replace(nameN), maps.Absent(nameN).Pos(0))

	// in
	inN := node.Scalar(p.in)
	maps.KV(n, "in", maps.Present().Move(1).Replace(inN), maps.Absent(inN).Pos(1))

	// description
	maps.KV(n, "description", maps.Present().Move(2), maps.Absent(node.Scalar(internal.PlaceHolder)).Pos(2))

	// required
	reqN := node.Scalar(b2s(p.required))
	maps.KV(n, "required", maps.Present().Move(3).Replace(reqN), maps.Absent(reqN).Pos(3))

	// schema
	schemaN := maps.KV(n, "schema", maps.Present().Move(4), maps.Absent(node.Map()).Pos(4))

	// schema -> type
	typeN := node.Scalar(p.typ)
	maps.KV(schemaN, "type", maps.Present().Move(0).Replace(typeN), maps.Absent(typeN).Pos(0))

	return n
}

func findParam(n *yaml.Node, name string) (int, *yaml.Node) {
	for i, nod := range n.Content {
		nameN := maps.KV(nod, "name", maps.Present().Get(), nil)
		if nameN == nil {
			continue
		}
		if nameN.Value == name {
			return i, nod
		}
	}

	return 0, nil
}

func removePathParams(n *yaml.Node, keep []string) {
	if n.Kind != yaml.SequenceNode {
		return
	}
	var newContent []*yaml.Node
	for _, node := range n.Content {
		inN := maps.KV(node, "in", maps.Present().Get(), nil)
		if inN == nil {
			continue
		}
		in := inN.Value

		var paramName string
		nameN := maps.KV(node, "name", maps.Present().Get(), nil)
		if nameN != nil {
			paramName = nameN.Value
		}

		if in != "path" ||
			(in == "path" && slices.Contains(keep, paramName)) {
			newContent = append(newContent, node)
		}
	}
	n.Content = newContent
}

func pathN(n *yaml.Node, p Path) *yaml.Node {
	lastI := len(n.Content) - 1
	for i := 0; i < len(n.Content); i += 2 {
		normPathVal := normPathRE.ReplaceAllString(n.Content[i].Value, "{}")
		if normPathVal == p.normUrl {
			n.Content[i].Value = p.url
			if lastI < i {
				n.Content = append(n.Content, node.Map())
			}
			return n.Content[i+1]
		}
	}

	pN := maps.KV(n, p.url, nil, maps.Absent(node.Map()).Alpha())
	return pN
}
