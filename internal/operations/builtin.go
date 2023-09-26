package operations

var builtinTypes = map[string]struct{}{
	"bool":       {},
	"uint8":      {},
	"uint16":     {},
	"uint32":     {},
	"uint64":     {},
	"int8":       {},
	"int16":      {},
	"int32":      {},
	"int64":      {},
	"float32":    {},
	"float64":    {},
	"complex64":  {},
	"complex128": {},
	"string":     {},
	"int":        {},
	"uint":       {},
	"uintptr":    {},
	"byte":       {},
	"rune":       {},
	"any":        {},
}

// https://spec.openapis.org/oas/v3.0.3#data-types
func builtinToType(builtin string) (string, string) {
	switch builtin {
	case "bool":
		return "boolean", ""
	case "uint8":
		return "integer", "int32"
	case "uint16":
		return "integer", "int32"
	case "uint32":
		return "integer", "int32"
	case "uint64":
		return "integer", "int64"
	case "int8":
		return "integer", "int32"
	case "int16":
		return "integer", "int32"
	case "int32":
		return "integer", "int32"
	case "int64":
		return "integer", "int64"
	case "float32":
		return "number", "float"
	case "float64":
		return "number", "double"
	case "complex64":
		return "string", ""
	case "complex128":
		return "string", ""
	case "string":
		return "string", ""
	case "int":
		return "integer", "int64"
	case "uint":
		return "integer", "int64"
	case "uintptr":
		return "integer", "int64"
	case "byte":
		return "integer", "byte"
	case "rune":
		return "string", ""
	case "any":
		return "string", ""
	default:
		return "string", ""
	}
}
