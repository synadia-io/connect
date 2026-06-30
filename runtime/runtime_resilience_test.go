package runtime_test

import (
	"context"
	"encoding/base64"
	"log/slog"
	"testing"
	"time"

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

// The connector currently aborts permanently when a credential goes stale,
// because NatsConfig sets no resilience options and re-reads the RAW base64 env
// JWT instead of the decoded value. The fix: NatsOptions must source creds from
// the (decoded) Runtime fields and keep reconnecting through auth errors.
func TestNatsOptionsUsesDecodedCredentialsAndResilience(t *testing.T) {
	g := NewWithT(t)

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

	// Resilience: an auth error must NOT permanently abort reconnect, and the
	// client must keep trying rather than give up after the default budget.
	g.Expect(o.IgnoreAuthErrorAbort).To(BeTrue(), "IgnoreAuthErrorAbort must be set or a stale cred permanently closes the connection")
	g.Expect(o.MaxReconnect).To(Equal(-1), "reconnect must be unbounded so the connector survives transient outages")
	g.Expect(o.CustomReconnectDelayCB).ToNot(BeNil(), "an exponential backoff handler avoids a hot loop and a thundering herd")

	// Credentials must come from the DECODED Runtime field, not the raw base64
	// env var. The JWT callback re-reads the field on each (re)connect, which is
	// the seam a future credential refresh hangs off.
	g.Expect(o.UserJWT).ToNot(BeNil())
	jwt, err := o.UserJWT()
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(jwt).To(Equal(decodedJWT))

	// The signature callback must sign with the provided seed.
	g.Expect(o.SignatureCB).ToNot(BeNil())
	nonce := []byte("server-nonce")
	sig, err := o.SignatureCB(nonce)
	g.Expect(err).ToNot(HaveOccurred())
	vk, err := nkeys.FromPublicKey(pub)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(vk.Verify(nonce, sig)).To(Succeed())
}

// The reconnect delay must grow exponentially (with jitter) and be capped, so
// the connector backs off instead of hot-looping and many connectors do not
// reconnect in lockstep.
func TestNatsReconnectDelayIsExponentialAndCapped(t *testing.T) {
	g := NewWithT(t)

	rt := runtime.NewRuntime(runtime.WithNatsUrl("nats://example:4222"))
	opts, err := rt.NatsOptions()
	g.Expect(err).ToNot(HaveOccurred())
	o := applyOptions(g, opts)

	g.Expect(o.CustomReconnectDelayCB).ToNot(BeNil())

	d1 := o.CustomReconnectDelayCB(1)
	d2 := o.CustomReconnectDelayCB(2)
	d3 := o.CustomReconnectDelayCB(3)

	g.Expect(d1).To(BeNumerically(">", 0))
	g.Expect(d2).To(BeNumerically(">", d1), "delay must grow with attempts")
	g.Expect(d3).To(BeNumerically(">", d2), "delay must grow with attempts")

	// A large attempt count must be capped, not unbounded.
	dCapped := o.CustomReconnectDelayCB(1000)
	g.Expect(dCapped).To(BeNumerically("<=", 31*time.Second), "delay must be capped")
}

// A credential refresh (e.g. the control plane re-mints before expiry) must be
// observed on the NEXT (re)connect: nats.go invokes UserJWT/SignatureCB on every
// handshake, so the callbacks must read the runtime's current credentials rather
// than a value captured at option-build time.
func TestRefreshedCredentialsAreUsedOnNextHandshake(t *testing.T) {
	g := NewWithT(t)

	newUser := func() (seed []byte, pub string) {
		kp, err := nkeys.CreateUser()
		g.Expect(err).ToNot(HaveOccurred())
		s, err := kp.Seed()
		g.Expect(err).ToNot(HaveOccurred())
		p, err := kp.PublicKey()
		g.Expect(err).ToNot(HaveOccurred())
		return s, p
	}

	seed1, pub1 := newUser()
	seed2, pub2 := newUser()
	g.Expect(pub1).ToNot(Equal(pub2))

	rt := runtime.NewRuntime(
		runtime.WithNatsUrl("nats://example:4222"),
		runtime.WithNatsJwt("jwt-1"),
		runtime.WithNatsSeed(string(seed1)),
	)
	opts, err := rt.NatsOptions()
	g.Expect(err).ToNot(HaveOccurred())
	o := applyOptions(g, opts)

	nonce := []byte("server-nonce")
	verify := func(pub string) func([]byte) error {
		vk, err := nkeys.FromPublicKey(pub)
		g.Expect(err).ToNot(HaveOccurred())
		return func(sig []byte) error { return vk.Verify(nonce, sig) }
	}

	// Before refresh: original creds.
	jwt, err := o.UserJWT()
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(jwt).To(Equal("jwt-1"))
	sig, err := o.SignatureCB(nonce)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(verify(pub1)(sig)).To(Succeed())

	// Control plane drops fresh creds onto the runtime.
	rt.SetCredentials("jwt-2", string(seed2))

	// The SAME option callbacks (the ones nats re-invokes on reconnect) now
	// yield the refreshed credentials.
	jwt, err = o.UserJWT()
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(jwt).To(Equal("jwt-2"), "next handshake must use the refreshed jwt")
	sig, err = o.SignatureCB(nonce)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(verify(pub2)(sig)).To(Succeed(), "next handshake must sign with the refreshed seed")
	g.Expect(verify(pub1)(sig)).ToNot(Succeed(), "stale seed must no longer be used")
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
