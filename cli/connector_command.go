package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/choria-io/fisk"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/synadia-io/connect/client"
	"github.com/synadia-io/connect/model"
	terminal "golang.org/x/term"
	"gopkg.in/yaml.v3"
)

func init() {
	registerCommand("connector", 1, configureConnectorCommand)
}

type connectorCommand struct {
	listFilterKinds []string
	id              string
	format          string
	declaration     string
	interactive     bool
	placementTags   []string
	replicas        int
	pull            bool

	deploymentId string
	force        bool
}

func configureConnectorCommand(parentCmd commandHost) {
	c := &connectorCommand{}

	connectorCmd := parentCmd.Command("connector", "Manage connectors").Alias("c")

	list := connectorCmd.Command("list", "List all connectors").Alias("ls").Action(c.listConnectors)
	list.Flag("kind", "Filter by kind").Default().StringsVar(&c.listFilterKinds)

	info := connectorCmd.Command("get", "Show detailed information about a connector").Alias("info").Action(c.getConnector)
	info.Arg("id", "The id of the connector to describe").Required().StringVar(&c.id)
	info.Flag("format", "The format in which to render the information").Default("human").EnumVar(&c.format, "human", "json", "yaml")

	remove := connectorCmd.Command("remove", "Remove a connector").Alias("delete").Action(c.removeConnector)
	remove.Arg("id", "The id of the connector to remove").Required().StringVar(&c.id)

	createCmd := connectorCmd.Command("create", "Create a new connector").Alias("add").Action(c.createConnector)
	createCmd.Arg("id", "The id of the connector to create").Required().StringVar(&c.id)
	createCmd.Arg("declaration", "The connector declaration").Default("!nil!").StringVar(&c.declaration)
	createCmd.Flag("interactive", "Create/update the connector using your editor").Short('i').BoolVar(&c.interactive)

	editCmd := connectorCmd.Command("edit", "Edit a connector").Action(c.updateConnector)
	editCmd.Arg("id", "The id of the connector to edit").Required().StringVar(&c.id)

	deployCmd := connectorCmd.Command("deploy", "Deploy a connector").Action(c.deployConnector)
	deployCmd.Arg("id", "The id of the connector to deploy").Required().StringVar(&c.id)
	deployCmd.Flag("tag", "Placement tags").Short('t').StringsVar(&c.placementTags)
	deployCmd.Flag("replicas", "Number of connector instances").Short('r').Default("1").IntVar(&c.replicas)
	deployCmd.Flag("pull", "Indicate whether the image should be pulled").Default("true").BoolVar(&c.pull)

	redeployCmd := connectorCmd.Command("redeploy", "Redeploy a connector").Action(c.redeployConnector)
	redeployCmd.Arg("id", "The id of the connector to redeploy").Required().StringVar(&c.id)

	undeploy := connectorCmd.Command("undeploy", "Undeploy a connector").Action(c.undeployConnector)
	undeploy.Arg("id", "The id of the connector to undeploy").Required().StringVar(&c.id)
	undeploy.Flag("force", "Force undeployment").UnNegatableBoolVar(&c.force)

	logsCmd := connectorCmd.Command("logs", "List current logs for a connector").Alias("l").Action(c.getConnectorLogs)
	logsCmd.Arg("id", "The id of the connector to get logs for").Required().StringVar(&c.id)

	eventsCmd := connectorCmd.Command("events", "List current events for a connector").Alias("e").Action(c.getConnectorEvents)
	eventsCmd.Arg("id", "The id of the connector to get events for").Required().StringVar(&c.id)

	metricsCmd := connectorCmd.Command("metrics", "List current metrics for a connector").Alias("m").Action(c.getConnectorMetrics)
	metricsCmd.Arg("id", "The id of the connector to get metrics for").Required().StringVar(&c.id)
}

func (c *connectorCommand) listConnectors(pc *fisk.ParseContext) error {
	filter := client.ConnectorFilter{}
	if len(c.listFilterKinds) > 0 {
		for _, k := range c.listFilterKinds {
			filter.Kinds = append(filter.Kinds, model.ConnectorKind(k))
		}
	}

	w := table.NewWriter()
	w.AppendHeader(table.Row{"ConnectorId", "Kind", "Description", "Status"})
	w.SetStyle(table.StyleLight)
	err := controlClient().ListConnectors(filter, func(connector *client.ConnectorInfo, hasMore bool) error {
		if connector == nil {
			return nil
		}

		active := color.YellowString("inactive")
		if connector.IsActive {
			active = color.GreenString("active")
		}

		w.AppendRow(table.Row{connector.ConnectorId, connector.Kind, connector.Description, active})
		return nil
	})
	if err != nil {
		color.Red("Could not list connectors: %s", err)
		os.Exit(1)
	}

	fmt.Println(w.Render())
	return nil
}

