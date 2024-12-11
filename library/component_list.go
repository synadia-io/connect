package library

import (
	"encoding/json"
	"github.com/synadia-io/connect/model"
)

type ComponentFilter struct {
	RuntimeId string                `json:"runtime_id"`
	VersionId string                `json:"version_id"`
	Status    model.ComponentStatus `json:"status"`
	Kind      model.ComponentKind   `json:"kind"`
}

type ComponentInfo struct {
	RuntimeId string                `json:"runtime_id"`
	VersionId string                `json:"version_id"`
	Name      string                `json:"name"`
	Label     string                `json:"label"`
	Kind      model.ComponentKind   `json:"kind"`
	Status    model.ComponentStatus `json:"status"`
}

type ComponentCursor func(c *ComponentInfo, hasMore bool) error

func (c *client) ListComponents(filter ComponentFilter, cursor ComponentCursor, opts ...Opt) error {
	return c.RequestList(c.serviceSubject("COMPONENT.LIST"), filter, func(resp []byte, hasMore bool) error {
		var c *ComponentInfo
		if err := json.Unmarshal(resp, &c); err != nil {
			return err
		}

		return cursor(c, hasMore)
	}, opts...)
}
