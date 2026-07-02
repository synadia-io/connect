package runtime_test

import (
	"context"
	"encoding/base64"
	"log/slog"
	"testing"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	"github.com/synadia-io/connect/model"
	"github.com/synadia-io/connect/runtime"

	. "github.com/onsi/gomega"
)

// minimalCfg is a base64-encoded YAML connector spec sufficient for Launch to
// decode and unmarshal into model.Steps.
func minimalCfg() string {
	return base64.StdEncoding.EncodeToString([]byte("source:\n    type: generate\n"))
}

// applyOptions folds a []nats.Option into a fresh nats.Options so the resulting
// configuration can be asserted without opening a real connection.
func applyOptions(g Gomega, opts []nats.Option) *nats.Options {
	o := &nats.Options{}
	for _, opt := range opts {
		g.Expect(opt(o)).To(Succeed())
	}
	return o
}

// The connector must not give up on a stale credential: keep reconnecting
// through auth errors and past the default attempt cap, sourcing creds from the
// DECODED runtime fields (not the raw base64 env value). Reconnect backoff is
// left to the nats.go client default (fixed wait + jitter).
func TestNatsOptionsUsesDecodedCredentialsAndResilience(t *testing.T) {
	g := NewWithT(t)
	t.Setenv(runtime.NatsCredsFileVar, "") // force the env-creds path

	kp, err := nkeys.CreateUser()
	g.Expect(err).ToNot(HaveOccurred())
	seed, err := kp.Seed()
	g.Expect(err).ToNot(HaveOccurred())
	pub, err := kp.PublicKey()
	g.Expect(err).ToNot(HaveOccurred())

	const decodedJWT = "this-is-the-decoded-jwt"

	rt := runtime.NewRuntime(
		runtime.WithNatsUrl("nats://example:4222"),
		runtime.WithNatsJwt(decodedJWT),
		runtime.WithNatsSeed(string(seed)),
	)

	opts, err := rt.NatsOptions()
	g.Expect(err).ToNot(HaveOccurred())
	o := applyOptions(g, opts)

	g.Expect(o.IgnoreAuthErrorAbort).To(BeTrue(), "a stale cred must not permanently close the connection")
	g.Expect(o.MaxReconnect).To(Equal(-1), "reconnect must be unbounded (the default caps at 60)")

	// Credentials must come from the DECODED field, not the raw base64 env var.
	g.Expect(o.UserJWT).ToNot(BeNil())
	jwt, err := o.UserJWT()
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(jwt).To(Equal(decodedJWT))

	g.Expect(o.SignatureCB).ToNot(BeNil())
	nonce := []byte("server-nonce")
	sig, err := o.SignatureCB(nonce)
	g.Expect(err).ToNot(HaveOccurred())
	vk, err := nkeys.FromPublicKey(pub)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(vk.Verify(nonce, sig)).To(Succeed())
}

// Launch unconditionally overwrote r.Logger with slog.Default(), discarding a
// logger supplied via WithLogger.
func TestLaunchPreservesInjectedLogger(t *testing.T) {
	g := NewWithT(t)

	injected := slog.New(slog.NewTextHandler(nil, &slog.HandlerOptions{Level: slog.LevelError}))

	rt := runtime.NewRuntime(runtime.WithLogger(injected))

	var seen *slog.Logger
	err := rt.Launch(context.Background(), func(_ context.Context, r *runtime.Runtime, _ model.Steps) error {
		seen = r.Logger
		return nil
	}, minimalCfg())
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(seen).To(BeIdenticalTo(injected), "an injected logger must survive Launch")
}

// When no logger is injected, Launch must build one honoring the configured
// LogLevel (parsed from CONNECT_LOG_LEVEL) rather than the default Info handler.
func TestLaunchHonorsConfiguredLogLevelWhenNoLoggerInjected(t *testing.T) {
	g := NewWithT(t)

	rt := runtime.NewRuntime(runtime.WithLogLevel(slog.LevelWarn))

	var seen *slog.Logger
	err := rt.Launch(context.Background(), func(_ context.Context, r *runtime.Runtime, _ model.Steps) error {
		seen = r.Logger
		return nil
	}, minimalCfg())
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(seen).ToNot(BeNil())
	ctx := context.Background()
	g.Expect(seen.Enabled(ctx, slog.LevelWarn)).To(BeTrue(), "warn must be enabled at LevelWarn")
	g.Expect(seen.Enabled(ctx, slog.LevelInfo)).To(BeFalse(), "info must be suppressed at LevelWarn")
}

// With CONNECT_LOG_LEVEL unset, the effective log level must default to Info, not
// Debug: Launch now builds its handler at r.LogLevel, so a Debug default would
// silently flip every connector into debug logging. This exercises the real
// deploy path (FromEnv -> NewRuntime -> Launch).
func TestDefaultLogLevelIsInfoWhenEnvUnset(t *testing.T) {
	g := NewWithT(t)
	t.Setenv(runtime.LogLevelEnvVar, "") // env unset -> default level

	rt, err := runtime.FromEnv()
	g.Expect(err).ToNot(HaveOccurred())

	var seen *slog.Logger
	err = rt.Launch(context.Background(), func(_ context.Context, r *runtime.Runtime, _ model.Steps) error {
		seen = r.Logger
		return nil
	}, minimalCfg())
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(seen).ToNot(BeNil())

	ctx := context.Background()
	g.Expect(seen.Enabled(ctx, slog.LevelInfo)).To(BeTrue(), "info must be enabled by default")
	g.Expect(seen.Enabled(ctx, slog.LevelDebug)).To(BeFalse(), "debug must be suppressed by default (no CONNECT_LOG_LEVEL)")
}
