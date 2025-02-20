package convert

import (
    "github.com/synadia-io/connect/model"
    "github.com/synadia-io/connect/spec"
)

func MetricsToSpec(m *model.Metrics) *spec.MetricsSpec {
    if m == nil {
        return nil
    }

    return &spec.MetricsSpec{
        Path: m.Path,
        Port: m.Port,
    }
}

func ConvertStepsToSpec(steps model.Steps) spec.StepsSpec {
    result := spec.StepsSpec{}

    if steps.Consumer != nil {
        ncfg := steps.Consumer.Nats

        result.Consumer = &spec.ConsumerStepSpec{
            Subject: steps.Consumer.Subject,
            Queue:   steps.Consumer.Queue,
            Nats: spec.NatsConfigSpec{
                Url:         ncfg.Url,
                AuthEnabled: ncfg.AuthEnabled,
                Jwt:         ncfg.Jwt,
                Seed:        ncfg.Seed,
            },
        }

        if steps.Consumer.Jetstream != nil {
            js := steps.Consumer.Jetstream
            result.Consumer.Jetstream = &spec.ConsumerStepSpecJetstream{
                Bind:          js.Bind,
                DeliverPolicy: js.DeliverPolicy,
                Durable:       js.Durable,
                MaxAckPending: js.MaxAckPending,
                MaxAckWait:    js.MaxAckWait,
            }
        }
    }

    if steps.Producer != nil {
        ncfg := steps.Producer.Nats

        result.Producer = &spec.ProducerStepSpec{
            Subject: steps.Producer.Subject,
            Nats: spec.NatsConfigSpec{
                Url:         ncfg.Url,
                AuthEnabled: ncfg.AuthEnabled,
                Jwt:         ncfg.Jwt,
                Seed:        ncfg.Seed,
            },
        }

        if steps.Producer.Jetstream != nil {
            js := steps.Producer.Jetstream
            result.Producer.Jetstream = &spec.ProducerStepSpecJetstream{
                MsgId:   js.MsgId,
                AckWait: js.AckWait,
            }

            if steps.Producer.Jetstream.Batching != nil {
                result.Producer.Jetstream.Batching = &spec.ProducerStepSpecJetstreamBatching{
                    Count:    js.Batching.Count,
                    ByteSize: js.Batching.ByteSize,
                }
            }
        }
    }

    if steps.Source != nil {
        result.Source = &spec.SourceStepSpec{
            Type:   steps.Source.Type,
            Config: spec.SourceStepSpecConfig(steps.Source.Config),
        }
    }

    if steps.Sink != nil {
        result.Sink = &spec.SinkStepSpec{
            Type:   steps.Sink.Type,
            Config: spec.SinkStepSpecConfig(steps.Sink.Config),
        }
    }

    if steps.Transformer != nil {
        t := ConvertTransformerToSpec(*steps.Transformer)
        result.Transformer = &t
    }

    return result
}

func ConvertTransformerToSpec(transformer model.TransformerStep) spec.TransformerStepSpec {
    result := spec.TransformerStepSpec{}

    if transformer.Composite != nil {
        result.Composite = &spec.TransformerStepSpecComposite{}

        for _, t := range transformer.Composite.Sequential {
            result.Composite.Sequential = append(result.Composite.Sequential, ConvertTransformerToSpec(t))
        }
    }

    if transformer.Mapping != nil {
        result.Mapping = &spec.TransformerStepSpecMapping{
            Sourcecode: transformer.Mapping.Sourcecode,
        }
    }

    if transformer.Service != nil {
        ncfg := transformer.Service.Nats

        result.Service = &spec.TransformerStepSpecService{
            Endpoint: transformer.Service.Endpoint,
            Timeout:  transformer.Service.Timeout,
            Nats: spec.NatsConfigSpec{
                Url:         ncfg.Url,
                AuthEnabled: ncfg.AuthEnabled,
                Jwt:         ncfg.Jwt,
                Seed:        ncfg.Seed,
            },
        }
    }

    return result
}
