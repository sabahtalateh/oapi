package internal

import (
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
	"golang.org/x/tools/go/packages"
)

var (
	ErrImportDirNotFound = errors.New("import directory not found")
)

type mod struct {
	root       string
	modFile    *modfile.File
	discovered map[string]packageInfo
}

type packageInfo struct {
	pkgName string
	dir     string
}

var m *mod

func initMod(dir string) (*mod, error) {
	if dir == "/" || dir == "." || dir == "" {
		return nil, fmt.Errorf("file should reside within module. (go mod init ..)")
	}

	bb, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if os.IsNotExist(err) {
		return initMod(filepath.Dir(dir))
	}
	if err != nil {
		return nil, err
	}

	modF, err := modfile.Parse("go.mod", bb, nil)
	if err != nil {
		return nil, err
	}

	return &mod{root: dir, modFile: modF, discovered: map[string]packageInfo{}}, nil
}

func ImportDir(ctx Context, alias string) (string, error) {
	var err error

	if m == nil {
		m, err = initMod(filepath.Dir(ctx.WorkDir))
		if err != nil {
			return "", err
		}
	}

	// suppose package dir equals to last package element (dir: c; package a/b/c)
	for _, imp := range ctx.Imports {
		val := UnquoteImport(imp.Path.Value)

		if imp.Name != nil && imp.Name.Name == alias {
			return pkgDirForImport(val)
		}

		if filepath.Base(val) == alias {
			return pkgDirForImport(val)
		}
	}

	// if dir not equals to last package element (dir d; package a/b/c)
	// then read ".go" file from dir to read actual package
	for _, imp := range ctx.Imports {
		val := UnquoteImport(imp.Path.Value)

		pkgName, err := readPackageName(val)
		if err != nil {
			return "", err
		}

		if pkgName == alias {
			return pkgDirForImport(val)
		}
	}

	return "", errors.Join(ErrImportDirNotFound, errors.New(alias))
}

func readPackageName(imp string) (string, error) {
	var (
		err error
		dir string
	)

	if strings.HasPrefix(imp, m.modFile.Module.Mod.Path) {
		dir = strings.Replace(imp, m.modFile.Module.Mod.Path, m.root, 1)
	} else {
		dir = filepath.Join(m.root, "vendor", imp)
	}

	var pkgName string
	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if pkgName != "" {
			return nil
		}
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		inf, err := d.Info()
		if err != nil {
			return err
		}
		if inf.IsDir() || !strings.HasSuffix(inf.Name(), ".go") {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}

		if !strings.HasSuffix(node.Name.Name, "_test") {
			pkgName = node.Name.Name
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if pkgName == "" {
		pkgs, err := packages.Load(nil, imp)
		if err != nil {
			return "", err
		}
		if len(pkgs) == 0 {
			return "", fmt.Errorf("malformed import: %s", imp)
		}
		pkg := pkgs[0]
		if len(pkg.Errors) != 0 {
			return "", errors.Join(fmt.Errorf("malformed import: %s", imp), pkg.Errors[0])
		}
		m.discovered[imp] = packageInfo{pkgName: pkg.Name, dir: filepath.Dir(pkg.GoFiles[0])}
		return pkg.Name, nil
	}

	m.discovered[imp] = packageInfo{pkgName: pkgName, dir: dir}
	return pkgName, nil
}

func pkgDirForImport(imp string) (string, error) {
	var err error

	if strings.HasPrefix(imp, m.modFile.Module.Mod.Path) {
		return strings.Replace(imp, m.modFile.Module.Mod.Path, m.root, 1), nil
	}

	if d, ok := m.discovered[imp]; ok {
		return d.dir, nil
	}

	for _, req := range m.modFile.Require {
		if strings.HasPrefix(imp, req.Mod.Path) {
			pkgDir := filepath.Join(m.root, "vendor", req.Mod.Path)
			_, err = os.Stat(pkgDir)
			if os.IsNotExist(err) {
				break
			}
			return pkgDir, nil
		}
	}

	pkgs, err := packages.Load(nil, imp)
	if err != nil {
		return "", err
	}
	if len(pkgs) == 0 {
		return "", fmt.Errorf("malformed import: %s", imp)
	}
	pkg := pkgs[0]
	if len(pkg.Errors) != 0 {
		return "", errors.Join(fmt.Errorf("malformed import: %s", imp), pkg.Errors[0])
	}
	m.discovered[imp] = packageInfo{pkgName: pkg.Name, dir: filepath.Dir(pkg.GoFiles[0])}
	return filepath.Dir(pkg.GoFiles[0]), nil
}

func UnquoteImport(imp string) string {
	v := strings.Trim(imp, "\"")
	v = strings.Trim(v, "'")
	v = strings.Trim(v, "`")

	return v
}
