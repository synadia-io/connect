package runtime

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/nats-io/nats.go"
	"github.com/synadia-io/connect/model"
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
		r.loggerSet = true
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
		LogLevel: slog.LevelInfo,
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

	// loggerSet records whether a logger was supplied via WithLogger, so Launch
	// does not clobber an injected logger.
	loggerSet bool
}

// logger returns the runtime's logger, falling back to the default if one has
// not been configured yet (e.g. NatsConfig called before Launch).
func (r *Runtime) logger() *slog.Logger {
	if r.Logger != nil {
		return r.Logger
	}
	return slog.Default()
}

// NatsOptions builds the NATS connection options. The connection is configured
// to survive credential and connectivity loss rather than permanently aborting:
//
//   - IgnoreAuthErrorAbort keeps the client reconnecting through auth errors, so
//     a rotated/expired credential is not treated as a permanent failure (which
//     today forces a manual restart).
//   - MaxReconnects(-1) overrides the default 60-attempt cap so a long outage
//     does not close the connection for good. Reconnect backoff is the client
//     default (fixed wait + jitter).
//   - Credentials prefer a creds file: nats.UserCredentials re-reads it on every
//     (re)connect handshake, so a re-minted credential is picked up on the next
//     reconnect automatically. Otherwise the decoded JWT+seed from the
//     environment is used (static, no refresh).
//   - Lifecycle handlers surface disconnects/reconnects/closure that are silent
//     today.
//   - NoCallbacksAfterClientClose suppresses the lifecycle handlers once we
//     intentionally Close()/Drain() the connection, so an orderly shutdown does
//     not emit a misleading disconnect Warn.
func (r *Runtime) NatsOptions() ([]nats.Option, error) {
	log := r.logger()

	opts := []nats.Option{
		nats.MaxReconnects(-1),
		nats.IgnoreAuthErrorAbort(),
		nats.NoCallbacksAfterClientClose(),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Warn("nats disconnected", slog.Any("err", err))
		}),
		nats.ReconnectHandler(func(c *nats.Conn) {
			log.Info("nats reconnected", slog.String("url", c.ConnectedUrl()))
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			log.Error("nats connection closed")
		}),
		nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
			log.Error("nats async error", slog.Any("err", err))
		}),
	}

	if path := os.Getenv(NatsCredsFileVar); path != "" {
		opts = append(opts, nats.UserCredentials(path))
	} else if r.NatsJwt != "" {
		opts = append(opts, nats.UserJWTAndSeed(r.NatsJwt, r.NatsSeed))
	}

	return opts, nil
}

func (r *Runtime) NatsConfig() (*nats.Conn, error) {
	opts, err := r.NatsOptions()
	if err != nil {
		return nil, err
	}
	return nats.Connect(r.NatsUrl, opts...)
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

	// Build a logger honoring the configured LogLevel, unless one was injected
	// via WithLogger (which must be preserved).
	if !r.loggerSet {
		r.Logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: r.LogLevel}))
	}

	return workload(ctx, r, steps)
}

func (r *Runtime) Close() {}
