package builders

import (
	"github.com/synadia-io/connect/spec"
)

const DefaultNatsUrl = "nats://localhost:4222"

type NatsConfigBuilder struct {
	nats *spec.NatsConfigSpec
}

func NatsConfig(url string) *NatsConfigBuilder {
	return &NatsConfigBuilder{
		nats: &spec.NatsConfigSpec{
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

func (b *NatsConfigBuilder) Build() spec.NatsConfigSpec {
	return *b.nats
}
