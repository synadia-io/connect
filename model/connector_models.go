// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package model

import "encoding/json"
import "fmt"
import "reflect"

// A composite transformer which can be used to combine several transformers
type CompositeTransformerStep struct {
	// Sequential corresponds to the JSON schema field "sequential".
	Sequential []TransformerStep `json:"sequential" yaml:"sequential" mapstructure:"sequential"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *CompositeTransformerStep) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["sequential"]; raw != nil && !ok {
		return fmt.Errorf("field sequential in CompositeTransformerStep: required")
	}
	type Plain CompositeTransformerStep
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	*j = CompositeTransformerStep(plain)
	return nil
}

type Connector struct {
	// The unique id of the connector
	ConnectorId string `json:"connector_id" yaml:"connector_id" mapstructure:"connector_id"`

	// A description of the connector
	Description string `json:"description" yaml:"description" mapstructure:"description"`

	// The id of the connector's runtime
	RuntimeId string `json:"runtime_id" yaml:"runtime_id" mapstructure:"runtime_id"`

	// Steps corresponds to the JSON schema field "steps".
	Steps Steps `json:"steps" yaml:"steps" mapstructure:"steps"`
}

type ConnectorStartOptions struct {
	// The environment variables to set
	EnvVars ConnectorStartOptionsEnvVars `json:"env_vars,omitempty" yaml:"env_vars,omitempty" mapstructure:"env_vars,omitempty"`

	// The placement tags influencing the placement of connector instances
	PlacementTags []string `json:"placement_tags,omitempty" yaml:"placement_tags,omitempty" mapstructure:"placement_tags,omitempty"`

	// Whether the image to run the connector should be pulled
	Pull bool `json:"pull,omitempty" yaml:"pull,omitempty" mapstructure:"pull,omitempty"`

	// The authentication to use when pulling the image
	PullAuth *ConnectorStartOptionsPullAuth `json:"pull_auth,omitempty" yaml:"pull_auth,omitempty" mapstructure:"pull_auth,omitempty"`

	// The number of replicas to run
	Replicas int `json:"replicas,omitempty" yaml:"replicas,omitempty" mapstructure:"replicas,omitempty"`

	// The timeout for the start operation
	Timeout string `json:"timeout,omitempty" yaml:"timeout,omitempty" mapstructure:"timeout,omitempty"`
}

// The environment variables to set
type ConnectorStartOptionsEnvVars map[string]string

// The authentication to use when pulling the image
type ConnectorStartOptionsPullAuth struct {
	// Whether to use authentication
	Enabled bool `json:"enabled,omitempty" yaml:"enabled,omitempty" mapstructure:"enabled,omitempty"`

	// The password to use
	Password *string `json:"password,omitempty" yaml:"password,omitempty" mapstructure:"password,omitempty"`

	// The username to use
	Username *string `json:"username,omitempty" yaml:"username,omitempty" mapstructure:"username,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ConnectorStartOptionsPullAuth) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	type Plain ConnectorStartOptionsPullAuth
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["enabled"]; !ok || v == nil {
		plain.Enabled = false
	}
	*j = ConnectorStartOptionsPullAuth(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ConnectorStartOptions) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	type Plain ConnectorStartOptions
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["env_vars"]; !ok || v == nil {
		plain.EnvVars = map[string]string{}
	}
	if v, ok := raw["placement_tags"]; !ok || v == nil {
		plain.PlacementTags = []string{}
	}
	if v, ok := raw["pull"]; !ok || v == nil {
		plain.Pull = false
	}
	if v, ok := raw["replicas"]; !ok || v == nil {
		plain.Replicas = 1.0
	}
	if v, ok := raw["timeout"]; !ok || v == nil {
		plain.Timeout = "5m"
	}
	*j = ConnectorStartOptions(plain)
	return nil
}

