package model

type Version struct {
	RuntimeId string    `json:"runtime_id"`
	VersionId string    `json:"version_id"`
	Workload  *Workload `json:"workload,omitempty"`
}

type Workload struct {
	Location string           `json:"location,omitempty"`
	Metrics  *MetricsEndpoint `json:"metrics,omitempty"`
}
