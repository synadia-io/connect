package model

import "github.com/nats-io/nats.go"

type NatsConfig struct {
	Url         string `json:"url,omitempty" yaml:"url"`
	AuthEnabled bool   `json:"auth_enabled,omitempty" yaml:"authEnabled,omitempty"`
	Jwt         string `json:"jwt,omitempty" yaml:"jwt,omitempty"`
	Seed        string `json:"seed,omitempty" yaml:"seed,omitempty"`
	Username    string `json:"username,omitempty" yaml:"username,omitempty"`
	Password    string `json:"password,omitempty" yaml:"password,omitempty"`
}

func (c *NatsConfig) Opts() []nats.Option {
	opts := []nats.Option{
		nats.Name("connect"),
	}

	if c.AuthEnabled {
		if c.Username != "" && c.Password != "" {
			opts = append(opts, nats.UserInfo(c.Username, c.Password))
		} else {
			opts = append(opts, nats.UserJWTAndSeed(c.Jwt, c.Seed))
		}
	}

	return opts
}