type ConnectorStatus struct {
	// The number of pending instances
	Pending int `json:"pending" yaml:"pending" mapstructure:"pending"`

	// The number of running instances
	Running int `json:"running" yaml:"running" mapstructure:"running"`

	// The number of stopped instances
	Stopped int `json:"stopped" yaml:"stopped" mapstructure:"stopped"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ConnectorStatus) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	type Plain ConnectorStatus
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["pending"]; !ok || v == nil {
		plain.Pending = 0.0
	}
	if v, ok := raw["running"]; !ok || v == nil {
		plain.Running = 0.0
	}
	if v, ok := raw["stopped"]; !ok || v == nil {
		plain.Stopped = 0.0
	}
	*j = ConnectorStatus(plain)
	return nil
}

type ConnectorSummary struct {
	// The unique id of the connector
	ConnectorId string `json:"connector_id" yaml:"connector_id" mapstructure:"connector_id"`

	// A description of the connector
	Description string `json:"description" yaml:"description" mapstructure:"description"`

	// Instances corresponds to the JSON schema field "instances".
	Instances ConnectorSummaryInstances `json:"instances" yaml:"instances" mapstructure:"instances"`

	// The id of the connector's runtime
	RuntimeId string `json:"runtime_id" yaml:"runtime_id" mapstructure:"runtime_id"`
}

type ConnectorSummaryInstances struct {
	// The number of pending instances
	Pending int `json:"pending" yaml:"pending" mapstructure:"pending"`

	// The number of running instances
	Running int `json:"running" yaml:"running" mapstructure:"running"`

	// The number of stopped instances
	Stopped int `json:"stopped" yaml:"stopped" mapstructure:"stopped"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ConnectorSummaryInstances) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	type Plain ConnectorSummaryInstances
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["pending"]; !ok || v == nil {
		plain.Pending = 0.0
	}
	if v, ok := raw["running"]; !ok || v == nil {
		plain.Running = 0.0
	}
	if v, ok := raw["stopped"]; !ok || v == nil {
		plain.Stopped = 0.0
	}
	*j = ConnectorSummaryInstances(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ConnectorSummary) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["connector_id"]; raw != nil && !ok {
		return fmt.Errorf("field connector_id in ConnectorSummary: required")
	}
	if _, ok := raw["description"]; raw != nil && !ok {
		return fmt.Errorf("field description in ConnectorSummary: required")
	}
	if _, ok := raw["instances"]; raw != nil && !ok {
		return fmt.Errorf("field instances in ConnectorSummary: required")
	}
	if _, ok := raw["runtime_id"]; raw != nil && !ok {
		return fmt.Errorf("field runtime_id in ConnectorSummary: required")
	}
	type Plain ConnectorSummary
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	*j = ConnectorSummary(plain)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Connector) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["connector_id"]; raw != nil && !ok {
		return fmt.Errorf("field connector_id in Connector: required")
	}
	if _, ok := raw["description"]; raw != nil && !ok {
		return fmt.Errorf("field description in Connector: required")
	}
	if _, ok := raw["runtime_id"]; raw != nil && !ok {
		return fmt.Errorf("field runtime_id in Connector: required")
	}
	if _, ok := raw["steps"]; raw != nil && !ok {
		return fmt.Errorf("field steps in Connector: required")
	}
	type Plain Connector
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	*j = Connector(plain)
	return nil
}

// The consumer reading messages from NATS
type ConsumerStep struct {
	// The JetStream configuration
	Jetstream *ConsumerStepJetstream `json:"jetstream,omitempty" yaml:"jetstream,omitempty" mapstructure:"jetstream,omitempty"`

	// Nats corresponds to the JSON schema field "nats".
	Nats NatsConfig `json:"nats" yaml:"nats" mapstructure:"nats"`

	// The queue this connector should join
	Queue *string `json:"queue,omitempty" yaml:"queue,omitempty" mapstructure:"queue,omitempty"`

	// The subject to read messages from
	Subject string `json:"subject" yaml:"subject" mapstructure:"subject"`
}

