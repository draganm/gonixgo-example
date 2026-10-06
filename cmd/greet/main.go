// Command greet prints a greeting. It is the example program for gonixgo.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/draganm/gonixgo-example/internal/banner"
	"github.com/draganm/gonixgo-example/internal/greeting"
	"github.com/draganm/gonixgo-example/internal/zstd"
)

// version is set at link time with -X main.version=...
var version = "dev"

func main() {
	name := flag.String("name", "world", "who to greet")
	compress := flag.Bool("compress", false, "write the greeting zstd-compressed, for piping into `zstd -d`")
	showVersion := flag.Bool("version", false, "print the version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s (zstd %s)\n", version, zstd.Version())
		return
	}
	if *compress {
		frame, err := zstd.Compress([]byte(banner.Text() + greeting.For(*name) + "\n"))
		if err != nil {
			fmt.Fprintln(os.Stderr, "greet:", err)
			os.Exit(1)
		}
		os.Stdout.Write(frame)
		return
	}
	fmt.Print(banner.Text())
	color.Green(greeting.For(*name))
}
