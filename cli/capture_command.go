package cli

import (
	"fmt"
	"os"

	"github.com/choria-io/fisk"
	"github.com/fatih/color"
	"github.com/synadia-io/connect/client"
	"github.com/synadia-io/connect/model"
)

func init() {
	registerCommand("capture", 4, configureCaptureCommand)
}

type captureCommand struct {
	connectorId  string
	deploymentId string
	instanceId   string
	history      bool
}

func configureCaptureCommand(parentCmd commandHost) {
	c := &captureCommand{}

	logsCmd := parentCmd.Command("logs", "Show instance logs").Alias("log").Action(c.captureLogs)
	logsCmd.Arg("cid", "The connector id of instances for which to show the logs").StringVar(&c.connectorId)
	logsCmd.Arg("did", "The deployment id of instances for which to show the logs").StringVar(&c.deploymentId)
	logsCmd.Arg("id", "The instance id of an instance for which to show the logs").StringVar(&c.instanceId)
	logsCmd.Flag("history", "Show instance logs history").UnNegatableBoolVar(&c.history)

	eventsCmd := parentCmd.Command("events", "Show instance events").Alias("evt").Action(c.captureEvents)
	eventsCmd.Arg("cid", "The connector id of instances for which to show the events").StringVar(&c.connectorId)
	eventsCmd.Arg("did", "The deployment id of instances for which to show the events").StringVar(&c.deploymentId)
	eventsCmd.Arg("id", "The instance id of an instance for which to show the events").StringVar(&c.instanceId)
	eventsCmd.Flag("history", "Show instance events history").UnNegatableBoolVar(&c.history)

	metricsCmd := parentCmd.Command("metrics", "Show instance metrics").Alias("evt").Action(c.captureMetrics)
	metricsCmd.Arg("cid", "The connector id of instances for which to show the metrics").StringVar(&c.connectorId)
	metricsCmd.Arg("did", "The deployment id of instances for which to show the metrics").StringVar(&c.deploymentId)
	metricsCmd.Arg("id", "The instance id of an instance for which to show the metrics").StringVar(&c.instanceId)
	metricsCmd.Flag("history", "Show instance metrics history").UnNegatableBoolVar(&c.history)
}

func (c *captureCommand) captureLogs(pc *fisk.ParseContext) error {
	filter := client.CaptureFilter{
		ConnectorId:  c.connectorId,
		DeploymentId: c.deploymentId,
		InstanceId:   c.instanceId,
	}

	client := controlClient()

	if c.history {
		logs, err := client.GetLogs(c.connectorId, c.deploymentId, c.instanceId)
		if err != nil {
			color.Red("Could not get log history")
		} else {
			for _, log := range logs {
				RenderLog(*log)
			}
		}
	}

	fmt.Println("Capturing logs, press Ctrl+C to stop")
	s, err := client.CaptureLogs(filter, func(item model.InstanceLog) {
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
		InstanceId:   c.instanceId,
	}

	client := controlClient()

	if c.history {
		events, err := client.GetEvents(c.connectorId, c.deploymentId, c.instanceId)
		if err != nil {
			color.Red("Could not get event history")
		} else {
			for _, event := range events {
				RenderEvent(*event)
			}
		}
	}

	fmt.Println("Capturing events, press Ctrl+C to stop")
	s, err := client.CaptureEvents(filter, func(item model.InstanceEvent) {
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
		InstanceId:   c.instanceId,
	}

	client := controlClient()

	if c.history {
		metrics, err := client.GetMetrics(c.connectorId, c.deploymentId, c.instanceId)
		if err != nil {
			color.Red("Could not get metric history")
		} else {
			for _, metric := range metrics {
				RenderMetric(*metric)
			}
		}
	}

	fmt.Println("Capturing metrics, press Ctrl+C to stop")
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
