package client

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"github.com/synadia-io/connect/model"
)

type MetricsCursor func(item model.InstanceMetric)

func (c *client) CaptureMetrics(filter CaptureFilter, cursor MetricsCursor) (*nats.Subscription, error) {
	return c.capture("METRICS", filter, func(msg *nats.Msg) {
		// -- parse the message
		var item model.InstanceMetric
		if err := json.Unmarshal(msg.Data, &item); err != nil {
			return
		}

		cursor(item)
	})
}
