package main

import (
	"flag"
	"fmt"

	"github.com/shazow/go-selfcompile"
)

func main() {
	var plugin string
	flag.StringVar(&plugin, "plugin", "", "plugin to install")
	flag.Parse()

	if plugin != "" {
		fmt.Println("Installing plugin: ", plugin)
		c := selfcompile.SelfCompile{
			Install:       "github.com/shazow/go-selfcompile/example/abinary",
			RestoreAssets: RestoreAssets,
		}
		c.Plugin(plugin)
		err := c.Compile()
		if err != nil {
			fmt.Println("Error: ", err.Error())
			return
		}
		fmt.Println("Success.")
		return
	}

	fmt.Println("Just doing binary things.")
}
