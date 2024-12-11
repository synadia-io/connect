package library

import "encoding/json"

type VersionFilter struct {
	RuntimeId string `json:"runtime_id"`
}

type VersionInfo struct {
	VersionId        string `json:"version_id"`
	WorkloadLocation string `json:"workload_location"`
	MetricsEnabled   bool   `json:"metrics_enabled"`
}

type VersionCursor func(v *VersionInfo, hasMore bool) error

func (c *client) ListVersions(filter VersionFilter, cursor VersionCursor, opts ...Opt) error {
	return c.RequestList(c.serviceSubject("VERSION.LIST"), filter, func(resp []byte, hasMore bool) error {
		var v *VersionInfo
		if err := json.Unmarshal(resp, &v); err != nil {
			return err
		}

		return cursor(v, hasMore)
	}, opts...)
}
