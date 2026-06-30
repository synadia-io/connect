package runtime

const (
	NatsUrlVar  = "NEX_WORKLOAD_NATS_SERVERS"
	NatsJwtVar  = "NEX_WORKLOAD_NATS_B64_JWT"
	NatsSeedVar = "NEX_WORKLOAD_NATS_NKEY"

	GroupEnvVar     = "NEX_WORKLOAD_GROUP"
	NamespaceEnvVar = "NEX_WORKLOAD_NAMESPACE"
	InstanceEnvVar  = "NEX_WORKLOAD_ID"

	// NatsCredsFileVar, when set, points at a decorated NATS creds file that the
	// runtime loads at launch and watches for credential refreshes.
	NatsCredsFileVar = "NEX_WORKLOAD_NATS_CREDS_FILE"

	LogLevelEnvVar = "CONNECT_LOG_LEVEL"
)
