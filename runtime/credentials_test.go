package runtime_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nats-io/nkeys"
	"github.com/synadia-io/connect/model"
	"github.com/synadia-io/connect/runtime"

	. "github.com/onsi/gomega"
)

// credsContent builds a standard decorated NATS creds file (JWT block first,
// seed block second) that nkeys.ParseDecorated* can read.
func credsContent(jwt, seed string) string {
	return fmt.Sprintf(
		"-----BEGIN NATS USER JWT-----\n%s\n------END NATS USER JWT------\n\n"+
			"-----BEGIN USER NKEY SEED-----\n%s\n------END USER NKEY SEED------\n",
		jwt, seed,
	)
}

func newUserSeed(g Gomega) (seed string, pub string) {
	kp, err := nkeys.CreateUser()
	g.Expect(err).ToNot(HaveOccurred())
	s, err := kp.Seed()
	g.Expect(err).ToNot(HaveOccurred())
	p, err := kp.PublicKey()
	g.Expect(err).ToNot(HaveOccurred())
	return string(s), p
}

func TestLoadCredentialsFileParsesJwtAndSeed(t *testing.T) {
	g := NewWithT(t)

	seed, _ := newUserSeed(g)
	path := filepath.Join(t.TempDir(), "nats.creds")
	g.Expect(os.WriteFile(path, []byte(credsContent("jwt-1", seed)), 0o600)).To(Succeed())

	jwt, gotSeed, err := runtime.LoadCredentialsFile(path)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(jwt).To(Equal("jwt-1"))
	g.Expect(gotSeed).To(Equal(seed))
}

func TestLoadCredentialsFileRejectsGarbage(t *testing.T) {
	g := NewWithT(t)

	path := filepath.Join(t.TempDir(), "bad.creds")
	g.Expect(os.WriteFile(path, []byte("not a creds file"), 0o600)).To(Succeed())

	_, _, err := runtime.LoadCredentialsFile(path)
	g.Expect(err).To(HaveOccurred())
}

// The watcher is part B: it loads creds at startup and, when the file changes
// (the control plane re-mints before expiry), calls SetCredentials so the next
// reconnect handshake picks up the fresh credential.
func TestWatchCredentialsFileAppliesInitialAndUpdatedCredentials(t *testing.T) {
	g := NewWithT(t)

	seed1, _ := newUserSeed(g)
	seed2, _ := newUserSeed(g)

	path := filepath.Join(t.TempDir(), "nats.creds")
	g.Expect(os.WriteFile(path, []byte(credsContent("jwt-1", seed1)), 0o600)).To(Succeed())

	rt := runtime.NewRuntime()
	// Bind the option callbacks to rt; they read the runtime's live credentials.
	opts, err := rt.NatsOptions()
	g.Expect(err).ToNot(HaveOccurred())
	o := applyOptions(g, opts)
	currentJWT := func() string {
		j, _ := o.UserJWT()
		return j
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() { _ = rt.WatchCredentialsFile(ctx, path, 5*time.Millisecond) }()

	// Initial load is applied.
	g.Eventually(currentJWT).Should(Equal("jwt-1"))

	// A re-minted creds file is observed on the next poll.
	g.Expect(os.WriteFile(path, []byte(credsContent("jwt-2", seed2)), 0o600)).To(Succeed())
	g.Eventually(currentJWT).Should(Equal("jwt-2"))
}

// When NEX_WORKLOAD_NATS_CREDS_FILE is set, Launch must adopt the file's
// credentials before running the workload (and keep watching for refreshes).
func TestLaunchAdoptsCredentialsFromCredsFileEnv(t *testing.T) {
	g := NewWithT(t)

	seed, _ := newUserSeed(g)
	path := filepath.Join(t.TempDir(), "nats.creds")
	g.Expect(os.WriteFile(path, []byte(credsContent("file-jwt", seed)), 0o600)).To(Succeed())
	t.Setenv(runtime.NatsCredsFileVar, path)

	rt := runtime.NewRuntime()
	var jwt string
	err := rt.Launch(context.Background(), func(_ context.Context, r *runtime.Runtime, _ model.Steps) error {
		opts, err := r.NatsOptions()
		g.Expect(err).ToNot(HaveOccurred())
		o := applyOptions(g, opts)
		jwt, _ = o.UserJWT()
		return nil
	}, minimalCfg())
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(jwt).To(Equal("file-jwt"), "Launch must adopt credentials from the creds file")
}
