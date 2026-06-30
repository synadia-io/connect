package runtime

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/nats-io/nkeys"
)

// defaultCredsWatchInterval is how often the credentials file is polled for a
// refresh when no interval is supplied.
const defaultCredsWatchInterval = 30 * time.Second

// LoadCredentialsFile reads a standard decorated NATS creds file (a JWT block
// and a user NKEY seed block) and returns the JWT and seed.
func LoadCredentialsFile(path string) (jwt, seed string, err error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", "", fmt.Errorf("failed to read credentials file: %w", err)
	}

	jwt, err = nkeys.ParseDecoratedJWT(contents)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse jwt from credentials file: %w", err)
	}

	kp, err := nkeys.ParseDecoratedNKey(contents)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse seed from credentials file: %w", err)
	}
	defer kp.Wipe()

	seedBytes, err := kp.Seed()
	if err != nil {
		return "", "", fmt.Errorf("failed to extract seed from credentials file: %w", err)
	}

	return jwt, string(seedBytes), nil
}

// WatchCredentialsFile loads the credentials at path and applies them to the
// runtime via SetCredentials, then polls for changes and re-applies on every
// change until ctx is cancelled. This is the refresh driver: when the control
// plane re-mints a credential and rewrites the file, the next reconnect
// handshake picks it up (NatsOptions reads the runtime's credentials live).
//
// A failed read/parse is logged and skipped — the existing credentials are kept
// rather than dropped, so a transient bad write never kills the connector.
func (r *Runtime) WatchCredentialsFile(ctx context.Context, path string, interval time.Duration) error {
	if interval <= 0 {
		interval = defaultCredsWatchInterval
	}
	log := r.logger()

	var lastJWT, lastSeed string
	apply := func() {
		jwt, seed, err := LoadCredentialsFile(path)
		if err != nil {
			log.Warn("failed to load credentials file", slog.String("path", path), slog.Any("err", err))
			return
		}
		if jwt == lastJWT && seed == lastSeed {
			return
		}
		lastJWT, lastSeed = jwt, seed
		r.SetCredentials(jwt, seed)
		log.Info("nats credentials refreshed from file", slog.String("path", path))
	}

	apply()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			apply()
		}
	}
}
