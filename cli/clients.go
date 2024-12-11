package cli

import (
	"github.com/choria-io/fisk"
	"github.com/nats-io/nats.go"
	"github.com/synadia-io/connect/client"
	"github.com/synadia-io/connect/library"
	"sync"
)

var (
	nc   *nats.Conn
	lock sync.Mutex
)

func natsClient() *nats.Conn {
	lock.Lock()
	defer lock.Unlock()

	if nc == nil {
		var err error
		nc, err = opts().Config.Connect()
		fisk.FatalIfError(err, "Could not connect to NATS")
	}

	return nc
}

func libraryClient() library.Client {
	return library.NewClient(natsClient(), opts().Trace)
}

func controlClient() client.Client {
	cc, err := client.NewClient(natsClient(), opts().Trace)
	fisk.FatalIfError(err, "Could not create control client")
	return cc
}
