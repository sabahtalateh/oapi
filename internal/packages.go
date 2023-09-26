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

type mod struct {
	root       string
	modFile    *modfile.File
	discovered map[string]packageInfo
}

type packageInfo struct {
	alias string
	dir   string
}

var m *mod

func findMod(dir string) (*mod, error) {
	if dir == "/" || dir == "." || dir == "" {
		return nil, fmt.Errorf("file should reside within module. (go mod init ..)")
	}

	bb, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	if os.IsNotExist(err) {
		return findMod(filepath.Dir(dir))
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

func ReadPackageName(l Location, imp string) (string, error) {
	var err error

	if m == nil {
		m, err = findMod(filepath.Dir(l.File))
		if err != nil {
			return "", err
		}
	}

	var dir string
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
		m.discovered[imp] = packageInfo{alias: pkg.Name, dir: filepath.Dir(pkg.GoFiles[0])}
		return pkg.Name, nil
	}

	m.discovered[imp] = packageInfo{alias: pkgName, dir: dir}
	return pkgName, nil
}

func PkgDirForImport(l Location, imp string) (string, error) {
	var err error

	if m == nil {
		m, err = findMod(filepath.Dir(l.File))
		if err != nil {
			return "", err
		}
	}

	if d, ok := m.discovered[imp]; ok {
		return d.dir, nil
	}

	for _, req := range m.modFile.Require {
		if strings.HasPrefix(imp, req.Mod.Path) {
			pkgPath := filepath.Join(m.root, "vendor", req.Mod.Path)
			_, err = os.Stat(pkgPath)
			if os.IsNotExist(err) {
				return "", fmt.Errorf("vendor dir may be outdated. try run `go mod vendor`")
			}
			return pkgPath, nil
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
	m.discovered[imp] = packageInfo{alias: pkg.Name, dir: filepath.Dir(pkg.GoFiles[0])}
	return filepath.Dir(pkg.GoFiles[0]), nil
}
