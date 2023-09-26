package operations

import "errors"

const responseUsage = `usage: on func. together with ` + "`" + `go:generate oapi path ..` + "`" + `
//go:generate oapi response {HTTP Code | default} {Content-Type (can be omitted, default is applications/json)} {Schema Name}
examples:
//go:generate oapi response 200 applications/json Hello
//go:generate oapi response default Error`

const requestUsage = `usage: on func. together with ` + "`" + `go:generate oapi path ..` + "`" + `
//go:generate oapi request {Content-Type (can be omitted, default is applications/json)} {Schema Name}
examples:
//go:generate oapi request applications/json Hello
//go:generate oapi request Hello`

var (
	errUnknownOperation  = errors.New("unknown operation")
	errMalformedPath     = errors.New("malformed path")
	errMalformedRequest  = errors.New("malformed request")
	errMalformedResponse = errors.New("malformed response")
	errMalformedSchema   = errors.New("malformed schema")
	errPathUsage         = errors.New("usage: on func\n//go:generate oapi path GET /hello")
	errRequestUsage      = errors.New(requestUsage)
	errResponseUsage     = errors.New(responseUsage)
	errSchemaUsage       = errors.New("usage: on type\n//go:generate schema Hello")
	errUsage             = errors.Join(errPathUsage, errSchemaUsage, errResponseUsage)
)

func e(s string) error {
	return errors.New(s)
}
