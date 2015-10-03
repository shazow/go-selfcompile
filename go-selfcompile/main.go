// go-selfcompile binary is a helper wrapper around go-bindata for embedding
// the necessary assets to use SelfCompile.
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jteeuwen/go-bindata"
)

var version = "dev"

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

func goInputs(goroot string, gotooldir string) []bindata.InputConfig {
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

// selfPath returns the path of the current running process's binary.
func selfPath() (string, error) {
	path, err := exec.LookPath(os.Args[0])
	if err == nil && path != "" {
		return path, err
	}
	return filepath.Abs(os.Args[0])
}

func exit(code int, msg string) {
	fmt.Fprintf(os.Stderr, "go-selfcompile: %s\n", msg)
	os.Exit(code)
}

type options struct {
	ShowVersion bool
	SkipSource  bool
	Out         string
}

func main() {
	opts := options{}
	flag.BoolVar(&opts.ShowVersion, "version", false, "print version and exit")
	flag.BoolVar(&opts.SkipSource, "skip-source", false, "skip embedding package (will have to specify SelfCompile.Install target)")
	flag.StringVar(&opts.Out, "out", "bindata_selfcompile.go", "write bindata to this file")
	flag.Parse()

	if opts.ShowVersion {
		exit(0, fmt.Sprintf("version %s", version))
	}

	cfg := bindata.NewConfig()
	cfg.Output = opts.Out
	cfg.Debug = true // Assets don't need to be bundled in source, only in the built binary.

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

	selfcompilePath, err := selfPath()
	if err != nil {
		exit(1, fmt.Sprintf("failed detecting path of go-selfcompile"))
	}

	// Default paths
	cfg.Prefix = goroot
	cfg.Input = goInputs(goroot, gotooldir)

	// Bundle go-selfcompile.
	cfg.Input = append(cfg.Input, bindata.InputConfig{
		Path: selfcompilePath,
		Name: "bin/go-selfcompile",
	})

	if !opts.SkipSource {
		// Append source to cfg.Input with some default ignore settings.
		// TODO: ...
		exit(2, fmt.Sprintf("not implemented yet: embedding source"))
	}

	err = bindata.Translate(cfg)
	if err != nil {
		exit(1, err.Error())
	}
}
