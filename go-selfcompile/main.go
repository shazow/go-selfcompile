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

func detectGOROOT() (string, error) {
	cmd := exec.Command("go", "env")
	defer cmd.Wait()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	if err := cmd.Start(); err != nil {
		return "", err
	}

	in := bufio.NewScanner(stdout)
	for in.Scan() {
		line := in.Text()
		if !strings.HasPrefix(line, "GOROOT=") {
			continue
		}
		return line[len("GOROOT=\\") : len(line)-1], nil
	}
	if err := in.Err(); err != nil {
		return "", err
	}
	return "", errDetectGoRoot
}

func inputConfigs(goroot string) []bindata.InputConfig {
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
		bindata.InputConfig{Path: filepath.Join(goroot, "pkg", "tool", "asm")},
		bindata.InputConfig{Path: filepath.Join(goroot, "pkg", "tool", "compile")},
		bindata.InputConfig{Path: filepath.Join(goroot, "pkg", "tool", "link")},
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

	goroot, err := detectGOROOT()
	if err != nil {
		exit(1, fmt.Sprintf("failed detecting GOROOT: %v", err))
	}

	// Default paths
	cfg.Input = inputConfigs(goroot)
	cfg.Prefix = goroot

	err = bindata.Translate(cfg)
	if err != nil {
		exit(1, err.Error())
	}
}
