package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunCompiledTests(t *testing.T) {

	dirname = "example/"

	// Remove any existing opt files from the directory
	fs, _ := filepath.Glob(path.Join(dirname, "*_opt.go"))
	for _, f := range fs {
		os.Remove(f)
	}
	fs, _ = filepath.Glob(path.Join(dirname, "*_opt_test.go"))
	for _, f := range fs {
		os.Remove(f)
	}

	main()

	os.Chdir(dirname)

	args := []string{"test"}
	args = append(args, os.Args[1:]...)

	cmd := exec.Command("go", args...)
	result, err := cmd.CombinedOutput()
	fmt.Println(string(result))
	assert.NoError(t, err)
}
