package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/synadia-io/connect/cli/render"
	"github.com/synadia-io/connect/client"
	"github.com/synadia-io/connect/model"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"gopkg.in/yaml.v3"
)

func RenderHumanConnector(v *model.Connector) {
	//var err error
	cols := render.New("")
	cols.Println()
	cols.Println(v.Description)
	cols.Println()

	cols.AddSectionTitle("Properties")
	cols.AddRow("Id", v.Id)
	cols.AddRow("Workload", v.Workload)

	if len(v.DeploymentIds) > 0 {
		cols.AddSectionTitle("Deployments")

		w := table.NewWriter()
		w.SetStyle(table.StyleRounded)

		w.AppendHeader(table.Row{"Id", "Name", "Status"})
		err := controlClient().ListDeployments(client.DeploymentFilter{
			ConnectorId: v.Id,
		}, func(d *client.DeploymentInfo, hasMore bool) error {
			if d == nil {
				return nil
			}

			var deploymentStatus string
			if d.Status == nil {
				deploymentStatus = ""
			} else {
				deploymentStatus = fmt.Sprintf("%s %s %s %s",
					color.BlueString(fmt.Sprintf("%d", d.Status.Pending)),
					color.GreenString(fmt.Sprintf("%d", d.Status.Running)),
					color.YellowString(fmt.Sprintf("%d", d.Status.Stopped)),
					color.RedString(fmt.Sprintf("%d", d.Status.Errored)),
				)
			}

			w.AppendRow(table.Row{d.DeploymentId, d.Replicas, deploymentStatus})

			return nil
		})

		if err != nil {
			cols.Println(color.RedString(err.Error()))
		}

		cols.Println(w.Render())
	}

	cols.Println()
	cols.Frender(os.Stdout)
}

func RenderJsonConnector(c *model.Connector) {
	b, err := json.Marshal(c)
	if err != nil {
		color.Red(err.Error())
		return
	}
	fmt.Println(string(b))
}

func RenderYamlConnector(c *model.Connector) {
	b, err := yaml.Marshal(c)
	if err != nil {
		color.Red(err.Error())
		return
	}
	fmt.Println(string(b))
}

func RenderEvent(item model.InstanceEvent) {
	line := strings.ToUpper(string(item.Type))

	if item.Error != "" {
		line = color.RedString(fmt.Sprintf("%s %s", line, item.Error))
	}

	fmt.Printf("%s %s %s %s %s\n",
		item.Timestamp.Format("2006/01/02 15:04:05"),
		item.ConnectorId,
		item.DeploymentId[len(item.DeploymentId)-6:],
		item.InstanceId[len(item.InstanceId)-6:],
		line)
}

func RenderLog(item model.InstanceLog) {
	l := strings.SplitN(item.Line, " ", 2)[1]
	if strings.Contains(l, "ERROR") {
		l = color.RedString(l)
	} else if strings.Contains(l, "WARN") {
		l = color.YellowString(l)
	} else if strings.Contains(l, "DEBUG") {
		l = color.BlueString(l)
	}

	fmt.Printf("%s %s %s %s %s\n",
		item.Timestamp.Format("2006/01/02 15:04:05"),
		item.ConnectorId,
		item.DeploymentId[len(item.DeploymentId)-6:],
		item.InstanceId[len(item.InstanceId)-6:],
		l)
}

func RenderMetric(item model.InstanceMetric) {
	// TODO: in the future we might want to do something fancier here
	fmt.Printf("%s\n", item.Data)
}
