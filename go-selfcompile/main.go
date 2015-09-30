// go-selfcompile binary is a helper wrapper around go-bindata for embedding
// the necessary assets to use SelfCompile.
package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jteeuwen/go-bindata"
)

var errDetectGoRoot = errors.New("failed to detect GOROOT")

func goEnv() (map[string]string, error) {
	// TODO: Load from os.Environ() too?
	env := map[string]string{}
	cmd := exec.Command("go", "env")
	defer cmd.Wait()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return env, err
	}
	if err := cmd.Start(); err != nil {
		return env, err
	}

	in := bufio.NewScanner(stdout)
	for in.Scan() {
		line := in.Text()
		parts := strings.SplitN(line, "=", 2)
		k, v := parts[0], parts[1]
		env[k] = strings.Trim(v, `"`)
	}
	return env, in.Err()
}

func inputConfigs(goroot string, gotooldir string) []bindata.InputConfig {
	return []bindata.InputConfig{
		// Minimum artifacts required for `go build` to work.
		// See: https://github.com/shazow/go-selfcompile/issues/2
		bindata.InputConfig{
			Path:      filepath.Join(goroot, "src"),
			Recursive: true,
		},
		bindata.InputConfig{
			Path:      filepath.Join(goroot, "pkg", "include"),
			Recursive: true,
		},
		bindata.InputConfig{Path: filepath.Join(gotooldir, "asm")},
		bindata.InputConfig{Path: filepath.Join(gotooldir, "compile")},
		bindata.InputConfig{Path: filepath.Join(gotooldir, "link")},
		bindata.InputConfig{Path: filepath.Join(goroot, "bin", "go")},
	}
}

func exit(code int, msg string) {
	fmt.Fprintf(os.Stderr, "go-selfcompile: %s\n", msg)
	os.Exit(code)
}

func main() {
	cfg := bindata.NewConfig()
	cfg.Output = "bindata_selfcompile.go"

	env, err := goEnv()
	if err != nil {
		exit(1, fmt.Sprintf("failed loading go env: %v", err))
	}

	goroot := env["GOROOT"]
	if goroot == "" {
		exit(1, fmt.Sprintf("failed detecting GOROOT"))
	}

	gotooldir := env["GOTOOLDIR"]
	if gotooldir == "" {
		exit(1, fmt.Sprintf("failed detecting GOTOOLDIR"))
	}

	// Default paths
	cfg.Input = inputConfigs(goroot, gotooldir)
	cfg.Prefix = goroot

	err = bindata.Translate(cfg)
	if err != nil {
		exit(1, err.Error())
	}
}
