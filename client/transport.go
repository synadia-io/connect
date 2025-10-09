package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
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

func NewTransport(nc *nats.Conn, trace bool) (*Transport, error) {
	account, err := getUserInfo(nc)
	if err != nil {
		return nil, fmt.Errorf("could not get user info: %s", err)
	}

	return NewTransportForAccount(nc, account.Account, trace), nil
}

func NewTransportForAccount(nc *nats.Conn, account string, trace bool) *Transport {
	return &Transport{
		nc:      nc,
		trace:   trace,
		account: account,
	}
}

type Transport struct {
	nc      *nats.Conn
	trace   bool
	account string
}

func (t *Transport) Account() string {
	return t.account
}

func (t *Transport) Request(subject string, payload any, opts ...Opt) ([]byte, error) {
	options := DefaultRequestOpts()
	for _, opt := range opts {
		opt(options)
	}

	// -- encode the Request
	req, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("unable to marshal Request: %v", err)
	}

	if t.trace {
		fmt.Println(">>> ", subject, " [", string(req), "]")
	}

	resp, err := t.nc.Request(subject, req, options.Timeout)
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

func (t *Transport) RequestJson(subject string, payload any, target any, opts ...Opt) (bool, error) {
	b, err := t.Request(subject, payload, opts...)
	if err != nil {
		return false, err
	}

	if b == nil || len(b) == 0 {
		return false, nil
	}

	if err := json.Unmarshal(b, target); err != nil {
		return false, fmt.Errorf("could not parse response: %s", err)
	}

	return true, nil
}

func (t *Transport) RequestList(subject string, payload any, h ResponseHandler, opts ...Opt) error {
	options := DefaultRequestOpts()
	for _, opt := range opts {
		opt(options)
	}

	inb := nats.NewInbox()
	if t.trace {
		fmt.Println("??? ", inb)
	}

	sub, err := t.nc.SubscribeSync(inb)
	if err != nil {
		return err
	}

	// -- encode the Request
	req, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("unable to marshal Request: %v", err)
	}

	if t.trace {
		fmt.Println(">>> ", subject, " [", string(req), "]")
	}

	if err := t.nc.PublishRequest(subject, inb, req); err != nil {
		return fmt.Errorf("unable to publish Request: %v", err)
	}
	if err := t.nc.Flush(); err != nil {
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

func (t *Transport) Close() {
	t.nc.Close()
}

func getUserInfo(nc *nats.Conn) (*server.UserInfo, error) {
	resp, err := nc.Request("$SYS.REQ.USER.INFO", nil, time.Second)
	if err != nil {
		return nil, fmt.Errorf("could not get user info: %s", err)
	}

	var res = struct {
		Data   *server.UserInfo  `json:"data"`
		Server server.ServerInfo `json:"server"`
		Error  *server.ApiError  `json:"error"`
	}{}

	if err = json.Unmarshal(resp.Data, &res); err != nil {
		return nil, fmt.Errorf("could not parse user info: %s", err)
	}

	if res.Error != nil {
		return nil, fmt.Errorf("could not get user info: %s", res.Error)
	}

	return res.Data, nil
}
