package library

import (
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

type Client interface {
	ListRuntimes(cursor RuntimeCursor, opts ...Opt) error
	GetRuntime(runtimeId string, opts ...Opt) (*model.Runtime, error)

	ListVersions(filter VersionFilter, cursor VersionCursor, opts ...Opt) error
	GetLatestVersion(runtimeId string, opts ...Opt) (*model.Version, error)
	GetVersion(runtimeId string, versionId string, opts ...Opt) (*model.Version, error)

	ListComponents(filter ComponentFilter, cursor ComponentCursor, opts ...Opt) error
	GetComponent(runtimeId string, versionId string, kind model.ComponentKind, name string, opts ...Opt) (*model.Component, error)

	Close()
}

func NewClient(nc *nats.Conn, trace bool) Client {
	return &client{
		nc:    nc,
		trace: trace,
	}
}

type client struct {
	nc      *nats.Conn
	trace   bool
	account string
}

func (c *client) serviceSubject(subject string) string {
	return fmt.Sprintf("$CONLIB.%s", subject)
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
		if err := h(msg.Data, hasMore); err != nil {
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
