// Package cli holds algox's commands. main.go calls Main.
package cli

import (
	"flag"
	"fmt"
	"os"

	"algox/internal/gui"
)

const usage = `ALGOX

Usage:
  algox              open the desktop app
  algox app          open the desktop app
  algox web [flags]  serve the page to a browser, with no desktop app
  algox version      print the version

Flags of web:
`

// Main runs the command named by os.Args.
func Main(version string) {
	if err := run(version, os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "algox:", err)
		os.Exit(1)
	}
}

func run(version string, args []string) error {
	cmd := ""
	if len(args) > 0 {
		cmd, args = args[0], args[1:]
	}
	switch cmd {
	case "", "app":
		return gui.Run(version)
	case "web":
		return web(version, args)
	case "version", "-v", "--version":
		fmt.Println(version)
		return nil
	case "help", "-h", "--help":
		fmt.Print(usage)
		webFlags().PrintDefaults()
		return nil
	}
	return fmt.Errorf("unknown command %q; see `algox help`", cmd)
}

type webOpts struct {
	fs   *flag.FlagSet
	host *string
	port *int
}

func webFlags() *flag.FlagSet { return newWebOpts().fs }

func newWebOpts() webOpts {
	fs := flag.NewFlagSet("web", flag.ContinueOnError)
	return webOpts{
		fs:   fs,
		host: fs.String("host", "127.0.0.1", `address to listen on; "" or 0.0.0.0 for every address`),
		port: fs.Int("port", gui.DefaultPort, "port to listen on"),
	}
}

func web(version string, args []string) error {
	o := newWebOpts()
	if err := o.fs.Parse(args); err != nil {
		return err
	}
	// ALGOX_WEB_KEY keeps one key across runs, for a service.
	return gui.ServeWeb(version, *o.host, *o.port, os.Getenv("ALGOX_WEB_KEY"))
}
