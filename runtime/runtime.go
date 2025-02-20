package runtime

import (
    "context"
    "encoding/base64"
    "fmt"
    "github.com/synadia-io/connect/model"
    "gopkg.in/yaml.v3"
    "log/slog"
    "os"
)

const AccountEnvVar = "CONNECT_ACCOUNT"
const ConnectorIdEnvVar = "CONNECT_CONNECTOR_ID"
const InstanceIdEnvVar = "CONNECT_INSTANCE_ID"

func FromEnv() (*Runtime, error) {
    account := os.Getenv(AccountEnvVar)
    if account == "" {
        return nil, fmt.Errorf("%s environment variable not found", AccountEnvVar)
    }

    connectorId := os.Getenv(ConnectorIdEnvVar)
    if connectorId == "" {
        return nil, fmt.Errorf("%s environment variable not found", ConnectorIdEnvVar)
    }

    instanceId := os.Getenv(InstanceIdEnvVar)
    if instanceId == "" {
        return nil, fmt.Errorf("%s environment variable not found", InstanceIdEnvVar)
    }

    return NewRuntime(account, connectorId, instanceId, slog.LevelDebug), nil
}

func NewRuntime(account string, connectorId string, instanceId string, logLevel slog.Level) *Runtime {
    return &Runtime{
        Account:     account,
        ConnectorId: connectorId,
        InstanceId:  instanceId,
        LogLevel:    logLevel,
    }
}

type Workload func(ctx context.Context, runtime *Runtime, steps model.Steps) error

type Runtime struct {
    // Account is the account for which this connector is running
    Account string

    // ConnectorId is the connector id
    ConnectorId string

    // InstanceId identifies this instance of the connector
    InstanceId string

    // LogLevel is the log level for the runtime
    LogLevel slog.Level

    // Logger is the logger for the runtime and only set after launch
    Logger *slog.Logger
}

func (r *Runtime) Launch(ctx context.Context, workload Workload, cfg string) error {
    cfgb, err := base64.StdEncoding.DecodeString(cfg)
    if err != nil {
        return fmt.Errorf("failed to decode config: %w", err)
    }

    // -- decode the connector config
    var steps model.Steps
    if err := yaml.Unmarshal(cfgb, &steps); err != nil {
        return fmt.Errorf("failed to decode connector config: %w", err)
    }

    r.Logger = slog.Default()

    return workload(ctx, r, steps)
}

func (r *Runtime) Close() {}
