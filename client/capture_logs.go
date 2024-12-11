package client

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/synadia-io/connect/model"
)

type LogCursor func(item model.InstanceLog)

func (c *client) CaptureLogs(filter CaptureFilter, cursor LogCursor) (*nats.Subscription, error) {
	return c.capture("LOGS", filter, func(msg *nats.Msg) {
		// -- parse the message
		var item model.InstanceLog
		if err := json.Unmarshal(msg.Data, &item); err != nil {
			return
		}

		cursor(item)
	})
}
