package selfcompile

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

const srcdir = "_self"
const vendordir = "_vendor"
const pluginfile = "plugin_selfcompile.go"
const tmpprefix = "go-selfcompile"

var errRestoreAssets = errors.New("missing RestoreAssets")

type RestoreAssets func(dir, name string) error

// SelfCompile provides controls for registering new plugins and re-compiling
// the binary.
type SelfCompile struct {
	pkg     string // Package to use for plugins (empty will default to "main")
	plugins []string

	// Main package source URL to install on recompile (if not bundled).
	Install string

	// Automatically call SelfCompile.Cleanup() after Compile() is done.
	AutoCleanup bool

	// Parameters used to setup the temporary workdir.
	Prefix    string // Prefix for TempDir, used to stage recompiling assets.
	Root      string // Root of TempDir (empty will use OS default).
	workdir   string // Full path to the work dir once it has been created.
	srcdir    string // Full path to source dir of our package.
	vendordir string // Full path to GOPATH dir for our dependencies.

	// RestoreAssets is the function generated by bindata to restore the assets
	// recursively within a given directory.
	RestoreAssets RestoreAssets
}

// Plugin registers a plugin to self-compile. Make sure to register the full
// set of plugins that need to be enabled during the compile, not just new
// plugins.
func (c *SelfCompile) Plugin(p string) {
	c.plugins = append(c.plugins, p)
}

// Compile the program's source with the registered plugins.
func (c *SelfCompile) Compile() (err error) {
	err = c.setup()
	if err != nil {
		return
	}
	if c.AutoCleanup {
		defer func() {
			err = combineErrors(c.Cleanup(), err)
		}()
	}

	if c.Install == "" {
		// TODO: Handle bundled source if c.Install is not defined
		err = errors.New("not implemented: Bundled source, must specify Install target.")
		return
	}

	logger.Println("Compiling workdir:", c.workdir)

	err = c.goRun("get", c.Install)
	if err != nil {
		return err
	}

	self, err := selfPath()
	if err != nil {
		return err
	}

	logger.Println("Replacing binary:", self)
	err = c.copyFile(filepath.Join(c.vendordir, "bin", filepath.Base(c.Install)), self)
	return
}

func (c *SelfCompile) copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_TRUNC, 0777)
	if err != nil {
		return err
	}
	defer out.Close()
	// TODO: Use a buffered copy?
	_, err = io.Copy(out, in)
	return err
}

func (c *SelfCompile) goRun(args ...string) error {
	// FIXME: Default to env = os.Environ()?
	binpath := filepath.Join(c.workdir, "bin")
	env := []string{
		fmt.Sprintf("PATH=%s:%s", binpath, os.Getenv("PATH")),
		fmt.Sprintf("GOROOT=%s", c.workdir),
		fmt.Sprintf("GOPATH=%s", c.vendordir),
	}

	cmd := exec.Cmd{
		Path: filepath.Join(binpath, "go"),
		Args: append([]string{"go"}, args...),
		Env:  env,
		Dir:  c.workdir,

		// TODO: Eat outputs and do something with them?
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	return cmd.Run()
}

// stubPlugins will generate import stub files for the registered plugins.
func (c *SelfCompile) stubPlugins() error {
	// TODO: Use a mock fs to test: https://talks.golang.org/2012/10things.slide#8
	path := filepath.Join(c.srcdir, pluginfile)
	fd, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fd.Close()

	w := bufio.NewWriter(fd)
	defer w.Flush()

	p := plugin{
		Package: c.pkg,
		Imports: c.plugins,
	}
	_, err = p.WriteTo(w)
	return err
}

// setup will create a fresh temporary directory and inflate all the binary
// data with the appropriate layout inside of it.
func (c *SelfCompile) setup() error {
	var err error
	if c.RestoreAssets == nil {
		return errRestoreAssets
	}

	prefix := tmpprefix
	if c.Prefix != "" {
		prefix = c.Prefix
	}
	c.workdir, err = ioutil.TempDir(c.Root, prefix)
	if err != nil {
		return err
	}
	logger.Printf("Initializing workdir: %s", c.workdir)

	c.vendordir = filepath.Join(c.workdir, vendordir)
	if c.Install == "" {
		// Assume we embedded the source
		c.srcdir = filepath.Join(c.workdir, srcdir)
	} else {
		c.srcdir = filepath.Join(c.vendordir, "src", c.Install)
	}

	// Restore all the assets recursively
	err = c.RestoreAssets(c.workdir, "")
	if err != nil {
		return err
	}

	if c.Install != "" {
		// Fetch source
		err := c.goRun("get", "-d", c.Install)
		if err != nil {
			return err
		}
	}

	// Generate plugin stubs in srcdir
	err = c.stubPlugins()
	if err != nil {
		return err
	}

	// go generate
	err = c.goRun("generate", c.Install)
	if err != nil {
		return err
	}

	return nil
}

// Cleanup will delete any temporary files created for the workdir, good idea to
// call this as a defer after calling setup().
func (c *SelfCompile) Cleanup() error {
	if c.workdir == "" {
		// No workdir setup, nothing to clean up.
		return nil
	}
	logger.Printf("Cleaning up: %s", c.workdir)
	return os.RemoveAll(c.workdir)
}

// selfPath returns the path of the current running process's binary.
func selfPath() (string, error) {
	path, err := exec.LookPath(os.Args[0])
	if err == nil && path != "" {
		return path, err
	}
	return filepath.Abs(os.Args[0])
}
