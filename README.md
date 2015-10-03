[![GoDoc](https://godoc.org/github.com/shazow/go-selfcompile?status.svg)](https://godoc.org/github.com/shazow/go-selfcompile)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/shazow/go-selfcompile/master/LICENSE)
[![Build Status](https://travis-ci.org/shazow/go-selfcompile.svg?branch=master)](https://travis-ci.org/shazow/go-selfcompile)

# go-selfcompile

Build self-recompiling Go binaries for embedding new plugins at runtime.

**Status**: `v0` (no stability guarantee); this is a proof of concept, if you'd
like to depend on go-selfcompile then please open an issue with your API
requirements.

Check the project's [upcoming milestones](https://github.com/shazow/go-selfcompile/milestones)
to get a feel for what's prioritized.


## Why?

If you ship a Go-built binary to a user and want to make it easy to install
third-party plugins, what do you do?

Until now, the user would need to install the Go compiler and runtime, make a
stub to import the plugin, re-build the binary with the new dependency and use
that.

go-selfcompile facilitates bundling the Go compiler and runtime, creating the
plugin import stub, recompiling, and replacing the original binary with just a
call to `SelfCompile.Compile()`!


## How does it work?

### Plugins

Let's start with plugins: We define a *plugin* as a package which does something
inside `init() { ... }`. Your system would provide some way for plugins to
register themselves on init, then all you'll need to do is import them and off
you go.

Example of a plugin: [example/aplugin](https://github.com/shazow/go-selfcompile/tree/master/example/aplugin)

### Self-compiling Binary

Next, to use go-selfcompile in your binary you'll need to do two things:

1. **Generate the bundled asset container for the Go compiler and runtime.**
   You'll need our handy `go-selfcompile` binary that you can install with
   `go get github.com/shazow/go-selfcompile/...`.

   Somewhere near your `func main() { ... }`, add a go generate stanza:

   ```go
   //go:generate go-selfcompile --skip-source
   ```

   Now run `go generate` to build the bundle container.

2. **Add the SelfCompile handler in your command line flow.**
   Check [the documentation](https://godoc.org/github.com/shazow/go-selfcompile#SelfCompile)
   for all the options, but a bare minimum would look something like this:

   ```go
    c := selfcompile.SelfCompile{
        Install:       "github.com/shazow/go-selfcompile/example/abinary",
        RestoreAssets: RestoreAssets,
    }
    // Add a plugin from the CLI call
    c.Plugin(plugin)
    // Initiate the compiling with the plugin stubs
    if err := c.Compile(); err != nil { ... }
    // Delete the temporary directory used for compiling
    if err := c.Cleanup(); err != nil { ... }
    ```

Example of a self-compiling binary: [example/abinary](https://github.com/shazow/go-selfcompile/tree/master/example/abinary)

### Examples

If you're trying out the built-in examples, it will look something like this:

```
$ example-abinary
Just doing binary things

$ example-abinary --plugin "github.com/shazow/go-selfcompile/example/aplugin"
Installing plugin: github.com/shazow/go-selfcompile/example/aplugin
[selfcompile] 2015/10/03 15:10:08 Initializing workdir: /tmp/go-selfcompile690079187
[selfcompile] 2015/10/03 15:10:21 Compiling workdir: /tmp/go-selfcompile690079187
[selfcompile] 2015/10/03 15:10:43 Replacing binary: /usr/local/bin/example-abinary
[selfcompile] 2015/10/03 15:10:44 Cleaning up: /tmp/go-selfcompile690079187
Success.

$ example-abinary
aplugin activated.
Just doing binary things
```

Fancy, right?

## Developing

There's an end-to-end integration flow setup in the `Makefile`. You can run it with `make example-aplugin`.


## Sponsors

This project was made possible thanks to [Glider Labs](http://gliderlabs.com/).


## License

MIT
