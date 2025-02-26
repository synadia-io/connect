package convert

import (
	"github.com/synadia-io/connect/v2/model"
	"github.com/synadia-io/connect/v2/spec"
)

func ConvertStepsFromSpec(sp spec.StepsSpec) model.Steps {
	result := model.Steps{}

	if sp.Consumer != nil {
		ncfg := sp.Consumer.Nats

		result.Consumer = &model.ConsumerStep{
			Subject: sp.Consumer.Subject,
			Queue:   sp.Consumer.Queue,
			Nats: model.NatsConfig{
				Url:         ncfg.Url,
				AuthEnabled: ncfg.AuthEnabled,
				Jwt:         ncfg.Jwt,
				Seed:        ncfg.Seed,
			},
		}

		if sp.Consumer.Jetstream != nil {
			js := sp.Consumer.Jetstream
			result.Consumer.Jetstream = &model.ConsumerStepJetstream{
				Bind:          js.Bind,
				DeliverPolicy: js.DeliverPolicy,
				Durable:       js.Durable,
				MaxAckPending: js.MaxAckPending,
				MaxAckWait:    js.MaxAckWait,
			}
		}
	}

	if sp.Producer != nil {
		ncfg := sp.Producer.Nats

		result.Producer = &model.ProducerStep{
			Subject: sp.Producer.Subject,
			Nats: model.NatsConfig{
				Url:         ncfg.Url,
				AuthEnabled: ncfg.AuthEnabled,
				Jwt:         ncfg.Jwt,
				Seed:        ncfg.Seed,
			},
		}

		if sp.Producer.Jetstream != nil {
			js := sp.Producer.Jetstream
			result.Producer.Jetstream = &model.ProducerStepJetstream{
				MsgId:   js.MsgId,
				AckWait: js.AckWait,
			}

			if js.Batching != nil {
				result.Producer.Jetstream.Batching = &model.ProducerStepJetstreamBatching{
					Count:    js.Batching.Count,
					ByteSize: js.Batching.ByteSize,
				}
			}
		}
	}

	if sp.Source != nil {
		result.Source = &model.SourceStep{
			Type:   sp.Source.Type,
			Config: model.SourceStepConfig(sp.Source.Config),
		}
	}

	if sp.Sink != nil {
		result.Sink = &model.SinkStep{
			Type:   sp.Sink.Type,
			Config: model.SinkStepConfig(sp.Sink.Config),
		}
	}

	if sp.Transformer != nil {
		t := ConvertTransformerFromSpec(*sp.Transformer)
		result.Transformer = &t
	}

	return result
}

func ConvertTransformerFromSpec(sp spec.TransformerStepSpec) model.TransformerStep {
	result := model.TransformerStep{}

	if sp.Service != nil {
		ncfg := sp.Service.Nats

		result.Service = &model.ServiceTransformerStep{
			Endpoint: sp.Service.Endpoint,
			Timeout:  sp.Service.Timeout,
			Nats: model.NatsConfig{
				Url:         ncfg.Url,
				AuthEnabled: ncfg.AuthEnabled,
				Jwt:         ncfg.Jwt,
				Seed:        ncfg.Seed,
			},
		}
	}

	if sp.Mapping != nil {
		result.Mapping = &model.MappingTransformerStep{
			Sourcecode: sp.Mapping.Sourcecode,
		}
	}

	if sp.Composite != nil {
		result.Composite = &model.CompositeTransformerStep{}

		for _, t := range sp.Composite.Sequential {
			result.Composite.Sequential = append(result.Composite.Sequential, ConvertTransformerFromSpec(t))
		}
	}

	return result
}
