package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/choria-io/fisk"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/synadia-io/connect/v2/convert"
	"github.com/synadia-io/connect/v2/docker"
	"github.com/synadia-io/connect/v2/standalone"
	"github.com/synadia-io/connect/v2/validation"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type standaloneCommand struct {
	opts *Options

	// Common flags
	connectorName string // Standardized connector name parameter

	// Run command flags
	image            string
	envVars          map[string]string
	dockerOpts       string
	envFile          string
	envFileSetByUser bool
	follow           bool
	remove           bool

	// Template flags
	templateName string
	outputFile   string

	// Runtime flags
	runtimeID          string
	runtimeRegistry    string
	runtimeName        string
	runtimeDescription string
	runtimeAuthor      string
}

func ConfigureStandaloneCommand(parentCmd commandHost, opts *Options) {
	c := &standaloneCommand{
		opts:    opts,
		envVars: make(map[string]string),
	}

	standaloneCmd := parentCmd.Command("standalone", "Run connectors in standalone mode without NATS services")

	// Validate command
	validateCmd := standaloneCmd.Command("validate", "Validate a connector definition").Action(c.validateConnector)
	validateCmd.Arg("name", "Connector name (will look for <name>.connector.yml)").Required().StringVar(&c.connectorName)

	// Run command
	runCmd := standaloneCmd.Command("run", "Run a connector locally using Docker").Action(c.runConnector)
	runCmd.Arg("name", "Connector name (will look for <name>.connector.yml)").Required().StringVar(&c.connectorName)
	runCmd.Flag("image", "Override Docker image (uses runtime configuration by default)").StringVar(&c.image)
	runCmd.Flag("env", "Environment variables to set").Short('e').StringMapVar(&c.envVars)
	runCmd.Flag("docker-opts", "Custom docker options to set").StringVar(&c.dockerOpts)
	runCmd.Flag("env-file", "Read environment variables from file").Default(".env").IsSetByUser(&c.envFileSetByUser).StringVar(&c.envFile)
	runCmd.Flag("follow", "Follow logs after starting").Short('f').BoolVar(&c.follow)
	runCmd.Flag("rm", "Remove container when it exits").BoolVar(&c.remove)

	// Stop command
	stopCmd := standaloneCmd.Command("stop", "Stop a running connector").Action(c.stopConnector)
	stopCmd.Arg("name", "Connector name").Required().StringVar(&c.connectorName)

	// Remove command
	removeCmd := standaloneCmd.Command("remove", "Remove a connector (stop and delete container)").Alias("rm").Action(c.removeConnector)
	removeCmd.Arg("name", "Connector name").Required().StringVar(&c.connectorName)

	// Logs command
	logsCmd := standaloneCmd.Command("logs", "Show logs for a connector").Action(c.showLogs)
	logsCmd.Arg("name", "Connector name").Required().StringVar(&c.connectorName)
	logsCmd.Flag("follow", "Follow logs").Short('f').BoolVar(&c.follow)

	// List command
	standaloneCmd.Command("list", "List running connectors").Alias("ls").Action(c.listConnectors)

	// Create command
	createCmd := standaloneCmd.Command("create", "Create a new connector definition file").Action(c.createConnector)
	createCmd.Arg("name", "Connector name (will create <name>.connector.yml)").Required().StringVar(&c.connectorName)
	createCmd.Flag("file", "Override output file path").StringVar(&c.outputFile)
	createCmd.Flag("template", "Template to use").StringVar(&c.templateName)

	// Template subcommands
	templateCmd := standaloneCmd.Command("template", "Manage connector templates")

	templateCmd.Command("list", "List available templates").Alias("ls").Action(c.listTemplates)

	templateGetCmd := templateCmd.Command("get", "Get a specific template").Action(c.getTemplate)
	templateGetCmd.Arg("name", "Template name").Required().StringVar(&c.templateName)
	templateGetCmd.Flag("output", "Output file").Short('o').Default("./ConnectFile.yaml").StringVar(&c.outputFile)

	// Runtime subcommands
	runtimeCmd := standaloneCmd.Command("runtime", "Manage connector runtimes")

	runtimeCmd.Command("list", "List available runtimes").Alias("ls").Action(c.listRuntimes)

	runtimeAddCmd := runtimeCmd.Command("add", "Add a new runtime").Action(c.addRuntime)
	runtimeAddCmd.Arg("id", "Runtime ID").Required().StringVar(&c.runtimeID)
	runtimeAddCmd.Arg("registry", "Registry path (without version tag)").Required().StringVar(&c.runtimeRegistry)
	runtimeAddCmd.Flag("name", "Runtime name").StringVar(&c.runtimeName)
	runtimeAddCmd.Flag("description", "Runtime description").StringVar(&c.runtimeDescription)
	runtimeAddCmd.Flag("author", "Runtime author").StringVar(&c.runtimeAuthor)

	runtimeRemoveCmd := runtimeCmd.Command("remove", "Remove a runtime").Alias("rm")
	runtimeRemoveCmd.Arg("id", "Runtime ID").Required().StringVar(&c.runtimeID)
	runtimeRemoveCmd.Action(c.removeRuntime)

	runtimeShowCmd := runtimeCmd.Command("show", "Show runtime details").Action(c.showRuntime)
	runtimeShowCmd.Arg("id", "Runtime ID").Required().StringVar(&c.runtimeID)
}

