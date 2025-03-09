package builders

import "github.com/synadia-io/connect/model"

const DefaultNatsUrl = "nats://localhost:4222"

type NatsConfigBuilder struct {
    nats *model.NatsConfig
}

func NatsConfig(url string) *NatsConfigBuilder {
    return &NatsConfigBuilder{
        nats: &model.NatsConfig{
            Url: url,
        },
    }
}

func (b *NatsConfigBuilder) Auth(jwt string, seed string) *NatsConfigBuilder {
    b.nats.AuthEnabled = true
    b.nats.Jwt = &jwt
    b.nats.Seed = &seed
    return b
}

func (b *NatsConfigBuilder) Build() model.NatsConfig {
    return *b.nats
}
