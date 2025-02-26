package client

import (
	"time"

	"github.com/nats-io/nats.go"
	"github.com/synadia-io/connect/v2/model"
)

type Client interface {
	Account() string

	ConnectorClient
	LibraryClient

	Close()
}

type ConnectorClient interface {
	ListConnectors(timeout time.Duration) ([]model.ConnectorSummary, error)
	GetConnector(id string, timeout time.Duration) (*model.Connector, error)
	GetConnectorStatus(id string, timeout time.Duration) (*model.ConnectorStatus, error)
	CreateConnector(id, description, runtimeId string, steps model.Steps, timeout time.Duration) (*model.Connector, error)
	PatchConnector(id string, patch string, timeout time.Duration) (*model.Connector, error)
	DeleteConnector(id string, timeout time.Duration) error

	ListConnectorInstances(id string, timeout time.Duration) ([]model.Instance, error)
	StartConnector(id string, startOpts *model.ConnectorStartOptions, timeout time.Duration) ([]model.Instance, error)
	StopConnector(id string, timeout time.Duration) ([]model.Instance, error)
}

type LibraryClient interface {
	ListRuntimes(timeout time.Duration) ([]model.RuntimeSummary, error)
	GetRuntime(id string, timeout time.Duration) (*model.Runtime, error)

	SearchComponents(filter *model.ComponentSearchFilter, timeout time.Duration) ([]model.ComponentSummary, error)
	GetComponent(runtimeId string, kind model.ComponentKind, id string, timeout time.Duration) (*model.Component, error)
}

func NewClient(nc *nats.Conn, trace bool) (Client, error) {
	t, err := NewTransport(nc, trace)
	if err != nil {
		return nil, err
	}

	return &client{
		t:               t,
		connectorClient: connectorClient{t: t},
		libraryClient:   libraryClient{t: t},
	}, nil
}

func NewClientForAccount(nc *nats.Conn, account string, trace bool) Client {
	t := NewTransportForAccount(nc, account, trace)

	return &client{
		t:               t,
		connectorClient: connectorClient{t: t},
		libraryClient:   libraryClient{t: t},
	}
}

type client struct {
	t *Transport
	connectorClient
	libraryClient
}

func (c *client) Account() string {
	return c.t.Account()
}

func (c *client) Close() {
	c.t.Close()
}
