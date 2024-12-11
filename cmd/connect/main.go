package main

import (
	"github.com/choria-io/fisk"
	"github.com/synadia-io/connect/cli"
	"log"
	"log/slog"
	"os"
	"runtime/debug"
)

var version = "0.0.0"

func main() {
	ncli := fisk.New("connect", "Synadia Connect CLI")
	ncli.Author("Synadia Authors <info.synadia.io>")
	ncli.UsageWriter(os.Stdout)
	ncli.Version(getVersion())
	ncli.HelpFlag.Short('h')

	opts, err := cli.ConfigureInApp(ncli, nil, true)
	if err != nil {
		return
	}
	cli.SetVersion(version)

	ncli.Flag("server", "NATS server urls").Short('s').Envar("NATS_URL").PlaceHolder("URL").StringVar(&opts.Servers)
	ncli.Flag("user", "Username or Token").Envar("NATS_USER").PlaceHolder("USER").StringVar(&opts.Username)
	ncli.Flag("password", "Password").Envar("NATS_PASSWORD").PlaceHolder("PASSWORD").StringVar(&opts.Password)
	ncli.Flag("connection-name", "Nickname to use for the underlying NATS Connection").Default("Connect CLI " + version).PlaceHolder("NAME").StringVar(&opts.ConnectionName)
	ncli.Flag("creds", "User credentials").Envar("NATS_CREDS").PlaceHolder("FILE").StringVar(&opts.Creds)
	ncli.Flag("nkey", "User NKEY").Envar("NATS_NKEY").PlaceHolder("FILE").StringVar(&opts.Nkey)
	ncli.Flag("tlscert", "TLS public certificate").Envar("NATS_CERT").PlaceHolder("FILE").ExistingFileVar(&opts.TlsCert)
	ncli.Flag("tlskey", "TLS private key").Envar("NATS_KEY").PlaceHolder("FILE").ExistingFileVar(&opts.TlsKey)
	ncli.Flag("tlsca", "TLS certificate authority chain").Envar("NATS_CA").PlaceHolder("FILE").ExistingFileVar(&opts.TlsCA)
	ncli.Flag("tlsfirst", "Perform TLS handshake before expecting the server greeting").BoolVar(&opts.TlsFirst)
	ncli.Flag("context", "Configuration context").Envar("NATS_CONTEXT").PlaceHolder("NAME").StringVar(&opts.CfgCtx)
	ncli.Flag("timeout", "Time to wait on responses from NATS").Default("5s").Envar("NATS_TIMEOUT").PlaceHolder("DURATION").DurationVar(&opts.Timeout)
	ncli.Flag("no-context", "Disable the selected context").UnNegatableBoolVar(&cli.SkipContexts)
	ncli.Flag("trace", "Enable request tracing").UnNegatableBoolVar(&opts.Trace)
	ncli.Flag("log-level", "Log level to use").Default("info").EnumVar(&opts.LogLevel, "error", "warn", "info", "debug", "trace")

	log.SetFlags(log.Ltime)
	switch opts.LogLevel {
	case "error":
		slog.SetLogLoggerLevel(slog.LevelError)
	case "warn":
		slog.SetLogLoggerLevel(slog.LevelWarn)
	case "info":
		slog.SetLogLoggerLevel(slog.LevelInfo)
	case "debug":
		slog.SetLogLoggerLevel(slog.LevelDebug)
	default:
		slog.SetLogLoggerLevel(slog.LevelInfo)
	}

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
