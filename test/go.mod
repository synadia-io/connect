module github.com/synadia-io/connect.test

go 1.24

toolchain go1.24.0

replace github.com/synadia-io/connect v1.0.0 => ../

require (
	github.com/google/uuid v1.6.0
	github.com/nats-io/jwt/v2 v2.7.3
	github.com/nats-io/nats-server/v2 v2.10.25
	github.com/nats-io/nats.go v1.39.1
	github.com/nats-io/nkeys v0.4.10
	github.com/nats-io/nuid v1.0.1
	github.com/onsi/ginkgo/v2 v2.22.2
	github.com/onsi/gomega v1.36.2
	github.com/rs/zerolog v1.33.0
	github.com/synadia-io/connect v1.0.0
	github.com/synadia-io/connect-node v1.0.1
)

require (
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-task/slim-sprig/v3 v3.0.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/pprof v0.0.0-20241210010833-40e02aabc2ad // indirect
	github.com/klauspost/compress v1.17.11 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/minio/highwayhash v1.0.3 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/net v0.33.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	golang.org/x/time v0.9.0 // indirect
	golang.org/x/tools v0.28.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
