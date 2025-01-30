package cli

import (
	"encoding/json"
	"fmt"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/synadia-io/connect/client"
	"strings"

	"github.com/synadia-io/connect/model"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"gopkg.in/yaml.v3"
)

func RenderHumanConnector(v *model.Connector) {
	fmt.Println(text.Colors{text.Underline, text.Bold}.Sprintf("%q %s", v.Id, v.Steps.Kind()))
	fmt.Println()

	if v.Description != "" {
		fmt.Println(v.Description)
		fmt.Println()
	}

	if v.Steps != nil {
		b, _ := yaml.Marshal(v.Steps)
		fmt.Println(string(b))
	}

	props := table.NewWriter()
	props.SetStyle(table.StyleRounded)
	props.AppendRow(table.Row{"Properties", "Properties"}, table.RowConfig{AutoMerge: true})
	props.AppendSeparator()
	props.AppendRow(table.Row{"Key", "Value"})
	props.AppendSeparator()

	metrics := "Disabled"
	if v.Metrics != nil {
		metrics = fmt.Sprintf(":%d%s", v.Metrics.Port, v.Metrics.Path)
	}
	props.AppendRow(table.Row{"Metrics", metrics})
	props.AppendRow(table.Row{"Workload", v.Workload})
	fmt.Println(props.Render())
	fmt.Println()

	if len(v.DeploymentIds) > 0 {
		tbl := table.NewWriter()
		tbl.SetStyle(table.StyleRounded)
		tbl.AppendRow(table.Row{"Deployments", "Deployments", "Deployments"}, table.RowConfig{AutoMerge: true})
		tbl.AppendSeparator()

		statusHeader := fmt.Sprintf("%s %s %s %s",
			color.BlueString("P"),
			color.GreenString("R"),
			color.YellowString("S"),
			color.RedString("!"),
		)
		tbl.AppendRow(table.Row{"Deployment ID", "Total", statusHeader})
		tbl.AppendSeparator()

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

			tbl.AppendRow(table.Row{d.DeploymentId, d.Replicas, deploymentStatus})

			return nil
		})

		if err != nil {
			errStr := color.RedString(err.Error())
			tbl.AppendRow(table.Row{errStr, errStr, errStr})
		}

		fmt.Println(tbl.Render())
	}

	////var err error
	//cols := render.New("")
	//cols.Println()
	//cols.Println(v.Description)
	//cols.Println()
	//
	//cols.AddSectionTitle("Properties")
	//cols.AddRow("Id", v.Id)
	//cols.AddRow("Workload", v.Workload)

	//if len(v.DeploymentIds) > 0 {
	//	cols.AddSectionTitle("Deployments")
	//
	//	w := table.NewWriter()
	//	w.SetStyle(table.StyleRounded)
	//
	//	w.AppendHeader(table.Row{"Id", "Name", "Status"})
	//	err := controlClient().ListDeployments(client.DeploymentFilter{
	//		ConnectorId: v.Id,
	//	}, func(d *client.DeploymentInfo, hasMore bool) error {
	//		if d == nil {
	//			return nil
	//		}
	//
	//		var deploymentStatus string
	//		if d.Status == nil {
	//			deploymentStatus = ""
	//		} else {
	//			deploymentStatus = fmt.Sprintf("%s %s %s %s",
	//				color.BlueString(fmt.Sprintf("%d", d.Status.Pending)),
	//				color.GreenString(fmt.Sprintf("%d", d.Status.Running)),
	//				color.YellowString(fmt.Sprintf("%d", d.Status.Stopped)),
	//				color.RedString(fmt.Sprintf("%d", d.Status.Errored)),
	//			)
	//		}
	//
	//		w.AppendRow(table.Row{d.DeploymentId, d.Replicas, deploymentStatus})
	//
	//		return nil
	//	})
	//
	//	if err != nil {
	//		cols.Println(color.RedString(err.Error()))
	//	}
	//
	//	cols.Println(w.Render())
	//}
	//
	//cols.Println()
	//cols.Frender(os.Stdout)
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

func RenderEvent(item model.InstanceEvent, detailed bool) {
	line := strings.ToUpper(string(item.Type))

	if item.Error != "" {
		line = color.RedString(fmt.Sprintf("%s %s", line, item.Error))
	}

	if detailed {
		fmt.Printf("%s %s %s %s %s\n",
			item.Timestamp.Format("2006/01/02 15:04:05"),
			item.ConnectorId,
			item.DeploymentId[len(item.DeploymentId)-6:],
			item.InstanceId[len(item.InstanceId)-6:],
			line)
	} else {
		fmt.Printf("%s %s %s\n",
			item.ConnectorId,
			item.InstanceId[len(item.InstanceId)-6:],
			line)
	}
}

func RenderLog(item model.InstanceLog, detailed bool) {
	l := strings.SplitN(item.Line, " ", 2)[1]
	if strings.Contains(l, "ERROR") {
		l = color.RedString(l)
	} else if strings.Contains(l, "WARN") {
		l = color.YellowString(l)
	} else if strings.Contains(l, "DEBUG") {
		l = color.BlueString(l)
	}

	if detailed {
		fmt.Printf("%s %s %s %s %s\n",
			item.Timestamp.Format("2006/01/02 15:04:05"),
			item.ConnectorId,
			item.DeploymentId[len(item.DeploymentId)-6:],
			item.InstanceId[len(item.InstanceId)-6:],
			l)
	} else {
		fmt.Printf("%s %s %s\n",
			item.ConnectorId,
			item.InstanceId[len(item.InstanceId)-6:],
			l)
	}
}

func RenderMetric(item model.InstanceMetric) {
	// TODO: in the future we might want to do something fancier here
	fmt.Printf("%s\n", item.Data)
}
