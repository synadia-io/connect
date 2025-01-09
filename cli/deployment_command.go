package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/choria-io/fisk"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/synadia-io/connect/cli/render"
	"github.com/synadia-io/connect/client"
	"github.com/synadia-io/connect/model"
)

func init() {
	registerCommand("deployment", 2, configureDeploymentCommand)
}

type deploymentCommand struct {
	connectorId    string
	deploymentId   string
	instanceStatus []string
}

func configureDeploymentCommand(parentCmd commandHost) {
	c := &deploymentCommand{}

	deploymentCmd := parentCmd.Command("deployment", "Interact with deployments").Alias("d")

	listCmd := deploymentCmd.Command("list", "List deployments for a connector").Alias("ls").Action(c.listDeployments)
	listCmd.Arg("connector", "The connector to list deployments for").StringVar(&c.connectorId)
	listCmd.Flag("status", "At least one instance needs to have this status").StringsVar(&c.instanceStatus)

	infoCmd := deploymentCmd.Command("get", "Show details about a deployment").Alias("info").Action(c.describeDeployment)
	infoCmd.Arg("connector", "The connector to show deployment details for").Required().StringVar(&c.connectorId)
	infoCmd.Arg("id", "The deployment id to show details for").StringVar(&c.deploymentId)

	purgeCmd := deploymentCmd.Command("purge", "Remove all stopped or failed deployments").Action(c.purgeDeployments)
	purgeCmd.Arg("connector", "The connector to purge deployments for").StringVar(&c.connectorId)

	logsCmd := deploymentCmd.Command("logs", "List current logs for a deployment").Action(c.getDeploymentLogs)
	logsCmd.Arg("connector", "The connector to get logs for").Required().StringVar(&c.connectorId)
	logsCmd.Arg("id", "The deployment id to get logs for").StringVar(&c.deploymentId)

	eventsCmd := deploymentCmd.Command("events", "List current events for a deployment").Alias("e").Action(c.getDeploymentEvents)
	eventsCmd.Arg("connector", "The connector to get events for").Required().StringVar(&c.connectorId)
	eventsCmd.Arg("id", "The deployment id to get events for").StringVar(&c.deploymentId)

	metricsCmd := deploymentCmd.Command("metrics", "List current metrics for a deployment").Alias("m").Action(c.getDeploymentMetrics)
	metricsCmd.Arg("connector", "The connector to get metrics for").Required().StringVar(&c.connectorId)
	metricsCmd.Arg("id", "The deployment id to get metrics for").StringVar(&c.deploymentId)
}

func (c *deploymentCommand) listDeployments(pc *fisk.ParseContext) error {
	filter := client.DeploymentFilter{
		ConnectorId: c.connectorId,
	}

	for _, status := range c.instanceStatus {
		filter.InstanceStatus = append(filter.InstanceStatus, model.InstanceStatus(status))
	}

	w := table.NewWriter()
	w.AppendHeader(table.Row{"Id", "Connector", "Replicas",
		color.GreenString("R"),
		color.YellowString("P"),
		color.RedString("!"),
		"S"})
	w.SetStyle(table.StyleLight)
	err := controlClient().ListDeployments(filter, func(deployment *client.DeploymentInfo, hasMore bool) error {
		if deployment == nil {
			return nil
		}

		running := ""
		if deployment.Status.Running > 0 {
			running = color.GreenString(fmt.Sprintf("%d", deployment.Status.Running))
		}

		pending := ""
		if deployment.Status.Pending > 0 {
			pending = color.YellowString(fmt.Sprintf("%d", deployment.Status.Pending))
		}

		failed := ""
		if deployment.Status.Errored > 0 {
			failed = color.RedString(fmt.Sprintf("%d", deployment.Status.Errored))
		}

		stopped := ""
		if deployment.Status.Stopped > 0 {
			stopped = fmt.Sprintf("%d", deployment.Status.Stopped)
		}

		w.AppendRow(table.Row{deployment.DeploymentId, deployment.ConnectorId, deployment.Replicas,
			running,
			pending,
			failed,
			stopped,
		})
		return nil
	})
	if err != nil {
		color.Red("Could not list deployments: %s", err)
		os.Exit(1)
	}

	fmt.Println(w.Render())
	return nil
}

