package tests

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func WriteOriginal(t *testing.T, orig string) {
	_, filename, _, _ := runtime.Caller(1)
	err := os.WriteFile(filepath.Join(filepath.Dir(filename), "oapi.yaml"), []byte(orig), os.ModePerm)
	require.NoError(t, err)

	t.Log("Original oapi.yaml written")
}

func GoGenerate(t *testing.T, dir string) error {
	execDir := buildExecutable(t)

	cmd := exec.Command("go", "generate", "./...")
	cmd.Dir = dir
	cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=$PATH:%s", execDir))
	cmd.Env = append(cmd.Env, fmt.Sprintf("HOME=%s", os.Getenv("HOME")))

	var stdOut bytes.Buffer
	var stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return errors.Join(err, errors.New(stdOut.String()), errors.New(stdErr.String()))
	}

	t.Log("go generate ./...")
	for _, s := range strings.Split(strings.TrimSpace(stdOut.String()), "\n") {
		t.Log(s)
	}
	return nil
}

func Dir() string {
	_, filename, _, _ := runtime.Caller(1)
	return filepath.Dir(filename)
}

func ReadOAPI(t *testing.T) (string, map[string]any) {
	_, filename, _, _ := runtime.Caller(1)

	bb, err := os.ReadFile(filepath.Join(filepath.Dir(filename), "oapi.yaml"))
	require.NoError(t, err)

	m := map[string]any{}
	err = yaml.Unmarshal(bb, &m)
	require.NoError(t, err)

	return TrimOApi(string(bb)), m
}

// buildExecutable returns executable dir
func buildExecutable(t *testing.T) string {
	_, filename, _, _ := runtime.Caller(2)

	var testDirParts []string
	parts := strings.Split(filename, string(os.PathSeparator))
	if parts[0] == "" {
		parts[0] = "/"
	}
	for _, part := range parts {
		testDirParts = append(testDirParts, part)
		if part == "tests" {
			break
		}
	}

	testsDir := filepath.Join(testDirParts...)
	mainDir := filepath.Dir(testsDir)
	_, err := exec.Command("go", "build", "-o", filepath.Join(filepath.Dir(filename), "oapi"), mainDir).Output()
	require.NoError(t, err)

	t.Logf("Executable built: %s", filepath.Join(testsDir, "oapi"))

	return filepath.Dir(filename)
}

func CleanGenerated(t *testing.T, dir string) {
	_, filename, _, _ := runtime.Caller(1)
	testDir := filepath.Dir(filename)

	err := os.Remove(filepath.Join(testDir, "oapi"))
	if !os.IsNotExist(err) {
		require.NoError(t, err)
	}

	t.Log("Executable removed")

	err = os.Remove(filepath.Join(dir, "oapi.yaml"))
	if !os.IsNotExist(err) {
		require.NoError(t, err)
	}

	t.Log("Generated api removed")
}

func TrimOApi(s string) string {
	outLines := strings.Split(s, "\n")
	for i := 0; i < len(outLines); i++ {
		outLines[i] = strings.TrimRightFunc(outLines[i], unicode.IsSpace)
	}
	return strings.TrimSpace(strings.Join(outLines, "\n"))
}
