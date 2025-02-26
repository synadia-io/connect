package test

import (
	"github.com/nats-io/jwt/v2"
	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nkeys"
	. "github.com/onsi/gomega"
)

func Account(name string) Acc {
	// Account with signing key
	acc, _ := nkeys.CreateAccount()
	accPk, _ := acc.PublicKey()
	accClaim := jwt.NewAccountClaims(accPk)
	accSk, _ := nkeys.CreateAccount()
	accSkPk, _ := accSk.PublicKey()
	accSkSeed, _ := accSk.Seed()
	accClaim.Name = name
	accClaim.SigningKeys.Add(accSkPk)
	accClaim.Limits.JetStreamLimits.MemoryStorage = -1
	accClaim.Limits.JetStreamLimits.DiskStorage = -1

	return Acc{
		Id:          accPk,
		KeyPair:     acc,
		SigningSeed: accSkSeed,
		Claims:      accClaim,
	}
}

type Acc struct {
	Id          string
	KeyPair     nkeys.KeyPair
	SigningSeed []byte
	Claims      *jwt.AccountClaims
}

func (a *Acc) Creds() (string, []byte) {
	user, _ := nkeys.CreateUser()
	userPk, _ := user.PublicKey()
	userSeed, _ := user.Seed()
	userClaim := jwt.NewUserClaims(userPk)
	userJwt, err := userClaim.Encode(a.KeyPair)
	Expect(err).ToNot(HaveOccurred())

	return userJwt, userSeed
}

func (a *Acc) Connect(srv *server.Server) *nats.Conn {
	userJwt, userSeed := a.Creds()

	natsUrl := srv.ClientURL()
	nc, err := nats.Connect(natsUrl, nats.UserJWTAndSeed(userJwt, string(userSeed)))
	Expect(err).ToNot(HaveOccurred())

	return nc
}
