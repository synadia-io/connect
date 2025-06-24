package cli

import (
    "fmt"
    "log/slog"
    "time"

    "github.com/choria-io/fisk"
    "github.com/nats-io/jsm.go/natscontext"
    "github.com/nats-io/nats.go"
    "github.com/synadia-io/connect/client"
)

var DefaultOptions *Options

func RegisterFlags(app *fisk.Application, version string, opts *Options) {
    app.Flag("server", "NATS server urls").Short('s').Envar("NATS_URL").PlaceHolder("URL").StringVar(&opts.Servers)
    app.Flag("user", "Username or Token").Envar("NATS_USER").PlaceHolder("USER").StringVar(&opts.Username)
    app.Flag("password", "Password").Envar("NATS_PASSWORD").PlaceHolder("PASSWORD").StringVar(&opts.Password)
    app.Flag("connection-name", "Nickname to use for the underlying NATS Connection").Default("NATS Vent CLI Item " + version).PlaceHolder("NAME").StringVar(&opts.ConnectionName)
    app.Flag("creds", "User credentials").Envar("NATS_CREDS").PlaceHolder("FILE").StringVar(&opts.Creds)
    app.Flag("jwt", "User JWT").Envar("NATS_JWT").PlaceHolder("JWT").StringVar(&opts.JWT)
    app.Flag("seed", "User Seed").Envar("NATS_SEED").PlaceHolder("SEED").StringVar(&opts.Seed)
    app.Flag("context", "Configuration context").Envar("NATS_CONTEXT").PlaceHolder("NAME").StringVar(&opts.ContextName)
    app.Flag("timeout", "Time to wait on responses from NATS").Default("5s").Envar("NATS_TIMEOUT").PlaceHolder("DURATION").DurationVar(&opts.Timeout)
    app.Flag("log-level", "Log level to use").Default("info").EnumVar(&opts.LogLevel, "error", "warn", "info", "debug", "trace")
    app.Flag("standalone", "Run in standalone mode without NATS services").BoolVar(&opts.Standalone)
}

// Options configure the CLI
type Options struct {
    // Servers is the list of servers to connect to
    Servers string
    // Username is the username or token to connect with
    Username string
    // Password is the password to connect with
    Password string
    // ConnectionName is the name to use for the underlying NATS connection
    ConnectionName string
    // Creds is nats credentials to authenticate with
    Creds string
    // JWT is the user JWT
    JWT string
    // Seed is the user seed
    Seed string
    // ContextName is the context name to use
    ContextName string

    // Timeout is how long to wait for operations
    Timeout time.Duration

    // LogLevel is the log level to use
    LogLevel string

    // Standalone enables standalone mode without NATS services
    Standalone bool
}

type AppContext struct {
    Nc             *nats.Conn
    DefaultTimeout time.Duration
    Client         client.Client
}

func (ac *AppContext) Close() {
    ac.Nc.Close()
}

func LoadOptions(opts *Options) (*AppContext, error) {
    if opts.Standalone {
        return nil, fmt.Errorf("this command requires NATS services and cannot be used in standalone mode. Use 'connect standalone' commands instead")
    }

    nc, err := loadNats(opts)
    if err != nil {
        return nil, fmt.Errorf("failed to connect to nats: %w", err)
    }

    // -- set the log level
    lvl := slog.LevelInfo
    trace := false
    switch opts.LogLevel {
    case "error":
        lvl = slog.LevelError
    case "warn":
        lvl = slog.LevelWarn
    case "info":
        lvl = slog.LevelInfo
    case "debug":
        lvl = slog.LevelDebug
    case "trace":
        lvl = slog.LevelDebug
        trace = true
    default:
        slog.Warn(fmt.Sprintf("unknown log level %q; reverting to INFO", opts.LogLevel))
    }
    slog.SetLogLoggerLevel(lvl)

    cl, err := client.NewClient(nc, trace)
    if err != nil {
        return nil, fmt.Errorf("failed to create client: %w", err)
    }

    return &AppContext{
        Nc:             nc,
        DefaultTimeout: opts.Timeout,
        Client:         cl,
    }, nil
}

func loadNats(opts *Options) (*nats.Conn, error) {
    if opts == nil {
        opts = DefaultOptions
    }

    if opts.Servers == "" && opts.ContextName == "" {
        opts.ContextName = natscontext.SelectedContext()
    }

    if opts.ContextName != "" {
        return natscontext.Connect(opts.ContextName)
    }

    nopts := []nats.Option{}
    if opts.ConnectionName != "" {
        nopts = append(nopts, nats.Name(opts.ConnectionName))
    }

    if opts.Username != "" {
        if opts.Password == "" {
            nopts = append(nopts, nats.Token(opts.Username))
        } else {
            nopts = append(nopts, nats.UserInfo(opts.Username, opts.Password))
        }
    }

    if opts.Creds != "" {
        nopts = append(nopts, nats.UserCredentials(opts.Creds))
    }

    if opts.JWT != "" && opts.Seed != "" {
        nopts = append(nopts, nats.UserJWTAndSeed(opts.JWT, opts.Seed))
    }

    return nats.Connect(opts.Servers, nopts...)
}
