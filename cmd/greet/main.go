// Command greet prints a greeting. It is the example program for gonixgo.
package main

import (
	"flag"
	"fmt"

	"github.com/fatih/color"

	"github.com/draganm/gonixgo-example/internal/banner"
	"github.com/draganm/gonixgo-example/internal/greeting"
)

// version is set at link time with -X main.version=...
var version = "dev"

func main() {
	name := flag.String("name", "world", "who to greet")
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		return
	}
	fmt.Print(banner.Text())
	color.Green(greeting.For(*name))
}
