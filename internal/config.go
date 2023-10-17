package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	mod2 "github.com/sabahtalateh/mod"
	"gopkg.in/yaml.v3"
)

const (
	confFileName = "oapi.conf.yaml"
)

const (
	Verb1 = 1 // no logs
	Verb2 = 2 // short logs
	Verb3 = 3 // full logs
)

type Config struct {
	OutFile      string `yaml:"out_file"`
	Executable   string `yaml:"executable"`
	Indent       int    `yaml:"indent"`
	LogVerbosity int    `yaml:"log_verbosity"`
}

var defaultConf = Config{
	OutFile:      "oapi/oapi.yaml",
	Executable:   "oapi",
	Indent:       2,
	LogVerbosity: Verb3,
}

func isModRoot(dir string) bool {
	_, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	return !os.IsNotExist(err)
}

// returns found file & dir where it was found
func findConfig(dir, prevDir string, first bool) (*os.File, string, error) {
	if !first && dir == prevDir {
		return nil, "", os.ErrNotExist
	}

	path := filepath.Join(dir, confFileName)
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		if isModRoot(dir) {
			return nil, "", os.ErrNotExist
		}
		return findConfig(filepath.Dir(dir), dir, false)
	}
	if err != nil {
		return nil, "", err
	}
	return f, dir, nil
}

func FindConfig(wd string) (Config, error) {
	var (
		conf Config
		err  error
	)

	confFile, confDir, err := findConfig(wd, wd, true)
	if err != nil {
		conf = defaultConf
		goModPath, err := mod2.ModFilePath(wd)
		if err != nil {
			return conf, fmt.Errorf("not resides within module: %s", wd)
		}
		confDir = filepath.Dir(goModPath)
	} else {
		err = yaml.NewDecoder(confFile).Decode(&conf)
		if err != nil {
			return conf, errors.Join(err, fmt.Errorf("malformed config: %s", filepath.Join(confDir, confFileName)))
		}
	}

	// set defaults
	if conf.OutFile == "" {
		conf.OutFile = defaultConf.OutFile
	}
	if conf.Executable == "" {
		conf.Executable = defaultConf.Executable
	}
	if conf.Indent == 0 {
		conf.Indent = defaultConf.Indent
	}
	if conf.Indent == 0 {
		conf.Indent = defaultConf.Indent
	}
	if conf.LogVerbosity == 0 {
		conf.LogVerbosity = defaultConf.LogVerbosity
	}

	if !filepath.IsAbs(conf.OutFile) {
		conf.OutFile, err = filepath.Abs(filepath.Join(confDir, conf.OutFile))
		if err != nil {
			return Config{}, err
		}
	}

	return conf, nil
}
