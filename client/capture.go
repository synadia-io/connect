package client

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

type CaptureFilter struct {
	ConnectorId  string `json:"connector_id"`
	DeploymentId string `json:"deployment_id"`
	InstanceId   string `json:"instance_id"`
}

type captureCursor func(msg *nats.Msg)

func (c *client) capture(kind string, filter CaptureFilter, cursor captureCursor) (*nats.Subscription, error) {
	acc := c.Account()

	cid := "*"
	if filter.ConnectorId != "" {
		cid = filter.ConnectorId
	}

	did := "*"
	if filter.DeploymentId != "" {
		did = filter.DeploymentId
	}

	iid := "*"
	if filter.InstanceId != "" {
		iid = filter.InstanceId
	}

	subject := fmt.Sprintf("$CONNECT.%s.%s.CONNECTOR.%s.DEPLOYMENT.%s.INSTANCE.%s", acc, kind, cid, did, iid)

	// subscribe
	s, err := c.nc.Subscribe(subject, func(msg *nats.Msg) {
		if msg == nil || msg.Data == nil {
			return
		}

		cursor(msg)
	})
	if err != nil {
		return nil, err
	}

	return s, nil
}
