package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/synadia-io/connect/model"
)

const HasMoreHeader = "Nats-Has-More"

type ResponseHandler func(resp []byte, hasMore bool) error

type RequestOpts struct {
	Timeout time.Duration
}

func DefaultRequestOpts() *RequestOpts {
	return &RequestOpts{
		Timeout: 5 * time.Second,
	}
}

type Opt func(opts *RequestOpts)

func WithTimeout(timeout time.Duration) Opt {
	return func(opts *RequestOpts) {
		opts.Timeout = timeout
	}
}

type Client interface {
	Account() string

	ListConnectors(filter ConnectorFilter, cursor ConnectorCursor, opts ...Opt) error
	GetConnector(id string, opts ...Opt) (*model.Connector, error)
	CreateConnector(id string, config model.ConnectorConfig, opts ...Opt) (*model.Connector, error)
	UpdateConnector(id string, updates model.ConnectorConfig, opts ...Opt) (*model.Connector, error)
	DeleteConnector(id string, opts ...Opt) (bool, error)

	DeployConnector(id string, opts ...DeployOpt) (*ConnectorDeployResult, error)
	RedeployConnector(id string, opts ...RedeployOpt) (*ConnectorRedeployResult, error)
	UndeployConnector(id string, opts ...UndeployOpt) (*ConnectorUndeployResult, error)

	ListDeployments(filter DeploymentFilter, cursor DeploymentCursor, opts ...Opt) error
	GetDeployment(connectorId string, deploymentId string, opts ...Opt) (*model.Deployment, error)
	PurgeDeployments(filter DeploymentFilter, cursor DeploymentPurgeCursor, opts ...Opt) error

	ListInstances(filter InstanceFilter, cursor InstanceCursor, opts ...Opt) error
	GetInstance(connectorId string, deploymentId string, instanceId string, opts ...Opt) (*model.Instance, error)

	CaptureMetrics(filter CaptureFilter, cursor MetricsCursor) (*nats.Subscription, error)
	CaptureEvents(filter CaptureFilter, cursor EventsCursor) (*nats.Subscription, error)
	CaptureLogs(filter CaptureFilter, cursor LogCursor) (*nats.Subscription, error)
}

func NewClient(nc *nats.Conn, trace bool) (Client, error) {
	ui, err := getUserInfo(nc)
	if err != nil {
		return nil, err
	}

	return NewClientForAccount(nc, ui.Account, trace), nil
}

func NewClientForAccount(nc *nats.Conn, account string, trace bool) Client {
	return &client{
		nc:      nc,
		account: account,
		trace:   trace,
	}
}

type client struct {
	nc      *nats.Conn
	trace   bool
	account string
}

func (c *client) serviceSubject(subject string) string {
	return fmt.Sprintf("$CONSVC.%s.%s", c.account, subject)
}

func (c *client) Account() string {
	return c.account
}

func (c *client) Request(subject string, payload any, opts ...Opt) ([]byte, error) {
	options := DefaultRequestOpts()
	for _, opt := range opts {
		opt(options)
	}

	// -- encode the Request
	req, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal Request: %v", err)
	}

	if c.trace {
		fmt.Println(">>> ", subject, " [", string(req), "]")
	}

	resp, err := c.nc.Request(subject, req, options.Timeout)
	if err != nil {
		return nil, fmt.Errorf("unable to get response: %v", err)
	}

	serviceErr := resp.Header.Get("Nats-Service-Error")
	serviceErrCode := resp.Header.Get("Nats-Service-Error-Code")
	if serviceErr != "" {
		return nil, fmt.Errorf("%s (%s)", serviceErr, serviceErrCode)
	}

	return resp.Data, nil
}

func (c *client) RequestList(subject string, payload any, h ResponseHandler, opts ...Opt) error {
	options := DefaultRequestOpts()
	for _, opt := range opts {
		opt(options)
	}

	inb := nats.NewInbox()
	if c.trace {
		fmt.Println("??? ", inb)
	}

	sub, err := c.nc.SubscribeSync(inb)
	if err != nil {
		return err
	}

	// -- encode the Request
	req, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("unable to marshal Request: %v", err)
	}

	if c.trace {
		fmt.Println(">>> ", subject, " [", string(req), "]")
	}

	if err := c.nc.PublishRequest(subject, inb, req); err != nil {
		return fmt.Errorf("unable to publish Request: %v", err)
	}
	if err := c.nc.Flush(); err != nil {
		return fmt.Errorf("unable to flush: %v", err)
	}

	for {
		msg, err := sub.NextMsg(options.Timeout)
		if err != nil {
			return fmt.Errorf("unable to get response: %v", err)
		}

		serviceErr := msg.Header.Get("Nats-Service-Error")
		serviceErrCode := msg.Header.Get("Nats-Service-Error-Code")
		if serviceErr != "" {
			return fmt.Errorf("%s (%s)", serviceErr, serviceErrCode)
		}

		hasMore := msg.Header.Get(HasMoreHeader) == "true"
		data := msg.Data
		if bytes.Equal(data, []byte("null")) {
			data = nil
		}
		if err := h(data, hasMore); err != nil {
			return err
		}

		if !hasMore {
			break
		}
	}

	return nil
}

func (c *client) Close() {
	c.nc.Close()
}
