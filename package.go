package selfcompile

import "regexp"

// Package will generate bindata necessary for SelfCompile to work. This
// includes a Go compiler for the current GOARCH, the runtime sources, and the
// source of the host binary.
type Package struct {
	// Path to the GOROOT directory for the desired GOARCH. Can get this value from running `go env`.
	GoRoot string
	// Path to the host binary's source code.
	BinRoot string
	// Ignore any filenames matching the regex pattern specified (passed to bindata).
	// Useful for ignoring version control directories and other development-specific assets.
	Ignore []*regexp.Regexp
}

func (p Package) Translate() error {
	return nil
}