func (c *connectorCommand) getConnector(pc *fisk.ParseContext) error {
	connector, err := controlClient().GetConnector(c.id)
	if err != nil {
		return err
	}

	if connector == nil {
		color.Red("Connector %s not found", c.id)
		os.Exit(1)
	}

	switch c.format {
	case "json":
		RenderJsonConnector(connector)
	case "yaml":
		RenderYamlConnector(connector)
	default:
		RenderHumanConnector(connector)
	}

	return nil
}

func (c *connectorCommand) removeConnector(pc *fisk.ParseContext) error {
	deleted, err := controlClient().DeleteConnector(c.id)
	if err != nil {
		color.Red("Could not remove connector %s: %s", c.id, err)
		os.Exit(1)
	}

	if !deleted {
		color.Yellow("Connector %s not removed", c.id)
	} else {
		color.Green("Connector %s removed", c.id)
	}

	return nil
}

func (c *connectorCommand) createConnector(pc *fisk.ParseContext) error {
	if c.interactive {
		c.interactiveCreate(c.id)
		return nil
	} else if c.declaration == "!nil!" && (terminal.IsTerminal(int(os.Stdout.Fd()))) {
		log.Println("Reading payload from STDIN")
		body, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		c.declaration = string(body)
	}

	var payload model.ConnectorConfig
	// Use yaml.Unmarshal to support loading both YAML and JSON as input.
	if err := yaml.Unmarshal([]byte(c.declaration), &payload); err != nil {
		color.Red("Could not parse connector configuration: %s", err)
		os.Exit(1)
	}

	connector, err := controlClient().CreateConnector(c.id, payload)
	if err != nil {
		color.Red("Could not save connector: %s", err)
		os.Exit(1)
	}

	fmt.Printf("Created connector %s\n", color.GreenString(connector.Id))

	return nil
}

func (c *connectorCommand) updateConnector(pc *fisk.ParseContext) error {
	connector, err := controlClient().GetConnector(c.id)
	if err != nil {
		return err
	}

	if connector == nil {
		color.Red("Connector %s not found", c.id)
		os.Exit(1)
	}

	ccfg, changed := c.interactiveEdit(*connector)
	if !changed {
		color.Red("no difference in configuration")
		os.Exit(1)
	}

	if _, err := controlClient().PatchConnector(connector.Id, ccfg); err != nil {
		color.Red("Could not update the connector: %s", err)
		os.Exit(1)
	}

	color.Green("Connector %s updated", connector.Id)
	return nil
}

func (c *connectorCommand) deployConnector(pc *fisk.ParseContext) error {
	o := []client.DeployOpt{
		client.WithDeployReplicas(c.replicas),
		client.WithDeployTimeout(opts().Timeout.String()),
		client.WithPull(c.pull),
	}

	if c.placementTags != nil {
		o = append(o, client.WithDeployPlacementTags(c.placementTags))
	}

	depl, err := controlClient().DeployConnector(c.id, o...)
	if err != nil {
		color.Red("Could not deploy connector %s: %s", c.id, err)
		os.Exit(1)
	}

	if depl.Error != "" {
		color.Red("Deploy error: %s", depl.Error)
		os.Exit(1)
	}

	if len(depl.InstanceErrors) > 0 {
		color.Red("Deploy instance errors:")
		for k, v := range depl.InstanceErrors {
			fmt.Printf("- %s: %s\n", k, v)
		}
	}

	fmt.Printf("Deployed connector %s as %s\n", color.GreenString(c.id), color.GreenString(depl.DeploymentId))

	return nil
}

func (c *connectorCommand) redeployConnector(pc *fisk.ParseContext) error {
	o := []client.RedeployOpt{
		client.WithRedeployReplicas(c.replicas),
		client.WithRedeployTimeout(opts().Timeout.String()),
	}

	if c.placementTags != nil {
		o = append(o, client.WithRedeployPlacementTags(c.placementTags))
	}

	depl, err := controlClient().RedeployConnector(c.id, o...)
	if err != nil {
		color.Red("Could not redeploy connector %s: %s", c.id, err)
		os.Exit(1)
	}

	if depl.Undeploy.Error != "" {
		color.Red("Undeploy error: %s", depl.Undeploy.Error)
		os.Exit(1)
	}

	if len(depl.Undeploy.InstanceErrors) > 0 {
		color.Red("Undeploy instance errors:")
		for k, v := range depl.Undeploy.InstanceErrors {
			fmt.Printf("- %s: %s\n", k, v)
		}
	}

	if depl.Deploy.Error != "" {
		color.Red("Deploy error: %s", depl.Deploy.Error)
		os.Exit(1)
	}

	if len(depl.Deploy.InstanceErrors) > 0 {
		color.Red("Deploy instance errors:")
		for k, v := range depl.Deploy.InstanceErrors {
			fmt.Printf("- %s: %s\n", k, v)
		}
	}

	fmt.Printf("Redeployed %s as %s\n", color.GreenString(c.id), color.GreenString(depl.DeploymentId))

	return nil
}

