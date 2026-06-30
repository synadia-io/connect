package runtime

import (
	"context"
	"encoding/base64"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/synadia-io/connect/model"
	"gopkg.in/yaml.v3"
)

const (
	// reconnectBaseDelay is the first backoff interval; subsequent attempts
	// double up to reconnectMaxDelay.
	reconnectBaseDelay = 1 * time.Second
	reconnectMaxDelay  = 30 * time.Second
)

// reconnectDelay implements capped exponential backoff with jitter. nats.go
// calls this after exhausting the server list, passing the attempt count; the
// returned value is the total sleep, so jitter is added here to de-synchronise
// many connectors reconnecting at once.
func reconnectDelay(attempts int) time.Duration {
	d := reconnectBaseDelay
	for i := 1; i < attempts && d < reconnectMaxDelay; i++ {
		d *= 2
	}
	if d > reconnectMaxDelay {
		d = reconnectMaxDelay
	}
	jitter := time.Duration(rand.Int64N(int64(reconnectBaseDelay / 4)))
	return d + jitter
}

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

	// loggerSet records whether a logger was supplied via WithLogger, so Launch
	// does not clobber an injected logger.
	loggerSet bool

	// credMu guards NatsJwt/NatsSeed against a concurrent credential refresh
	// (SetCredentials) racing the NATS reconnect goroutine.
	credMu sync.RWMutex
}

// logger returns the runtime's logger, falling back to the default if one has
// not been configured yet (e.g. NatsConfig called before Launch).
func (r *Runtime) logger() *slog.Logger {
	if r.Logger != nil {
		return r.Logger
	}
	return slog.Default()
}

// SetCredentials replaces the NATS credentials at runtime. Because NatsOptions'
// JWT/signature callbacks read these fields on every (re)connect handshake, the
// refreshed credentials are used on the next reconnect — the seam a credential
// refresher (e.g. a creds-file watcher) drives.
func (r *Runtime) SetCredentials(jwt, seed string) {
	r.credMu.Lock()
	defer r.credMu.Unlock()
	r.NatsJwt = jwt
	r.NatsSeed = seed
}

func (r *Runtime) currentCredentials() (jwt, seed string) {
	r.credMu.RLock()
	defer r.credMu.RUnlock()
	return r.NatsJwt, r.NatsSeed
}

// NatsOptions builds the NATS connection options from the runtime's (decoded)
// credentials. Credentials are read from the Runtime fields rather than the raw
// environment, and the connection is configured to survive credential and
// connectivity loss instead of permanently aborting:
//
//   - IgnoreAuthErrorAbort keeps the client reconnecting through auth errors, so
//     a rotated/expired credential can be picked up on a later attempt rather
//     than closing the connection for good (which today forces a manual restart).
//   - The JWT/signature callbacks re-read the runtime's credential fields on
//     every (re)connect — the seam a future credential refresh hangs off.
//   - Lifecycle handlers surface disconnects/reconnects/closure that are silent
//     today.
func (r *Runtime) NatsOptions() ([]nats.Option, error) {
	log := r.logger()

	opts := []nats.Option{
		nats.MaxReconnects(-1),
		nats.CustomReconnectDelay(reconnectDelay),
		nats.IgnoreAuthErrorAbort(),
		// Read credentials live on every (re)connect handshake so a refresh via
		// SetCredentials is observed on the next reconnect. An empty JWT yields
		// an anonymous connection (the signature callback is only invoked when a
		// JWT is present).
		nats.UserJWT(
			func() (string, error) {
				jwt, _ := r.currentCredentials()
				return jwt, nil
			},
			func(nonce []byte) ([]byte, error) {
				_, seed := r.currentCredentials()
				kp, err := nkeys.FromSeed([]byte(seed))
				if err != nil {
					return nil, fmt.Errorf("failed to load nats seed: %w", err)
				}
				defer kp.Wipe()
				return kp.Sign(nonce)
			},
		),
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
