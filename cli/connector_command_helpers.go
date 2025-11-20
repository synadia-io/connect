package cli

import (
	"fmt"
	"time"

	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/synadia-io/connect/v2/model"
)

// These are testable helper functions that can be called with a provided AppContext

func (c *connectorCommand) listConnectorsWithClient(appCtx *AppContext) error {
	resp, err := appCtx.Client.ListConnectors(c.opts.Timeout)
	if err != nil {
		return fmt.Errorf("failed to list connectors: %w", err)
	}

	if len(resp) == 0 {
		fmt.Println("No connectors found")
		return nil
	}

	tbl := table.NewWriter()
	tbl.SetStyle(table.StyleRounded)
	title := "Connectors"
	tbl.SetTitle(title)
	tbl.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, Name: "Name"},
		{Number: 2, Name: "Description"},
		{Number: 3, Name: "Runtime"},
		{Number: 4, Name: color.GreenString("\u25B6"), WidthMin: 3, WidthMax: 5, Align: text.AlignCenter, AlignHeader: text.AlignCenter},
		{Number: 5, Name: color.RedString("\u25FC"), WidthMin: 3, WidthMax: 5, Align: text.AlignCenter, AlignHeader: text.AlignCenter},
	})

	tbl.AppendHeader(table.Row{"Name", "Description", "Runtime", color.GreenString("\u25B6"), color.RedString("\u25FC")}, table.RowConfig{AutoMerge: true})

	for _, c := range resp {
		running := ""
		if c.Instances.Running > 0 {
			running = color.GreenString("%d", c.Instances.Running)
		}

		stopped := ""
		if c.Instances.Stopped > 0 {
			stopped = color.RedString("%d", c.Instances.Stopped)
		}

		tbl.AppendRow(table.Row{
			c.ConnectorId,
			text.WrapSoft(c.Description, 50),
			c.RuntimeId,
			running,
			stopped,
		})
	}

	fmt.Println(tbl.Render())
	return nil
}

func (c *connectorCommand) getConnectorWithClient(appCtx *AppContext) error {
	connector, err := appCtx.Client.GetConnector(c.id, c.opts.Timeout)
	if err != nil {
		return fmt.Errorf("failed to get connector: %w", err)
	}

	fmt.Println(renderConnector(*connector))
	return nil
}

func (c *connectorCommand) removeConnectorWithClient(appCtx *AppContext) error {
	err := appCtx.Client.DeleteConnector(c.id, c.opts.Timeout)
	if err != nil {
		return fmt.Errorf("failed to delete connector: %w", err)
	}

	fmt.Printf("Connector %s deleted\n", c.id)
	return nil
}

func (c *connectorCommand) connectorStatusWithClient(appCtx *AppContext) error {
	status, err := appCtx.Client.GetConnectorStatus(c.id, c.opts.Timeout)
	if err != nil {
		return fmt.Errorf("failed to get connector status: %w", err)
	}

	// Display status
	fmt.Printf("Connector %s status:\n", c.id)
	fmt.Printf("  Running: %d\n", status.Running)
	fmt.Printf("  Stopped: %d\n", status.Stopped)

	return nil
}

func (c *connectorCommand) startConnectorWithClient(appCtx *AppContext) error {
	// Validate timeout format
	_, err := time.ParseDuration(c.startTimeout)
	if err != nil {
		return fmt.Errorf("invalid timeout: %w", err)
	}

	envVars := make(model.ConnectorStartOptionsEnvVars)
	for k, v := range c.envVars {
		envVars[k] = v
	}

	startOpts := &model.ConnectorStartOptions{
		Pull:          !c.noPull,
		Replicas:      c.replicas,
		Timeout:       c.startTimeout,
		EnvVars:       envVars,
		PlacementTags: c.placementTags,
	}

	instances, err := appCtx.Client.StartConnector(c.id, startOpts, c.opts.Timeout)
	if err != nil {
		return fmt.Errorf("failed to start connector: %w", err)
	}

	fmt.Printf("Started %d instances of connector %s\n", len(instances), c.id)
	return nil
}

func (c *connectorCommand) stopConnectorWithClient(appCtx *AppContext) error {
	instances, err := appCtx.Client.StopConnector(c.id, c.opts.Timeout)
	if err != nil {
		return fmt.Errorf("failed to stop connector: %w", err)
	}

	fmt.Printf("Stopped %d instances of connector %s\n", len(instances), c.id)
	return nil
}

func (c *connectorCommand) copyConnectorWithClient(appCtx *AppContext) error {
	// Get the source connector
	connector, err := appCtx.Client.GetConnector(c.id, c.opts.Timeout)
	if err != nil {
		return fmt.Errorf("failed to get source connector: %w", err)
	}

	// Create the copy
	copied, err := appCtx.Client.CreateConnector(c.targetId, connector.Description, connector.RuntimeId, c.runtimeVersion, connector.Steps, c.opts.Timeout)
	if err != nil {
		return fmt.Errorf("failed to create connector copy: %w", err)
	}

	fmt.Printf("Copied connector %s to %s\n", c.id, copied.ConnectorId)
	return nil
}

func (c *connectorCommand) reloadConnectorWithClient(appCtx *AppContext) error {
	// Stop the connector
	_, err := appCtx.Client.StopConnector(c.id, c.opts.Timeout)
	if err != nil {
		return fmt.Errorf("failed to stop connector: %w", err)
	}

	// Start it again with default options
	startOpts := &model.ConnectorStartOptions{
		Pull:     true,
		Replicas: 1,
	}

	_, err = appCtx.Client.StartConnector(c.id, startOpts, c.opts.Timeout)
	if err != nil {
		return fmt.Errorf("failed to start connector: %w", err)
	}

	fmt.Printf("Reloaded connector %s\n", c.id)
	return nil
}
