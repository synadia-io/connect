package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/choria-io/fisk"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/synadia-io/connect/cli/render"
	"github.com/synadia-io/connect/client"
)

func init() {
	registerCommand("instance", 3, configureInstanceCommand)
}

type instanceCommand struct {
	connectorId    string
	deploymentId   string
	instanceId     string
	instanceStatus []string
}

func configureInstanceCommand(parentCmd commandHost) {
	c := &instanceCommand{}

	instanceCmd := parentCmd.Command("instance", "Interact with instances").Alias("i")

	listCmd := instanceCmd.Command("list", "List instances for a connector").Alias("ls").Action(c.listInstances)
	listCmd.Arg("connector", "The connector to list instances for").StringVar(&c.connectorId)
	listCmd.Arg("deployment", "The deployment to list instances for").StringVar(&c.deploymentId)

	infoCmd := instanceCmd.Command("get", "Show details about a instance").Alias("info").Action(c.describeInstance)
	infoCmd.Arg("connector", "The connector to show instance details for").Required().StringVar(&c.connectorId)
	infoCmd.Arg("deployment", "The deployment to show instance details for").Required().StringVar(&c.deploymentId)
	infoCmd.Arg("id", "The instance id to show details for").Required().StringVar(&c.instanceId)
}

func (c *instanceCommand) listInstances(pc *fisk.ParseContext) error {
	filter := client.InstanceFilter{
		ConnectorId:  c.connectorId,
		DeploymentId: c.deploymentId,
	}

	w := table.NewWriter()
	w.AppendHeader(table.Row{"Id", "Deployment", "Connector", "Status", "Uptime"})

	w.SetStyle(table.StyleLight)
	err := controlClient().ListInstances(filter, func(instance *client.InstanceInfo, hasMore bool) error {
		if instance == nil {
			return nil
		}

		w.AppendRow(table.Row{instance.InstanceId, instance.DeploymentId, instance.ConnectorId, instance.Status, instance.Uptime})
		return nil
	})
	if err != nil {
		color.Red("Could not list instances: %s", err)
		os.Exit(1)
	}

	fmt.Println(w.Render())
	return nil
}

func (c *instanceCommand) describeInstance(pc *fisk.ParseContext) error {
	instance, err := controlClient().GetInstance(c.connectorId, c.deploymentId, c.instanceId)
	if err != nil {
		return err
	}

	if instance == nil {
		color.Red("Instance %s not found for deployment %s and connector %s", c.instanceId, c.deploymentId, c.connectorId)
		os.Exit(1)
	}

	now := time.Now()
	loc := now.Location()

	cols := render.New("")
	cols.AddRow("Id", instance.InstanceId)
	cols.AddRow("DeploymentId", instance.DeploymentId)
	cols.AddRow("ConnectorId", instance.ConnectorId)

	cols.AddSectionTitle("Timeline")
	if instance.ScheduledAt != nil {
		cols.AddRow("Scheduled At", instance.ScheduledAt.In(loc))
	}
	if instance.StartedAt != nil {
		cols.AddRow("Started At", instance.StartedAt.In(loc))
	}
	if instance.FinishedAt != nil {
		cols.AddRow("Finished At", instance.FinishedAt.In(loc))
	}

	return cols.Frender(os.Stdout)
}
