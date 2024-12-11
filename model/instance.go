package model

import (
	"time"
)

var (
	InstanceScheduled InstanceStatus = "scheduled"
	InstanceCreated   InstanceStatus = "created"
	InstanceRunning   InstanceStatus = "running"
	InstanceStopped   InstanceStatus = "stopped"
	InstanceFailed    InstanceStatus = "failed"
)

var (
	ScheduledEventType InstanceEventType = "scheduled"
	CreatedEventType   InstanceEventType = "created"
	RunningEventType   InstanceEventType = "running"
	StoppedEventType   InstanceEventType = "stopped"
	FailedEventType    InstanceEventType = "failed"
)

const EventTypeHeaderName = "Nats-Event-Type"

type (
	InstanceStatus    string
	InstanceEventType string

	Instance struct {
		ConnectorId  string `json:"connector_id"`
		DeploymentId string `json:"deployment_id"`
		InstanceId   string `json:"instance_id"`

		Status InstanceStatus `json:"status"`

		ExitCode    int        `json:"exit_code,omitempty"`
		Error       string     `json:"error,omitempty"`
		ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
		StartedAt   *time.Time `json:"started_at,omitempty"`
		FinishedAt  *time.Time `json:"finished_at,omitempty"`
	}

	InstanceLog struct {
		ConnectorId  string    `json:"connector_id"`
		DeploymentId string    `json:"deployment_id"`
		InstanceId   string    `json:"instance_id"`
		Timestamp    time.Time `json:"timestamp"`
		Line         string    `json:"line"`
	}

	InstanceMetric struct {
		ConnectorId  string    `json:"connector_id"`
		DeploymentId string    `json:"deployment_id"`
		InstanceId   string    `json:"instance_id"`
		Timestamp    time.Time `json:"timestamp"`
		Data         []byte    `json:"data"`
	}

	InstanceEvent struct {
		ConnectorId  string            `json:"connector_id"`
		DeploymentId string            `json:"deployment_id"`
		InstanceId   string            `json:"instance_id"`
		Type         InstanceEventType `json:"type"`
		Timestamp    time.Time         `json:"timestamp"`
		ExitCode     int               `json:"exit_code,omitempty"`
		Error        string            `json:"error,omitempty"`
	}
)