// Helper functions for consistent naming
func (c *standaloneCommand) getFilePath() string {
	if c.outputFile != "" {
		return c.outputFile
	}
	return fmt.Sprintf("%s.connector.yml", c.connectorName)
}

func (c *standaloneCommand) getContainerName() string {
	return fmt.Sprintf("%s_connector", c.connectorName)
}

func (c *standaloneCommand) validateConnector(pc *fisk.ParseContext) error {
	validator := validation.NewValidator()
	filePath := c.getFilePath()

	fmt.Printf("Validating connector '%s' (file: %s)\n", c.connectorName, filePath)

	// Check file exists
	if err := validator.ValidateFileExists(filePath); err != nil {
		return err
	}

	// Check file extension
	if err := validator.ValidateFileExtension(filePath); err != nil {
		return err
	}

	// Validate connector definition
	if err := validator.ValidateConnectorFile(filePath); err != nil {
		color.Red("✗ Validation failed: %s", err.Error())
		return err
	}

	color.Green("✓ Connector '%s' is valid", c.connectorName)
	return nil
}

func (c *standaloneCommand) runConnector(pc *fisk.ParseContext) error {
	// First validate the connector
	if err := c.validateConnector(pc); err != nil {
		return err
	}

	filePath := c.getFilePath()
	containerName := c.getContainerName()

	// Load environment variables from file if specified
	envVars, err := LoadEnvFile(c.envFile, c.envFileSetByUser)
	if err != nil {
		return fmt.Errorf("failed to load env file: %w", err)
	}

	// Merge with command line env vars
	for k, v := range envVars {
		if _, exists := c.envVars[k]; !exists {
			c.envVars[k] = v
		}
	}

	// Load and parse the connector file
	connector, err := c.loadConnectorSpec(filePath)
	if err != nil {
		return fmt.Errorf("failed to load connector spec: %w", err)
	}

	// Determine image to use
	image := c.image
	if image == "" {
		// Resolve runtime reference to full image name
		rm := standalone.NewRuntimeManager()
		resolvedImage, err := rm.ResolveRuntimeImage(connector.RuntimeId)
		if err != nil {
			return fmt.Errorf("failed to resolve runtime '%s': %w", connector.RuntimeId, err)
		}
		image = resolvedImage
		fmt.Printf("Using runtime '%s' resolved to image '%s'\n", connector.RuntimeId, image)
	}

	// Convert spec to model steps
	steps := convert.ConvertStepsFromSpec(connector.Steps)

	// Create Docker runner
	runner := docker.NewRunner()

	// Validate Docker is available
	if err := runner.ValidateDockerAvailable(); err != nil {
		return fmt.Errorf("docker is not available: %w", err)
	}

	// Run the connector
	runOpts := &docker.RunOptions{
		ConnectorID: containerName,
		Image:       image,
		Steps:       steps,
		EnvVars:     c.envVars,
		DockerOpts:  c.dockerOpts,
		Follow:      c.follow,
		Remove:      c.remove,
		RuntimeID:   connector.RuntimeId,
	}

	fmt.Printf("Starting connector '%s' with image '%s'\n", c.connectorName, image)

	if err := runner.Run(context.Background(), runOpts); err != nil {
		return fmt.Errorf("failed to run connector: %w", err)
	}

	if !c.follow {
		color.Green("✓ Connector '%s' started successfully", c.connectorName)
		fmt.Printf("Use 'connect standalone logs %s' to view logs\n", c.connectorName)
		fmt.Printf("Use 'connect standalone stop %s' to stop the connector\n", c.connectorName)
	}

	return nil
}

func (c *standaloneCommand) stopConnector(pc *fisk.ParseContext) error {
	runner := docker.NewRunner()
	containerName := c.getContainerName()

	fmt.Printf("Stopping connector '%s'\n", c.connectorName)

	if err := runner.Stop(context.Background(), containerName); err != nil {
		return fmt.Errorf("failed to stop connector: %w", err)
	}

	color.Green("✓ Connector '%s' stopped", c.connectorName)
	return nil
}

func (c *standaloneCommand) showLogs(pc *fisk.ParseContext) error {
	runner := docker.NewRunner()
	containerName := c.getContainerName()

	fmt.Printf("Showing logs for connector '%s'\n", c.connectorName)

	return runner.Logs(context.Background(), containerName, c.follow)
}

func (c *standaloneCommand) removeConnector(pc *fisk.ParseContext) error {
	runner := docker.NewRunner()
	containerName := c.getContainerName()

	fmt.Printf("Removing connector '%s'\n", c.connectorName)

	if err := runner.Remove(context.Background(), containerName); err != nil {
		return fmt.Errorf("failed to remove connector: %w", err)
	}

	color.Green("✓ Connector '%s' removed", c.connectorName)
	return nil
}

func (c *standaloneCommand) listConnectors(pc *fisk.ParseContext) error {
	runner := docker.NewRunner()

	return runner.List(context.Background())
}

