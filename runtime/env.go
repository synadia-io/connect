package runtime

import (
	"fmt"
	"github.com/synadia-io/connect/model"
	"os"
	"strings"
)

const (
	NatsUrlVar      = "CONTROL_NATS_URL"
	NatsAuthVar     = "CONTROL_NATS_AUTH"
	NatsAuthJwtVar  = "CONTROL_NATS_AUTH_JWT"
	NatsAuthSeedVar = "CONTROL_NATS_AUTH_SEED"
)

func StringFromEnv(scope string, key string) string {
	envVar := fmt.Sprintf("%s_%s", strings.ToUpper(scope), strings.ToUpper(key))
	return os.Getenv(envVar)
}

func BoolFromEnv(scope string, key string) bool {
	envVar := fmt.Sprintf("%s_%s", strings.ToUpper(scope), strings.ToUpper(key))
	return os.Getenv(envVar) == "true"
}

func NatsConfigFromEnv(scope string) (*model.NatsConfig, error) {
	result := &model.NatsConfig{}

	result.Url = os.Getenv(NatsUrlVar)
	if result.Url == "" {
		return nil, fmt.Errorf("no %s nats config found", scope)
	}

	result.AuthEnabled = os.Getenv(NatsAuthVar) == "true"
	if result.AuthEnabled {
		result.Jwt = os.Getenv(NatsAuthJwtVar)
		result.Seed = os.Getenv(NatsAuthSeedVar)

		if result.Jwt == "" || result.Seed == "" {
			return nil, fmt.Errorf("jwt and seed must be provided for authenticated nats connection")
		}
	}

	return result, nil
}
