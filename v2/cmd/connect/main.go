package main

import (
	"os"
	"runtime/debug"
	"time"

	"github.com/choria-io/fisk"
	"github.com/synadia-io/connect/v2/cli"
)

var version = "0.0.0"
var opts *cli.Options

func main() {
	version = getVersion()

	ncli := fisk.New("connect", "Synadia Connect CLI")
	ncli.Author("Synadia Authors <info.synadia.io>")
	ncli.UsageWriter(os.Stdout)
	ncli.Version(version)
	ncli.HelpFlag.Short('h')

	opts = &cli.Options{
		Timeout: 5 * time.Second,
	}

	cli.RegisterFlags(ncli, version, opts)

	cli.ConfigureConnectorCommand(ncli, opts)
	cli.ConfigureLibraryCommand(ncli, opts)

	ncli.MustParseWithUsage(os.Args[1:])
}

func getVersion() string {
	if version != "0.0.0" {
		return version
	}

	nfo, ok := debug.ReadBuildInfo()
	if !ok || (nfo != nil && nfo.Main.Version == "") {
		return version
	}

	return nfo.Main.Version
}