func (c *standaloneCommand) createConnector(pc *fisk.ParseContext) error {
	fmt.Printf("Creating new connector definition: %s\n", c.connectorName)

	filePath := c.getFilePath()

	// Use existing template system
	template, err := c.selectTemplate()
	if err != nil {
		return fmt.Errorf("failed to select template: %w", err)
	}

	// Customize the template
	template.Description = fmt.Sprintf("Connector: %s", c.connectorName)
	// Keep the original runtime ID from the template

	// Write to file
	if err := c.writeConnectorFile(template, filePath); err != nil {
		return fmt.Errorf("failed to write connector file: %w", err)
	}

	color.Green("✓ Created connector definition: %s", filePath)
	fmt.Printf("Edit the file to customize your connector, then run:\n")
	fmt.Printf("  connect standalone validate %s\n", c.connectorName)
	fmt.Printf("  connect standalone run %s\n", c.connectorName)

	return nil
}

func (c *standaloneCommand) listTemplates(pc *fisk.ParseContext) error {
	fmt.Println("Available connector templates:")

	for _, template := range standaloneTemplates {
		fmt.Printf("  %s\t\t#%s\n", template.Description, template.ConnectorSpec.Description)
	}

	return nil
}

func (c *standaloneCommand) getTemplate(pc *fisk.ParseContext) error {
	var selectedTemplate *connectorTemplate

	for _, template := range standaloneTemplates {
		if strings.Contains(strings.ToLower(template.Description), strings.ToLower(c.templateName)) {
			selectedTemplate = &template
			break
		}
	}

	if selectedTemplate == nil {
		return fmt.Errorf("template '%s' not found", c.templateName)
	}

	// Copy template and customize
	spec := selectedTemplate.ConnectorSpec
	// Keep the original runtime ID from the template

	if err := c.writeConnectorFile(&spec, c.outputFile); err != nil {
		return fmt.Errorf("failed to write template: %w", err)
	}

	color.Green("✓ Created template file: %s", c.outputFile)
	return nil
}

func (c *standaloneCommand) listRuntimes(pc *fisk.ParseContext) error {
	rm := standalone.NewRuntimeManager()
	runtimes, err := rm.LoadRuntimes()
	if err != nil {
		return fmt.Errorf("failed to load runtimes: %w", err)
	}

	if len(runtimes) == 0 {
		fmt.Println("No runtimes configured")
		return nil
	}

	tbl := table.NewWriter()
	tbl.SetStyle(table.StyleRounded)
	tbl.SetTitle("Available Runtimes")
	tbl.AppendHeader(table.Row{"ID", "Name", "Registry", "Description"})

	for _, runtime := range runtimes {
		tbl.AppendRow(table.Row{
			runtime.ID,
			runtime.Name,
			runtime.Registry,
			text.WrapSoft(runtime.Description, 40),
		})
	}

	fmt.Println(tbl.Render())
	return nil
}

func (c *standaloneCommand) addRuntime(pc *fisk.ParseContext) error {
	rm := standalone.NewRuntimeManager()

	// Set defaults if not provided
	name := c.runtimeName
	if name == "" {
		name = cases.Title(language.English).String(c.runtimeID) + " Runtime"
	}

	description := c.runtimeDescription
	if description == "" {
		description = fmt.Sprintf("Custom runtime: %s", c.runtimeID)
	}

	runtime := standalone.Runtime{
		ID:          c.runtimeID,
		Name:        name,
		Description: description,
		Registry:    c.runtimeRegistry,
		Author:      c.runtimeAuthor,
	}

	if err := rm.AddRuntime(runtime); err != nil {
		return fmt.Errorf("failed to add runtime: %w", err)
	}

	color.Green("✓ Added runtime '%s' with registry '%s'", c.runtimeID, c.runtimeRegistry)
	return nil
}

func (c *standaloneCommand) removeRuntime(pc *fisk.ParseContext) error {
	rm := standalone.NewRuntimeManager()

	if err := rm.RemoveRuntime(c.runtimeID); err != nil {
		return fmt.Errorf("failed to remove runtime: %w", err)
	}

	color.Green("✓ Removed runtime '%s'", c.runtimeID)
	return nil
}

func (c *standaloneCommand) showRuntime(pc *fisk.ParseContext) error {
	rm := standalone.NewRuntimeManager()

	runtime, err := rm.GetRuntime(c.runtimeID)
	if err != nil {
		return fmt.Errorf("failed to get runtime: %w", err)
	}

	tbl := table.NewWriter()
	tbl.SetStyle(table.StyleRounded)
	tbl.SetTitle(fmt.Sprintf("Runtime: %s", runtime.ID))

	tbl.AppendRow(table.Row{"ID", runtime.ID})
	tbl.AppendRow(table.Row{"Name", runtime.Name})
	tbl.AppendRow(table.Row{"Description", runtime.Description})
	tbl.AppendRow(table.Row{"Registry", runtime.Registry})
	if runtime.Author != "" {
		tbl.AppendRow(table.Row{"Author", runtime.Author})
	}

	fmt.Println(tbl.Render())
	return nil
}
