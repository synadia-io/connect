package library

import (
	"encoding/json"
	"github.com/synadia-io/connect/model"
)

type (
	getVersionRequest struct {
		RuntimeId string `json:"runtime_id"`
		VersionId string `json:"version_id"`
	}

	getVersionResponse struct {
		Found    bool           `json:"found"`
		Revision uint64         `json:"revision"`
		Version  *model.Version `json:"version,omitempty"`
	}
)

func (c *client) GetLatestVersion(runtimeId string, opts ...Opt) (*model.Version, error) {
	return c.GetVersion(runtimeId, "latest", opts...)
}

func (c *client) GetVersion(runtimeId string, versionId string, opts ...Opt) (*model.Version, error) {
	req := getVersionRequest{
		RuntimeId: runtimeId,
		VersionId: versionId,
	}

	b, err := c.Request(c.serviceSubject("VERSION.GET"), req, opts...)
	if err != nil {
		return nil, err
	}

	var resp getVersionResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return nil, err
	}

	if !resp.Found {
		return nil, nil
	}

	return resp.Version, nil
}
