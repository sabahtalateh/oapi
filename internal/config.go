package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

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

var defaultC = Config{
	OutFile:      "oapi/oapi.yaml",
	Executable:   "oapi",
	Indent:       2,
	LogVerbosity: Verb1,
}

func isModRoot(dir string) bool {
	_, err := os.ReadFile(filepath.Join(dir, "go.mod"))
	return !os.IsNotExist(err)
}

// returns found file & dir where it was found
func findConfigFile(dir string) (*os.File, string, error) {
	path := filepath.Join(dir, confFileName)
	f, err := os.Open(path)
	if os.IsNotExist(err) {
		if isModRoot(dir) {
			return nil, dir, nil
		}

		dir = filepath.Dir(dir)
		if dir == "/" || dir == "." {
			return nil, "", fmt.Errorf("%s not found", confFileName)
		}

		return findConfigFile(dir)
	}
	return f, dir, nil
}

func ReadConfig(wd string) (Config, error) {
	var err error

	confFile, dir, err := findConfigFile(wd)
	if err != nil {
		return defaultC, err
	}

	var c Config
	err = yaml.NewDecoder(confFile).Decode(&c)
	if err != nil {
		return c, errors.Join(err, fmt.Errorf("malformed Config: %s", filepath.Join(dir, confFileName)))
	}

	// set defaults
	if c.OutFile == "" {
		c.OutFile = defaultC.OutFile
	}
	if c.Executable == "" {
		c.Executable = defaultC.Executable
	}
	if c.Indent == 0 {
		c.Indent = defaultC.Indent
	}
	if c.Indent == 0 {
		c.Indent = defaultC.Indent
	}
	if c.LogVerbosity == 0 {
		c.LogVerbosity = defaultC.LogVerbosity
	}

	if !filepath.IsAbs(c.OutFile) {
		c.OutFile, err = filepath.Abs(filepath.Join(dir, c.OutFile))
		if err != nil {
			return Config{}, err
		}
	}

	return c, nil
}
