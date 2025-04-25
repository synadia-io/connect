package convert

import (
    "github.com/synadia-io/connect/model"
    "github.com/synadia-io/connect/spec"
)

func ConvertStepsToSpec(steps model.Steps) spec.StepsSpec {
    result := spec.StepsSpec{}

    if steps.Consumer != nil {
        ncfg := steps.Consumer.Nats

        result.Consumer = &spec.ConsumerStepSpec{
            Nats: spec.NatsConfigSpec{
                Url:         ncfg.Url,
                AuthEnabled: ncfg.AuthEnabled,
                Jwt:         ncfg.Jwt,
                Seed:        ncfg.Seed,
            },
        }

        if steps.Consumer.Core != nil {
            result.Consumer.Core = &spec.ConsumerStepSpecCore{
                Queue:   steps.Consumer.Core.Queue,
                Subject: steps.Consumer.Core.Subject,
            }
        }

        if steps.Consumer.Stream != nil {
            result.Consumer.Stream = &spec.ConsumerStepSpecStream{
                Subject: steps.Consumer.Stream.Subject,
            }
        }

        if steps.Consumer.Kv != nil {
            result.Consumer.Kv = &spec.ConsumerStepSpecKv{
                Bucket: steps.Consumer.Kv.Bucket,
                Key:    steps.Consumer.Kv.Key,
            }
        }
    }

    if steps.Producer != nil {
        ncfg := steps.Producer.Nats

        result.Producer = &spec.ProducerStepSpec{
            Nats: spec.NatsConfigSpec{
                Url:         ncfg.Url,
                AuthEnabled: ncfg.AuthEnabled,
                Jwt:         ncfg.Jwt,
                Seed:        ncfg.Seed,
            },
            Threads: steps.Producer.Threads,
        }

        if steps.Producer.Core != nil {
            result.Producer.Core = &spec.ProducerStepSpecCore{
                Subject: steps.Producer.Core.Subject,
            }
        }

        if steps.Producer.Stream != nil {
            result.Producer.Stream = &spec.ProducerStepSpecStream{
                Subject: steps.Producer.Stream.Subject,
            }
        }

        if steps.Producer.Kv != nil {
            result.Producer.Kv = &spec.ProducerStepSpecKv{
                Bucket: steps.Producer.Kv.Bucket,
                Key:    steps.Producer.Kv.Key,
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

    if transformer.Explode != nil {
        result.Explode = &spec.TransformerStepSpecExplode{
            Delimiter: transformer.Explode.Delimiter,
            Format:    spec.TransformerStepSpecExplodeFormat(transformer.Explode.Format),
        }
    }

    if transformer.Combine != nil {
        result.Combine = &spec.TransformerStepSpecCombine{
            Path:   transformer.Combine.Path,
            Format: spec.TransformerStepSpecCombineFormat(transformer.Combine.Format),
        }
    }

    return result
}
