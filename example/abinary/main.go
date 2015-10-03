//go:generate go-selfcompile
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shazow/go-selfcompile"
)

func main() {
	var plugin string
	flag.StringVar(&plugin, "plugin", "", "plugin to install")

	var printBundled bool
	flag.BoolVar(&printBundled, "bundled", false, "print bundled assets and exit")

	flag.Parse()
	selfcompile.SetLogger(os.Stderr)

	if printBundled {
		fmt.Println("Embedded assets:")
		for _, name := range AssetNames() {
			fmt.Println(" *", name)
		}
		return
	}

	if plugin != "" {
		fmt.Println("Installing plugin: ", plugin)
		c := selfcompile.SelfCompile{
			Install:       "github.com/shazow/go-selfcompile/example/abinary",
			RestoreAssets: RestoreAssets,
		}
		c.Plugin(plugin)
		if err := c.Compile(); err != nil {
			fmt.Println("Compile failed:", err.Error())
			return
		}
		if err := c.Cleanup(); err != nil {
			fmt.Println("Cleanup failed:", err.Error())
			return
		}
		fmt.Println("Success.")
		return
	}

	fmt.Println("Just doing binary things.")
}
