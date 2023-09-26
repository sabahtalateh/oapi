package operations

import (
	"fmt"
	"strings"

	"github.com/sabahtalateh/oapi/internal"
)

func logSyncPath(ctx internal.Context, p Path) {
	if ctx.Verbosity == internal.Verb1 {
		return
	}

	if ctx.Verbosity == internal.Verb2 {
		fmt.Printf("path %s %s\n", strings.ToUpper(p.method), p.url)
		return
	}

	if ctx.Verbosity == internal.Verb3 {
		l := ctx.Location
		fmt.Printf("path %s %s\n%s:%d\n\n", strings.ToUpper(p.method), p.url, l.File, l.Line)
		return
	}
}

func logSyncRequest(ctx internal.Context, r Request) {
	if ctx.Verbosity == internal.Verb1 {
		return
	}

	if ctx.Verbosity == internal.Verb2 {
		fmt.Printf("request %s %s\n", r.contentType, r.schema)
		return
	}

	if ctx.Verbosity == internal.Verb3 {
		l := ctx.Location
		fmt.Printf("request %s %s\n%s:%d\n\n", r.contentType, r.schema, l.File, l.Line)
		return
	}
}

func logSyncResponse(ctx internal.Context, r Response) {
	if ctx.Verbosity == internal.Verb1 {
		return
	}

	if ctx.Verbosity == internal.Verb2 {
		fmt.Printf("response %s %s %s\n", r.response, r.contentType, r.schema)
		return
	}

	if ctx.Verbosity == internal.Verb3 {
		l := ctx.Location
		fmt.Printf("response %s %s %s\n%s:%d\n\n", r.response, r.contentType, r.schema, l.File, l.Line)
		return
	}
}

func logSyncSchema(ctx internal.Context, s Schema) {
	if ctx.Verbosity == internal.Verb1 {
		return
	}

	if ctx.Verbosity == internal.Verb2 {
		fmt.Printf("schema %s\n", s.Name)
		return
	}

	if ctx.Verbosity == internal.Verb3 {
		l := ctx.Location
		fmt.Printf("schema %s\n%s:%d\n\n", s.Name, l.File, l.Line)
		return
	}
}
