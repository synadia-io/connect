package main

import (
	"os"
	"time"

	"github.com/choria-io/fisk"
	"github.com/synadia-io/connect/v2/cli"
)

var opts *cli.Options

func main() {
	ncli := fisk.New("connect", "Synadia Connect CLI")
	ncli.Author("Synadia Authors <info.synadia.io>")
	ncli.UsageWriter(os.Stdout)
	ncli.Version(Version)
	ncli.HelpFlag.Short('h')

	opts = &cli.Options{
		Timeout: 5 * time.Second,
	}

	cli.RegisterFlags(ncli, Version, opts)

	configureVersionCommand(ncli)

	// Configure all commands - standalone will be added conditionally
	cli.ConfigureConnectorCommand(ncli, opts)
	cli.ConfigureLibraryCommand(ncli, opts)
	cli.ConfigureLogsCommand(ncli, opts)
	cli.ConfigureStandaloneCommand(ncli, opts)

	ncli.MustParseWithUsage(os.Args[1:])
}
