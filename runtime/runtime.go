package runtime

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/synadia-io/connect/v2/model"
	"gopkg.in/yaml.v3"
)

type Opt func(*Runtime)

func WithLogLevel(level slog.Level) Opt {
	return func(r *Runtime) {
		r.LogLevel = level
	}
}

func WithNamespace(ns string) Opt {
	return func(r *Runtime) {
		r.Namespace = ns
	}
}

func WithGroup(group string) Opt {
	return func(r *Runtime) {
		r.Connector = group
	}
}

func WithInstance(instance string) Opt {
	return func(r *Runtime) {
		r.Instance = instance
	}
}

func WithNatsUrl(url string) Opt {
	return func(r *Runtime) {
		r.NatsUrl = url
	}
}

func WithNatsJwt(jwt string) Opt {
	return func(r *Runtime) {
		r.NatsJwt = jwt
	}
}

func WithNatsSeed(seed string) Opt {
	return func(r *Runtime) {
		r.NatsSeed = seed
	}
}

func WithLogger(logger *slog.Logger) Opt {
	return func(r *Runtime) {
		r.Logger = logger
	}
}

func FromEnv() (*Runtime, error) {
	opts := []Opt{
		WithNamespace(os.Getenv(NamespaceEnvVar)),
		WithInstance(os.Getenv(InstanceEnvVar)),
		WithGroup(os.Getenv(GroupEnvVar)),
		WithNatsSeed(os.Getenv(NatsSeedVar)),
		WithNatsUrl(os.Getenv(NatsUrlVar)),
	}

	if ll := os.Getenv(LogLevelEnvVar); ll != "" {
		switch strings.ToLower(os.Getenv(LogLevelEnvVar)) {
		case "debug":
			opts = append(opts, WithLogLevel(slog.LevelDebug))
		case "info":
			opts = append(opts, WithLogLevel(slog.LevelInfo))
		case "warn":
			opts = append(opts, WithLogLevel(slog.LevelWarn))
		case "error":
			opts = append(opts, WithLogLevel(slog.LevelError))
		}
	}

	if jwt := os.Getenv(NatsJwtVar); jwt != "" {
		j, err := base64.StdEncoding.DecodeString(jwt)
		if err != nil {
			return nil, fmt.Errorf("failed to decode nats jwt: %w", err)
		}

		opts = append(opts, WithNatsJwt(string(j)))
	}

	return NewRuntime(opts...), nil
}

func NewRuntime(opts ...Opt) *Runtime {
	result := &Runtime{
		LogLevel: slog.LevelDebug,
		Logger:   slog.Default(),
	}

	for _, opt := range opts {
		opt(result)
	}

	return result
}

type Workload func(ctx context.Context, runtime *Runtime, steps model.Steps) error

type Runtime struct {
	Namespace string
	Connector string
	Instance  string

	NatsUrl  string
	NatsJwt  string
	NatsSeed string

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
