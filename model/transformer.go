package model

type Transformer struct {
	Composite *CompositeTransformer `json:"composite,omitempty"`
	Service   *ServiceTransformer   `json:"service,omitempty"`
	Mapping   *MappingTransformer   `json:"mapping,omitempty"`
}

type CompositeTransformer struct {
	Sequential []Transformer `json:"sequential"`
}

type ServiceTransformer struct {
	Endpoint   string     `json:"endpoint"`
	NatsConfig NatsConfig `json:"nats_config"`
	Timeout    string     `json:"timeout,omitempty"`
}

type MappingTransformer struct {
	Sourcecode string `json:"sourcecode"`
}
