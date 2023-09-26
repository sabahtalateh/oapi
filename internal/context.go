package internal

import (
	"go/ast"
	"go/token"
	"path/filepath"
)

type Context struct {
	Executable string
	Location   Location
	WorkDir    string
	OutFile    string
	Indent     int
	Verbosity  int

	Imports       []*ast.ImportSpec
	FSet          *token.FileSet
	CurrentSchema string
}

func NewContext(loc Location, c Config, outFile string) Context {
	return Context{
		Executable: c.Executable,
		Location:   loc,
		WorkDir:    filepath.Dir(loc.File),
		OutFile:    outFile,
		Indent:     c.Indent,
		Verbosity:  c.LogVerbosity,
	}
}

func (c Context) WithLocation(l Location) Context {
	return Context{
		Executable:    c.Executable,
		Location:      l,
		WorkDir:       c.WorkDir,
		OutFile:       c.OutFile,
		Indent:        c.Indent,
		Verbosity:     c.Verbosity,
		Imports:       c.Imports,
		FSet:          c.FSet,
		CurrentSchema: c.CurrentSchema,
	}
}

func (c Context) WithImports(ii []*ast.ImportSpec) Context {
	return Context{
		Executable:    c.Executable,
		Location:      c.Location,
		WorkDir:       c.WorkDir,
		OutFile:       c.OutFile,
		Indent:        c.Indent,
		Verbosity:     c.Verbosity,
		Imports:       ii,
		FSet:          c.FSet,
		CurrentSchema: c.CurrentSchema,
	}
}

func (c Context) WithFSet(fset *token.FileSet) Context {
	return Context{
		Executable:    c.Executable,
		Location:      c.Location,
		WorkDir:       c.WorkDir,
		OutFile:       c.OutFile,
		Indent:        c.Indent,
		Verbosity:     c.Verbosity,
		Imports:       c.Imports,
		FSet:          fset,
		CurrentSchema: c.CurrentSchema,
	}
}

func (c Context) WithSchemaType(schemaType string) Context {
	return Context{
		Executable:    c.Executable,
		Location:      c.Location,
		WorkDir:       c.WorkDir,
		OutFile:       c.OutFile,
		Indent:        c.Indent,
		Verbosity:     c.Verbosity,
		Imports:       c.Imports,
		FSet:          c.FSet,
		CurrentSchema: schemaType,
	}
}
