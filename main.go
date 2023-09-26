package main

import (
	"fmt"
	"os"

	"github.com/sabahtalateh/oapi/internal"
	"github.com/sabahtalateh/oapi/internal/operations"
)

func main() {
	wd, err := os.Getwd()
	check(err)

	c, err := internal.ReadConfig(wd)
	check(err)

	l, err := internal.Loc(wd)
	check(err)

	ctx := internal.NewContext(l, c, c.OutFile)
	op, err := operations.Parse(ctx, os.Args[1:])
	checkLoc(l, err)

	n, err := internal.ReadOutputNode(c.OutFile)
	checkLoc(l, err)

	op.Sync(ctx, n.Content)

	err = n.Write(c.Indent)
	checkLoc(l, err)
}

func checkLoc(l internal.Location, err error) {
	if err != nil {
		fmt.Printf("oapi\n%s\n\t%s:%d\n", err, l.File, l.Line)
		os.Exit(1)
	}
}

func check(err error) {
	if err != nil {
		fmt.Printf("oapi\n%s\n", err)
		os.Exit(1)
	}
}