// The JetStream configuration
type ConsumerStepJetstream struct {
	// Whether to bind to the durable
	Bind *bool `json:"bind,omitempty" yaml:"bind,omitempty" mapstructure:"bind,omitempty"`

	// The JetStream deliver policy
	DeliverPolicy *string `json:"deliver_policy,omitempty" yaml:"deliver_policy,omitempty" mapstructure:"deliver_policy,omitempty"`

	// The durable name
	Durable *string `json:"durable,omitempty" yaml:"durable,omitempty" mapstructure:"durable,omitempty"`

	// The maximum number of acks pending
	MaxAckPending *int `json:"max_ack_pending,omitempty" yaml:"max_ack_pending,omitempty" mapstructure:"max_ack_pending,omitempty"`

	// The maximum ack wait time
	MaxAckWait *string `json:"max_ack_wait,omitempty" yaml:"max_ack_wait,omitempty" mapstructure:"max_ack_wait,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ConsumerStep) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["nats"]; raw != nil && !ok {
		return fmt.Errorf("field nats in ConsumerStep: required")
	}
	if _, ok := raw["subject"]; raw != nil && !ok {
		return fmt.Errorf("field subject in ConsumerStep: required")
	}
	type Plain ConsumerStep
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	*j = ConsumerStep(plain)
	return nil
}

type Instance struct {
	// the id of the connector this instance belongs to
	ConnectorId string `json:"connector_id" yaml:"connector_id" mapstructure:"connector_id"`

	// The unique id of the instance
	Id string `json:"id" yaml:"id" mapstructure:"id"`

	// The message associated with the status. This can be an error message or just a
	// completion message
	Message *string `json:"message,omitempty" yaml:"message,omitempty" mapstructure:"message,omitempty"`

	// The status of the instance
	Status InstanceStatus `json:"status" yaml:"status" mapstructure:"status"`

	// The amount of time the instance has been running
	Uptime *string `json:"uptime,omitempty" yaml:"uptime,omitempty" mapstructure:"uptime,omitempty"`
}

type InstanceStatus string

const InstanceStatusPending InstanceStatus = "pending"
const InstanceStatusRunning InstanceStatus = "running"
const InstanceStatusStopped InstanceStatus = "stopped"
const InstanceStatusUnknown InstanceStatus = "unknown"

var enumValues_InstanceStatus = []interface{}{
	"pending",
	"running",
	"stopped",
	"unknown",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *InstanceStatus) UnmarshalJSON(value []byte) error {
	var v string
	if err := json.Unmarshal(value, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_InstanceStatus {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_InstanceStatus, v)
	}
	*j = InstanceStatus(v)
	return nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Instance) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["connector_id"]; raw != nil && !ok {
		return fmt.Errorf("field connector_id in Instance: required")
	}
	if _, ok := raw["id"]; raw != nil && !ok {
		return fmt.Errorf("field id in Instance: required")
	}
	if _, ok := raw["status"]; raw != nil && !ok {
		return fmt.Errorf("field status in Instance: required")
	}
	type Plain Instance
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	*j = Instance(plain)
	return nil
}

// A mapping transformer which can transform the message
type MappingTransformerStep struct {
	// The source code of the mapping transformer
	Sourcecode string `json:"sourcecode" yaml:"sourcecode" mapstructure:"sourcecode"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *MappingTransformerStep) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["sourcecode"]; raw != nil && !ok {
		return fmt.Errorf("field sourcecode in MappingTransformerStep: required")
	}
	type Plain MappingTransformerStep
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	*j = MappingTransformerStep(plain)
	return nil
}

