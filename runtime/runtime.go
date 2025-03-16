package runtime

import (
    "context"
    "encoding/base64"
    "fmt"
    "github.com/nats-io/nats.go"
    "github.com/synadia-io/connect/model"
    "gopkg.in/yaml.v3"
    "log/slog"
    "os"
    "strings"
)

func FromEnv() (*Runtime, error) {
    logLevel := slog.LevelDebug
    switch strings.ToLower(os.Getenv(LogLevelEnvVar)) {
    case "debug":
        logLevel = slog.LevelDebug
    case "info":
        logLevel = slog.LevelInfo
    case "warn":
        logLevel = slog.LevelWarn
    case "error":
        logLevel = slog.LevelError
    }

    return NewRuntime(logLevel,
        os.Getenv(NamespaceEnvVar),
        os.Getenv(GroupEnvVar),
        os.Getenv(InstanceEnvVar)), nil
}

func NewRuntime(logLevel slog.Level, namespace string, connector string, instance string) *Runtime {
    return &Runtime{
        Namespace: namespace,
        Connector: connector,
        Instance:  instance,
        LogLevel:  logLevel,
    }
}

type Workload func(ctx context.Context, runtime *Runtime, steps model.Steps) error

type Runtime struct {
    Namespace string
    Connector string
    Instance  string

    // LogLevel is the log level for the runtime
    LogLevel slog.Level

    // Logger is the logger for the runtime and only set after launch
    Logger *slog.Logger
}

func (r *Runtime) NatsConfig() (*nats.Conn, error) {
    return nats.Connect(os.Getenv(NatsUrlVar),
        nats.UserJWTAndSeed(os.Getenv(NatsJwtVar), os.Getenv(NatsSeedVar)),
    )
}

func (r *Runtime) Launch(ctx context.Context, workload Workload, cfg string) error {
    cfgb, err := base64.StdEncoding.DecodeString(cfg)
    if err != nil {
        return fmt.Errorf("failed to decode config: %w", err)
    }

    // -- decode the connector config
    var steps model.Steps
    if err := yaml.Unmarshal(cfgb, &steps); err != nil {
        return fmt.Errorf("failed to decode connector config: %w", err)
    }

    r.Logger = slog.Default()

    return workload(ctx, r, steps)
}

func (r *Runtime) Close() {}
