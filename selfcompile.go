package selfcompile

import (
	"io/ioutil"
	"os"

	"github.com/jteeuwen/go-bindata"
)

// SelfCompile provides controls for registering new plugins and re-compiling
// the binary.
type SelfCompile struct {
	plugins []string
	bindata *bindata.Config

	// Parameters used to setup the temporary workdir.
	Prefix  string // Prefix for TempDir, used to stage recompiling assets.
	Root    string // Root of TempDir (empty will use OS default).
	workdir string // Our workdir once it has been created.
}

// Plugin registers a new plugin to self-compile
func (c *Selfcompile) Plugin(p string) {
	c.plugins = append(c.plugins, p)
}

// Compile will recompile the program's source with the registered plugins.
func (c *Selfcompile) Compile() error {
	// TODO: ...
}

// stubPlugins will generate import stub files for the registered plugins.
func (c *Selfcompile) stubPlugins() error {
	// TODO: Use a mock fs to test: https://talks.golang.org/2012/10things.slide#8
}

// setup will create a fresh temporary directory and inflate all the binary
// data with the appropriate layout inside of it.
func (c *Selfcompile) setup() error {
	var err error
	c.workdir, err = ioutil.TempDir(c.root, c.prefix)
	if err != nil {
		return err
	}
	// TODO: ...
}

// cleanup will delete any temporary files created for the workdir, good idea to
// call this as a defer after calling setup().
func (c *Selfcompile) cleanup() error {
	if c.pworkdir == "" {
		// No workdir setup, nothing to clean up.
		return nil
	}
	return os.RemoveAll(c.workdir)
}