// Information on how to collect metrics. If not set, no metrics will be collected
type Metrics struct {
	// The path to collect metrics from
	Path *string `json:"path,omitempty" yaml:"path,omitempty" mapstructure:"path,omitempty"`

	// The port to collect metrics from
	Port int `json:"port" yaml:"port" mapstructure:"port"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *Metrics) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["port"]; raw != nil && !ok {
		return fmt.Errorf("field port in Metrics: required")
	}
	type Plain Metrics
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	*j = Metrics(plain)
	return nil
}

type NatsConfig struct {
	// Whether authentication is enabled
	AuthEnabled bool `json:"auth_enabled,omitempty" yaml:"auth_enabled,omitempty" mapstructure:"auth_enabled,omitempty"`

	// The JWT token used during authentication. Only applicable if auth_enabled is
	// true
	Jwt *string `json:"jwt,omitempty" yaml:"jwt,omitempty" mapstructure:"jwt,omitempty"`

	// The seed used during authentication. Only applicable if auth_enabled is true
	Seed *string `json:"seed,omitempty" yaml:"seed,omitempty" mapstructure:"seed,omitempty"`

	// The url of the nats server to connect to
	Url string `json:"url" yaml:"url" mapstructure:"url"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *NatsConfig) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["url"]; raw != nil && !ok {
		return fmt.Errorf("field url in NatsConfig: required")
	}
	type Plain NatsConfig
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["auth_enabled"]; !ok || v == nil {
		plain.AuthEnabled = false
	}
	*j = NatsConfig(plain)
	return nil
}

// The producer writing messages to NATS
type ProducerStep struct {
	// The JetStream configuration
	Jetstream *ProducerStepJetstream `json:"jetstream,omitempty" yaml:"jetstream,omitempty" mapstructure:"jetstream,omitempty"`

	// Nats corresponds to the JSON schema field "nats".
	Nats NatsConfig `json:"nats" yaml:"nats" mapstructure:"nats"`

	// The subject to write messages to
	Subject string `json:"subject" yaml:"subject" mapstructure:"subject"`

	// The number of threads used to write messages.
	Threads int `json:"threads,omitempty" yaml:"threads,omitempty" mapstructure:"threads,omitempty"`
}

// The JetStream configuration
type ProducerStepJetstream struct {
	// The ack wait time
	AckWait *string `json:"ack_wait,omitempty" yaml:"ack_wait,omitempty" mapstructure:"ack_wait,omitempty"`

	// The Batching Policy
	Batching *ProducerStepJetstreamBatching `json:"batching,omitempty" yaml:"batching,omitempty" mapstructure:"batching,omitempty"`

	// The message id to allow stream message deduplication
	MsgId *string `json:"msg_id,omitempty" yaml:"msg_id,omitempty" mapstructure:"msg_id,omitempty"`
}

// The Batching Policy
type ProducerStepJetstreamBatching struct {
	// The size of the batch
	ByteSize *int `json:"byte_size,omitempty" yaml:"byte_size,omitempty" mapstructure:"byte_size,omitempty"`

	// The number of messages to batch
	Count *int `json:"count,omitempty" yaml:"count,omitempty" mapstructure:"count,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ProducerStep) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["nats"]; raw != nil && !ok {
		return fmt.Errorf("field nats in ProducerStep: required")
	}
	if _, ok := raw["subject"]; raw != nil && !ok {
		return fmt.Errorf("field subject in ProducerStep: required")
	}
	type Plain ProducerStep
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["threads"]; !ok || v == nil {
		plain.Threads = 1.0
	}
	*j = ProducerStep(plain)
	return nil
}

// A service transformer sends each message to a nats service to be transformed
type ServiceTransformerStep struct {
	// The nats subject on which the service is receiving requests
	Endpoint string `json:"endpoint" yaml:"endpoint" mapstructure:"endpoint"`

	// Nats corresponds to the JSON schema field "nats".
	Nats NatsConfig `json:"nats" yaml:"nats" mapstructure:"nats"`

	// The timeout for the service call
	Timeout string `json:"timeout,omitempty" yaml:"timeout,omitempty" mapstructure:"timeout,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *ServiceTransformerStep) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["endpoint"]; raw != nil && !ok {
		return fmt.Errorf("field endpoint in ServiceTransformerStep: required")
	}
	if _, ok := raw["nats"]; raw != nil && !ok {
		return fmt.Errorf("field nats in ServiceTransformerStep: required")
	}
	type Plain ServiceTransformerStep
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	if v, ok := raw["timeout"]; !ok || v == nil {
		plain.Timeout = "5s"
	}
	*j = ServiceTransformerStep(plain)
	return nil
}

