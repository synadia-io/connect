package main

import (
    "github.com/choria-io/fisk"
    "github.com/synadia-io/connect/cli"
    "os"
    "time"
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

    cli.ConfigureConnectorCommand(ncli, opts)
    cli.ConfigureLibraryCommand(ncli, opts)
    cli.ConfigureLogsCommand(ncli, opts)

    ncli.MustParseWithUsage(os.Args[1:])
}
