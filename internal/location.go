package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

type Location struct {
	File string
	Line int
}

func Loc(wd string) (Location, error) {
	var err error

	l := Location{
		File: filepath.Join(wd, os.Getenv("GOFILE")),
	}
	if l.File == "" {
		return Location{}, fmt.Errorf("seems to be run not from //go:generate")
	}
	l.Line, err = strconv.Atoi(os.Getenv("GOLINE"))
	if err != nil {
		return Location{}, errors.Join(err, fmt.Errorf("seems to be run not from //go:generate"))
	}
	return l, err
}
