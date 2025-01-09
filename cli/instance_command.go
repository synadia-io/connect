package cli

import (
	"os"

	"github.com/choria-io/fisk"
	"github.com/fatih/color"
)

func init() {
	registerCommand("instance", 3, configureInstanceCommand)
}

type instanceCommand struct {
	connectorId  string
	deploymentId string
	instanceId   string
}

func configureInstanceCommand(parentCmd commandHost) {
	c := &instanceCommand{}

	instanceCmd := parentCmd.Command("instance", "Interact with instances").Alias("i")

	logsCmd := instanceCmd.Command("logs", "List current logs for an instance").Action(c.getInstanceLogs)
	logsCmd.Arg("connector", "The connector to get logs for").Required().StringVar(&c.connectorId)
	logsCmd.Arg("deployment", "The deployment to get logs for").StringVar(&c.deploymentId)
	logsCmd.Arg("id", "The instance to get logs for").StringVar(&c.instanceId)

	eventsCmd := instanceCmd.Command("events", "List current events for an instance").Alias("e").Action(c.getInstanceEvents)
	eventsCmd.Arg("connector", "The connector to get events for").Required().StringVar(&c.connectorId)
	eventsCmd.Arg("deployment", "The deployment to get events for").StringVar(&c.deploymentId)
	eventsCmd.Arg("id", "The instance to get events for").StringVar(&c.instanceId)

	metricsCmd := instanceCmd.Command("metrics", "List current metrics for an instance").Alias("m").Action(c.getInstanceMetrics)
	metricsCmd.Arg("connector", "The connector to get metrics for").Required().StringVar(&c.connectorId)
	metricsCmd.Arg("deployment", "The deployment to get metrics for").StringVar(&c.deploymentId)
	metricsCmd.Arg("id", "The instance to get metrics for").StringVar(&c.instanceId)
}

func (c *instanceCommand) getInstanceLogs(pc *fisk.ParseContext) error {
	logs, err := controlClient().GetLogs(c.connectorId, c.deploymentId, c.instanceId)
	if err != nil {
		color.Red("Could not get logs for instance %s: %s", c.instanceId, err)
		os.Exit(1)
	}

	printLogs(logs)

	return nil
}

func (c *instanceCommand) getInstanceEvents(pc *fisk.ParseContext) error {
	events, err := controlClient().GetEvents(c.connectorId, c.deploymentId, c.instanceId)
	if err != nil {
		color.Red("Could not get events for instance %s: %s", c.instanceId, err)
		os.Exit(1)
	}

	printEvents(events)

	return nil
}

func (c *instanceCommand) getInstanceMetrics(pc *fisk.ParseContext) error {
	metrics, err := controlClient().GetMetrics(c.connectorId, c.deploymentId, c.instanceId)
	if err != nil {
		color.Red("Could not get metrics for instance %s: %s", c.instanceId, err)
		os.Exit(1)
	}

	printMetrics(metrics)

	return nil
}
