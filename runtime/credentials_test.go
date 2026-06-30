package runtime_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/nats-io/nkeys"
	"github.com/synadia-io/connect/runtime"

	. "github.com/onsi/gomega"
)

// credsContent builds a standard decorated NATS creds file (JWT block first,
// seed block second) that nats.UserCredentials / nkeys.ParseDecorated* read.
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

// When NEX_WORKLOAD_NATS_CREDS_FILE is set, the connection uses
// nats.UserCredentials, which re-reads the creds file on every handshake. So a
// re-minted credential (the control plane rewriting the file) is picked up on
// the next reconnect with no refresh machinery on our side.
func TestNatsOptionsUsesCredentialsFileWithLiveReload(t *testing.T) {
	g := NewWithT(t)

	seed1, pub1 := newUserSeed(g)
	seed2, pub2 := newUserSeed(g)
	g.Expect(pub1).ToNot(Equal(pub2))

	path := filepath.Join(t.TempDir(), "nats.creds")
	g.Expect(os.WriteFile(path, []byte(credsContent("jwt-1", seed1)), 0o600)).To(Succeed())
	t.Setenv(runtime.NatsCredsFileVar, path)

	rt := runtime.NewRuntime(runtime.WithNatsUrl("nats://example:4222"))
	opts, err := rt.NatsOptions()
	g.Expect(err).ToNot(HaveOccurred())
	o := applyOptions(g, opts)

	nonce := []byte("server-nonce")
	verifiesWith := func(pub string, sig []byte) error {
		vk, err := nkeys.FromPublicKey(pub)
		g.Expect(err).ToNot(HaveOccurred())
		return vk.Verify(nonce, sig)
	}

	// Initial creds come from the file.
	jwt, err := o.UserJWT()
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(jwt).To(Equal("jwt-1"))
	sig, err := o.SignatureCB(nonce)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(verifiesWith(pub1, sig)).To(Succeed())

	// Control plane re-mints: rewrite the file in place.
	g.Expect(os.WriteFile(path, []byte(credsContent("jwt-2", seed2)), 0o600)).To(Succeed())

	// The same callbacks (invoked on each reconnect handshake) now read the new
	// file — no SetCredentials/watcher needed.
	jwt, err = o.UserJWT()
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(jwt).To(Equal("jwt-2"), "next handshake must read the refreshed creds file")
	sig, err = o.SignatureCB(nonce)
	g.Expect(err).ToNot(HaveOccurred())
	g.Expect(verifiesWith(pub2, sig)).To(Succeed(), "next handshake must sign with the refreshed seed")
	g.Expect(verifiesWith(pub1, sig)).ToNot(Succeed(), "stale seed must no longer be used")
}