// The external system that is the target for the messages
type SinkStep struct {
	// The configuration of the sink step
	Config SinkStepConfig `json:"config" yaml:"config" mapstructure:"config"`

	// The type of the sink step. This should be a sink that is included in the
	// connector's runtime
	Type string `json:"type" yaml:"type" mapstructure:"type"`
}

// The configuration of the sink step
type SinkStepConfig map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SinkStep) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["config"]; raw != nil && !ok {
		return fmt.Errorf("field config in SinkStep: required")
	}
	if _, ok := raw["type"]; raw != nil && !ok {
		return fmt.Errorf("field type in SinkStep: required")
	}
	type Plain SinkStep
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	*j = SinkStep(plain)
	return nil
}

// The external system that is the source of the messages
type SourceStep struct {
	// The configuration of the source step
	Config SourceStepConfig `json:"config" yaml:"config" mapstructure:"config"`

	// The type of the source step. This should be a source that is included in the
	// connector's runtime
	Type string `json:"type" yaml:"type" mapstructure:"type"`
}

// The configuration of the source step
type SourceStepConfig map[string]interface{}

// UnmarshalJSON implements json.Unmarshaler.
func (j *SourceStep) UnmarshalJSON(value []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(value, &raw); err != nil {
		return err
	}
	if _, ok := raw["config"]; raw != nil && !ok {
		return fmt.Errorf("field config in SourceStep: required")
	}
	if _, ok := raw["type"]; raw != nil && !ok {
		return fmt.Errorf("field type in SourceStep: required")
	}
	type Plain SourceStep
	var plain Plain
	if err := json.Unmarshal(value, &plain); err != nil {
		return err
	}
	*j = SourceStep(plain)
	return nil
}

type Steps struct {
	// Consumer corresponds to the JSON schema field "consumer".
	Consumer *ConsumerStep `json:"consumer,omitempty" yaml:"consumer,omitempty" mapstructure:"consumer,omitempty"`

	// Producer corresponds to the JSON schema field "producer".
	Producer *ProducerStep `json:"producer,omitempty" yaml:"producer,omitempty" mapstructure:"producer,omitempty"`

	// Sink corresponds to the JSON schema field "sink".
	Sink *SinkStep `json:"sink,omitempty" yaml:"sink,omitempty" mapstructure:"sink,omitempty"`

	// Source corresponds to the JSON schema field "source".
	Source *SourceStep `json:"source,omitempty" yaml:"source,omitempty" mapstructure:"source,omitempty"`

	// Transformer corresponds to the JSON schema field "transformer".
	Transformer *TransformerStep `json:"transformer,omitempty" yaml:"transformer,omitempty" mapstructure:"transformer,omitempty"`
}

// The transformer for messages flowing through the connector
type TransformerStep struct {
	// Composite corresponds to the JSON schema field "composite".
	Composite *CompositeTransformerStep `json:"composite,omitempty" yaml:"composite,omitempty" mapstructure:"composite,omitempty"`

	// Mapping corresponds to the JSON schema field "mapping".
	Mapping *MappingTransformerStep `json:"mapping,omitempty" yaml:"mapping,omitempty" mapstructure:"mapping,omitempty"`

	// Service corresponds to the JSON schema field "service".
	Service *ServiceTransformerStep `json:"service,omitempty" yaml:"service,omitempty" mapstructure:"service,omitempty"`
}
