// algox — the commands live in internal/cli; this is only where the binary
// starts.
package main

import "algox/internal/cli"

// version is set at build time: -ldflags "-X main.version=…".
var version = "dev"

func main() {
	cli.Main(version)
}
