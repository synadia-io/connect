package cli

import (
	"github.com/nats-io/jsm.go/natscontext"
	"time"
)

var DefaultOptions *Options

// Options configure the CLI
type Options struct {
	// Config is a nats configuration context
	Config *natscontext.Context

	// Servers is the list of servers to connect to
	Servers string
	// Username is the username or token to connect with
	Username string
	// Password is the password to connect with
	Password string
	// ConnectionName is the name to use for the underlying NATS connection
	ConnectionName string
	// Creds is nats credentials to authenticate with
	Creds string
	// Nkey is the file holding a nkey to connect with
	Nkey string
	// TlsCert is the TLS Public Certificate
	TlsCert string
	// TlsKey is the TLS Private Key
	TlsKey string
	// TlsCA is the certificate authority to verify the connection with
	TlsCA string
	// TlsFirst configures the TLSHandshakeFirst behavior in nats.go
	TlsFirst bool
	// CfgCtx is the context name to use
	CfgCtx string
	// Timeout is how long to wait for operations
	Timeout time.Duration
	Trace   bool

	// LogLevel is the log level to use
	LogLevel string
}
