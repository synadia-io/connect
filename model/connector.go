package model

type ConnectorKind string

const (
	Inlet       ConnectorKind = "inlet"
	Outlet      ConnectorKind = "outlet"
	UnknownKind ConnectorKind = ""
)

type Connector struct {
	ConnectorConfig

	Id            string   `json:"id" yaml:"id"`
	DeploymentIds []string `json:"deployment_ids,omitempty"`
}

type ConnectorConfig struct {
	Description string `json:"description,omitempty" yaml:"description"`

	Workload string           `json:"workload,omitempty" yaml:"workload,omitempty"`
	Metrics  *MetricsEndpoint `json:"metrics,omitempty" yaml:"metrics,omitempty"`

	Steps *Steps `json:"steps,omitempty" yaml:"steps,omitempty"`
}

type Steps struct {
	Source   *Source   `json:"source,omitempty" yaml:"source,omitempty"`
	Consumer *Consumer `json:"consumer,omitempty" yaml:"consumer,omitempty"`

	Transformer *Transformer `json:"transformer,omitempty" yaml:"transformer,omitempty"`

	Producer *Producer `json:"producer,omitempty" yaml:"producer,omitempty"`
	Sink     *Sink     `json:"sink,omitempty" yaml:"sink,omitempty"`
}

type MetricsEndpoint struct {
	Port int    `json:"port" yaml:"port,omitempty"`
	Path string `json:"path" yaml:"path,omitempty"`
}

type Source struct {
	Type   string         `json:"type" yaml:"type"`
	Config map[string]any `json:"config" yaml:"config"`
}

type Sink struct {
	Type   string         `json:"type" yaml:"type"`
	Config map[string]any `json:"config" yaml:"config"`
}

type Consumer struct {
	NatsConfig NatsConfig `json:"nats_config" yaml:"nats_config"`
	Subject    string     `json:"subject" yaml:"subject"`

	Queue string `json:"queue,omitempty" yaml:"queue,omitempty"`

	JetStream *ConsumerJetStreamOptions `json:"jetstream,omitempty"`
}

type ConsumerJetStreamOptions struct {
	DeliverPolicy string `json:"deliver_policy,omitempty" yaml:"deliverPolicy,omitempty"`
	MaxAckPending int    `json:"max_ack_pending,omitempty" yaml:"maxAckPending,omitempty"`
	MaxAckWait    string `json:"max_ack_wait,omitempty" yaml:"maxAckWait,omitempty"`
	Durable       string `json:"durable,omitempty" yaml:"durable,omitempty"`
	Bind          bool   `json:"bind,omitempty" yaml:"bind,omitempty"`
}

type Producer struct {
	NatsConfig NatsConfig `json:"nats_config" yaml:"nats_config"`
	Subject    string     `json:"subject" yaml:"subject"`
	Threads    uint       `json:"threads,omitempty" yaml:"threads,omitempty"`

	JetStream *ProducerJetStreamOptions `json:"jetstream,omitempty" yaml:"jetstream,omitempty"`
}

type ProducerJetStreamOptions struct {
	MsgId    string       `json:"msg_id,omitempty" yaml:"msgId,omitempty"`
	AckWait  string       `json:"ack_wait,omitempty" yaml:"ackWait,omitempty"`
	Batching *BatchPolicy `json:"batching,omitempty" yaml:"batching,omitempty"`
}

type BatchPolicy struct {
	Count    int `json:"count,omitempty" yaml:"count,omitempty"`
	ByteSize int `json:"byte_size,omitempty" yaml:"byteSize,omitempty"`
}

func (s *Steps) Kind() ConnectorKind {
	if s.Source != nil && s.Producer != nil {
		return Inlet
	}

	if s.Sink != nil && s.Consumer != nil {
		return Outlet
	}

	return UnknownKind

}
