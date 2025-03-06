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
            Nats: model.NatsConfig{
                Url:         ncfg.Url,
                AuthEnabled: ncfg.AuthEnabled,
                Jwt:         ncfg.Jwt,
                Seed:        ncfg.Seed,
            },
        }

        if sp.Consumer.Core != nil {
            result.Consumer.Core = &model.ConsumerStepCore{
                Queue:   sp.Consumer.Core.Queue,
                Subject: sp.Consumer.Core.Subject,
            }
        }

        if sp.Consumer.Stream != nil {
            result.Consumer.Stream = &model.ConsumerStepStream{
                Subject: sp.Consumer.Stream.Subject,
            }
        }

        if sp.Consumer.Kv != nil {
            result.Consumer.Kv = &model.ConsumerStepKv{
                Bucket: sp.Consumer.Kv.Bucket,
                Key:    sp.Consumer.Kv.Key,
            }
        }
    }

    if sp.Producer != nil {
        ncfg := sp.Producer.Nats

        result.Producer = &model.ProducerStep{
            Nats: model.NatsConfig{
                Url:         ncfg.Url,
                AuthEnabled: ncfg.AuthEnabled,
                Jwt:         ncfg.Jwt,
                Seed:        ncfg.Seed,
            },
            Threads: sp.Producer.Threads,
        }

        if sp.Producer.Core != nil {
            result.Producer.Core = &model.ProducerStepCore{
                Subject: sp.Producer.Core.Subject,
            }
        }

        if sp.Producer.Stream != nil {
            result.Producer.Stream = &model.ProducerStepStream{
                Subject: sp.Producer.Stream.Subject,
            }
        }

        if sp.Producer.Kv != nil {
            result.Producer.Kv = &model.ProducerStepKv{
                Bucket: sp.Producer.Kv.Bucket,
                Key:    sp.Producer.Kv.Key,
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