func (c *deploymentCommand) describeDeployment(pc *fisk.ParseContext) error {
	if c.deploymentId == "" {
		c.deploymentId = "latest"
	}

	deployment, err := controlClient().GetDeployment(c.connectorId, c.deploymentId)
	if err != nil {
		return err
	}

	if deployment == nil {
		color.Red("Deployment %s not found for connector %s", c.deploymentId, c.connectorId)
		os.Exit(1)
	}

	cols := render.New("")
	cols.AddSectionTitle("General Information")
	cols.AddRow("Id", deployment.DeploymentId)
	cols.AddRow("ConnectorId", deployment.ConnectorId)
	cols.AddRow("Replicas", deployment.Replicas)

	if deployment.PlacementTags != nil {
		cols.AddRow("Placement Tags", strings.Join(deployment.PlacementTags, " "))
	}

	cols.AddSectionTitle("Status")
	cols.AddRow("Pending", color.BlueString("%d", deployment.Status.Pending))
	cols.AddRow("Running", color.GreenString("%d", deployment.Status.Running))
	cols.AddRow("Stopped", color.YellowString("%d", deployment.Status.Stopped))
	cols.AddRow("Failed", color.RedString("%d", deployment.Status.Errored))

	cols.AddSectionTitle("Instances")
	w := table.NewWriter()
	w.AppendHeader(table.Row{"Id", "Status"})
	w.SetStyle(table.StyleLight)

	instanceFilter := client.InstanceFilter{
		ConnectorId:  deployment.ConnectorId,
		DeploymentId: deployment.DeploymentId,
	}

	err = controlClient().ListInstances(instanceFilter, func(instance *client.InstanceInfo, hasMore bool) error {
		if instance == nil {
			return nil
		}

		status := string(instance.Status)
		switch instance.Status {
		case model.InstanceRunning:
			status = color.GreenString(status)
		case model.InstanceStopped:
			status = color.YellowString(status)
		case model.InstanceFailed:
			status = color.RedString(status)
		}

		w.AppendRow(table.Row{instance.InstanceId, status})

		return nil
	})

	if err != nil {
		color.Red("Could not list executions: %s", err)
	} else {
		cols.Println(w.Render())
	}

	return cols.Frender(os.Stdout)
}

func (c *deploymentCommand) purgeDeployments(pc *fisk.ParseContext) error {
	filter := client.DeploymentFilter{
		ConnectorId: c.connectorId,
	}

	prompt := color.YellowString("This will remove all stopped or failed deployments, are you sure?")
	if !confirmAction(prompt) {
		fmt.Println("Aborted")
		return nil
	}

	w := table.NewWriter()
	w.AppendHeader(table.Row{"Deployment", "Connector", ""})
	w.SetStyle(table.StyleLight)
	err := controlClient().PurgeDeployments(filter, func(pi *client.DeploymentPurgeInfo, hasMore bool) error {
		if pi == nil {
			return nil
		}

		outcome := color.GreenString("Purged")
		if pi.Error != "" {
			outcome = color.RedString(pi.Error)
		}

		w.AppendRow(table.Row{pi.DeploymentId, pi.ConnectorId, outcome})
		return nil
	})
	if err != nil {
		color.Red("Could not purge deployments: %s", err)
		os.Exit(1)
	}

	fmt.Println(w.Render())
	return nil
}

func (c *deploymentCommand) getDeploymentLogs(pc *fisk.ParseContext) error {
	logs, err := controlClient().GetLogs(c.connectorId, c.deploymentId, "")
	if err != nil {
		color.Red("Could not get logs for deployment %s: %s", c.deploymentId, err)
		os.Exit(1)
	}

	printLogs(logs)

	return nil
}

func (c *deploymentCommand) getDeploymentEvents(pc *fisk.ParseContext) error {
	events, err := controlClient().GetEvents(c.connectorId, c.deploymentId, "")
	if err != nil {
		color.Red("Could not get events for deployment %s: %s", c.deploymentId, err)
		os.Exit(1)
	}

	printEvents(events)

	return nil
}

func (c *deploymentCommand) getDeploymentMetrics(pc *fisk.ParseContext) error {
	metrics, err := controlClient().GetMetrics(c.connectorId, c.deploymentId, "")
	if err != nil {
		color.Red("Could not get metrics for deployment %s: %s", c.deploymentId, err)
		os.Exit(1)
	}

	printMetrics(metrics)

	return nil
}

func confirmAction(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s (yes/no): ", s)

	response, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	response = strings.ToLower(strings.TrimSpace(response))

	return response == "y" || response == "yes"
}