func (c *connectorCommand) undeployConnector(pc *fisk.ParseContext) error {
	o := []client.UndeployOpt{
		client.WithUndeployTimeout(opts().Timeout.String()),
	}

	if c.force {
		o = append(o, client.WithUndeployForce(true))
	}

	depl, err := controlClient().UndeployConnector(c.id, o...)
	if err != nil {
		color.Red("Could not undeploy connector %s: %s", c.id, err)
		os.Exit(1)
	}

	if depl.Error != "" {
		color.Red("Undeploy error: %s", depl.Error)
		os.Exit(1)
	}

	if len(depl.InstanceErrors) > 0 {
		color.Red("Undeploy instance errors:")
		for k, v := range depl.InstanceErrors {
			fmt.Printf("- %s: %s\n", k, v)
		}
	}

	fmt.Printf("Undeployed %s\n", color.GreenString(c.id))

	return nil
}

func (c *connectorCommand) getConnectorLogs(pc *fisk.ParseContext) error {
	logs, err := controlClient().GetConnectorLogs(c.id, "", "")
	if err != nil {
		color.Red("Could not get logs for connector %s: %s", c.id, err)
		os.Exit(1)
	}

	printLogs(logs)

	return nil
}

func (c *connectorCommand) getConnectorEvents(pc *fisk.ParseContext) error {
	events, err := controlClient().GetConnectorEvents(c.id, "", "")
	if err != nil {
		color.Red("Could not get events for connector %s: %s", c.id, err)
		os.Exit(1)
	}

	printEvents(events)

	return nil
}

func (c *connectorCommand) getConnectorMetrics(pc *fisk.ParseContext) error {
	metrics, err := controlClient().GetConnectorMetrics(c.id, "", "")
	if err != nil {
		color.Red("Could not get metrics for connector %s: %s", c.id, err)
		os.Exit(1)
	}

	printMetrics(metrics)

	return nil
}

func (c *connectorCommand) interactiveEdit(connector model.Connector) (model.ConnectorConfig, bool) {
	configYml, err := yaml.Marshal(connector.ConnectorConfig)
	if err != nil {
		color.Red("could not serialize yaml file: %s", err)
	}

	tmpFile, err := os.CreateTemp("", "*.yaml")
	if err != nil {
		color.Red("could not create temporary file: %s", err)
		os.Exit(1)
	}
	defer os.Remove(tmpFile.Name())

	_, err = fmt.Fprint(tmpFile, string(configYml))
	if err != nil {
		color.Red("could not create temporary file: %s", err)
		os.Exit(1)
	}
	tmpFile.Close()

	err = editFile(tmpFile.Name())
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

	modifiedConfig, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

	var payload model.ConnectorConfig
	// Use yaml.Unmarshal to support loading both YAML and JSON as input.
	if err := yaml.Unmarshal(modifiedConfig, &payload); err != nil {
		color.Red("Could not parse the connector configuration: %s", err)
		os.Exit(1)
	}

	return payload, !bytes.Equal(configYml, modifiedConfig)
}

func (c *connectorCommand) interactiveCreate(id string) {
	cfg, err := c.selectConnectorTemplate()
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}

	ccfg, _ := c.interactiveEdit(model.Connector{
		ConnectorConfig: *cfg,
		Id:              c.id,
	})

	if _, err := controlClient().CreateConnector(id, ccfg); err != nil {
		color.Red("Could not create the connector: %s", err)
		os.Exit(1)
	}

	color.Green("Connector %s created", id)
}

func (c *connectorCommand) selectConnectorTemplate() (*model.ConnectorConfig, error) {
	version, err := libraryClient().GetLatestVersion("vanilla")
	if err != nil {
		return nil, fmt.Errorf("could not get the latest vanilla version: %s", err)
	}

	cfg := model.ConnectorConfig{
		Description: "A summary of what this connector does",
		Workload:    version.Workload.Location,
	}

	if version.Workload.Metrics != nil {
		cfg.Metrics = &model.MetricsEndpoint{
			Port: version.Workload.Metrics.Port,
			Path: version.Workload.Metrics.Path,
		}
	}

	templates := []struct {
		name   string
		modify func()
	}{
		{name: "Inlet - Reading from a Source to NATS", modify: func() {
			cfg.Steps = &model.Steps{
				Source:   &model.Source{},
				Producer: &model.Producer{},
			}
		}},
		{name: "Outlet - Writing from NATS to a Sink", modify: func() {
			cfg.Steps = &model.Steps{
				Consumer: &model.Consumer{},
				Sink:     &model.Sink{},
			}
		}},
	}

	var options []string
	mapping := make(map[string]func())
	for _, template := range templates {
		options = append(options, template.name)
		mapping[template.name] = template.modify
	}

	choice := ""
	err = survey.AskOne(&survey.Select{
		Message: "Connector Template",
		Options: options,
	}, &choice, survey.WithValidator(survey.Required))
	if err != nil {
		return nil, err
	}

	modify, ok := mapping[choice]
	if !ok {
		return nil, fmt.Errorf("template not found")
	}
	modify()

	return &cfg, nil
}
