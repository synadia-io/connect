package library

import (
	"encoding/json"
	"github.com/synadia-io/connect/model"
)

type RuntimeInfo struct {
	Id              string       `json:"id"`
	Label           string       `json:"label"`
	LatestVersionId string       `json:"latest_version_id"`
	Author          model.Author `json:"author"`
	Description     string       `json:"description"`
}

type RuntimeCursor func(rt *RuntimeInfo, hasMore bool) error

func (c *client) ListRuntimes(cursor RuntimeCursor, opts ...Opt) error {
	return c.RequestList(c.serviceSubject("RUNTIME.LIST"), nil, func(b []byte, hasMore bool) error {
		var rt *RuntimeInfo
		if err := json.Unmarshal(b, &rt); err != nil {
			return err
		}

		return cursor(rt, hasMore)
	}, opts...)
}
