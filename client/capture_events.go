package client

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/synadia-io/connect/model"
)

type EventsCursor func(item model.InstanceEvent)

func (c *client) CaptureEvents(filter CaptureFilter, cursor EventsCursor) (*nats.Subscription, error) {
	return c.capture("EVENTS", filter, func(msg *nats.Msg) {
		// -- parse the message
		var item model.InstanceEvent
		if err := json.Unmarshal(msg.Data, &item); err != nil {
			return
		}

		item.Type = model.InstanceEventType(msg.Header.Get(model.EventTypeHeaderName))

		cursor(item)
	})
}
