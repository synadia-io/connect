package cli

import (
	"fmt"
	"os"

	"github.com/choria-io/fisk"
	"github.com/synadia-io/connect/client"
	"github.com/synadia-io/connect/model"
)

func init() {
	registerCommand("capture", 4, configureCaptureCommand)
}

type captureCommand struct {
	connectorId  string
	deploymentId string
}

func configureCaptureCommand(parentCmd commandHost) {
	c := &captureCommand{}

	logsCmd := parentCmd.Command("logs", "Show instance logs").Alias("log").Action(c.captureLogs)
	logsCmd.Arg("cid", "The connector id of instances for which to show the logs").StringVar(&c.connectorId)
	logsCmd.Arg("did", "The deployment id of instances for which to show the logs").StringVar(&c.deploymentId)

	eventsCmd := parentCmd.Command("events", "Show instance events").Alias("evt").Action(c.captureEvents)
	eventsCmd.Arg("cid", "The connector id of instances for which to show the events").StringVar(&c.connectorId)
	eventsCmd.Arg("did", "The deployment id of instances for which to show the events").StringVar(&c.deploymentId)

	metricsCmd := parentCmd.Command("metrics", "Show instance metrics").Alias("evt").Action(c.captureMetrics)
	metricsCmd.Arg("cid", "The connector id of instances for which to show the metrics").StringVar(&c.connectorId)
	metricsCmd.Arg("did", "The deployment id of instances for which to show the metrics").StringVar(&c.deploymentId)
}

func (c *captureCommand) captureLogs(pc *fisk.ParseContext) error {
	filter := client.CaptureFilter{
		ConnectorId:  c.connectorId,
		DeploymentId: c.deploymentId,
	}

	fmt.Println("Capturing logs, press Ctrl+C to stop")
	s, err := controlClient().CaptureLogs(filter, func(item model.InstanceLog) {
		RenderLog(item)
	})
	if err != nil {
		return err
	}
	defer func() {
		if s != nil {
			_ = s.Unsubscribe()
		}
	}()

	sigs := make(chan os.Signal, 1)
	<-sigs

	return nil
}

func (c *captureCommand) captureEvents(pc *fisk.ParseContext) error {
	filter := client.CaptureFilter{
		ConnectorId:  c.connectorId,
		DeploymentId: c.deploymentId,
	}

	fmt.Println("Capturing events, press Ctrl+C to stop")
	s, err := controlClient().CaptureEvents(filter, func(item model.InstanceEvent) {
		RenderEvent(item)
	})
	if err != nil {
		return err
	}
	defer func() {
		if s != nil {
			_ = s.Unsubscribe()
		}
	}()

	sigs := make(chan os.Signal, 1)
	<-sigs

	return nil
}

func (c *captureCommand) captureMetrics(pc *fisk.ParseContext) error {
	filter := client.CaptureFilter{
		ConnectorId:  c.connectorId,
		DeploymentId: c.deploymentId,
	}

	s, err := controlClient().CaptureMetrics(filter, func(item model.InstanceMetric) {
		RenderMetric(item)
	})
	if err != nil {
		return err
	}
	defer func() {
		if s != nil {
			_ = s.Unsubscribe()
		}
	}()

	sigs := make(chan os.Signal, 1)
	<-sigs

	return nil
}
