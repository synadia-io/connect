package cli

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/jedib0t/go-pretty/v6/text"
    "github.com/synadia-io/connect/client"
    "github.com/synadia-io/connect/convert"
    "github.com/synadia-io/connect/spec"
    "os"
    "time"

    "github.com/AlecAivazis/survey/v2"
    "github.com/choria-io/fisk"
    "github.com/fatih/color"
    "github.com/jedib0t/go-pretty/v6/table"
    "github.com/mitchellh/mapstructure"
    "github.com/synadia-io/connect/model"
    "gopkg.in/yaml.v3"
    jsonpatch "github.com/evanphx/json-patch/v5"
)

type connectorCommand struct {
    opts *Options

    file          string
    fileSetByUser bool

    pull                  bool
    pullUsername          string
    pullUsernameSetByUser bool
    pullPassword          string
    pullPasswordSetByUser bool

    replicas      int
    placementTags []string
    envVars       map[string]string

    startTimeout string

    id          string
    description string
    image       string

    targetId string
    reload   bool

    runtime     string
    interactive bool
}

func ConfigureConnectorCommand(parentCmd commandHost, opts *Options) {
    c := &connectorCommand{
        opts:    opts,
        envVars: make(map[string]string),
    }

    connectorCmd := parentCmd.Command("connector", "Manage connectors").Alias("c")

    connectorCmd.Command("list", "list connectors").Alias("ls").Action(c.listConnectors)

    getCmd := connectorCmd.Command("get", "get a connector").Alias("show").Action(c.getConnector)
    getCmd.Arg("connector", "the name of the connector").Required().StringVar(&c.id)

    saveCmd := connectorCmd.Command("edit", "Add or modify a connector").Alias("create").Action(c.saveConnector)
    saveCmd.Arg("id", "The id of the connector to create or modify").Required().StringVar(&c.id)
    saveCmd.Flag("file", "use the connector definition from the given file").Short('f').IsSetByUser(&c.fileSetByUser).Default("./ConnectFile").StringVar(&c.file)
    saveCmd.Flag("runtime", "The runtime id").Default("wombat").StringVar(&c.runtime)

    copyCmd := connectorCmd.Command("copy", "Copy a connector").Action(c.copyConnector)
    copyCmd.Arg("id", "The id of the connector to copy").Required().StringVar(&c.id)
    copyCmd.Arg("target-id", "The id of the new connector").Required().StringVar(&c.targetId)

    deleteCmd := connectorCmd.Command("delete", "delete a connector").Alias("rm").Action(c.removeConnector)
    deleteCmd.Arg("connector", "the name of the connector").Required().StringVar(&c.id)

    statusCmd := connectorCmd.Command("status", "get the status of a connector").Action(c.connectorStatus)
    statusCmd.Arg("id", "the id of the connector to get status for").Required().StringVar(&c.id)

    startCmd := connectorCmd.Command("start", "Deploy a connector").Action(c.startConnector)
    startCmd.Arg("id", "The id of the connector to deploy").Required().StringVar(&c.id)
    startCmd.Flag("pull", "Whether to pull the image").Default("false").UnNegatableBoolVar(&c.pull)
    startCmd.Flag("pull-username", "Username for the pull").IsSetByUser(&c.pullUsernameSetByUser).StringVar(&c.pullUsername)
    startCmd.Flag("pull-password", "Password for the pull").IsSetByUser(&c.pullPasswordSetByUser).StringVar(&c.pullPassword)
    startCmd.Flag("replicas", "Number of replicas to start").Default("1").IntVar(&c.replicas)
    startCmd.Flag("placement-tags", "Placement tags to use").StringsVar(&c.placementTags)
    startCmd.Flag("env", "Environment variables to set").Short('e').StringMapVar(&c.envVars)
    startCmd.Flag("start-timeout", "How long to wait for the component to be started").Default("5m").StringVar(&c.startTimeout)

    stopCmd := connectorCmd.Command("stop", "stop a connector").Action(c.stopConnector)
    stopCmd.Arg("id", "The id of the connector to stop").Required().StringVar(&c.id)

    reloadCmd := connectorCmd.Command("reload", "reload a connector").Alias("restart").Action(c.reloadConnector)
    reloadCmd.Arg("id", "The id of the connector to reload").Required().StringVar(&c.id)
}

func (c *connectorCommand) listConnectors(pc *fisk.ParseContext) error {
    appCtx, err := LoadOptions(c.opts)
    fisk.FatalIfError(err, "failed to load options")
    defer appCtx.Close()

    resp, err := appCtx.Client.ListConnectors(c.opts.Timeout)
    fisk.FatalIfError(err, "failed to list connectors")

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
        {Number: 4, Name: color.YellowString("~"), WidthMin: 3, WidthMax: 5, Align: text.AlignCenter, AlignHeader: text.AlignCenter},
        {Number: 5, Name: color.GreenString("\u25B6"), WidthMin: 3, WidthMax: 5, Align: text.AlignCenter, AlignHeader: text.AlignCenter},
        {Number: 6, Name: color.RedString("\u25FC"), WidthMin: 3, WidthMax: 5, Align: text.AlignCenter, AlignHeader: text.AlignCenter},
    })

    tbl.AppendHeader(table.Row{"Name", "Description", "Runtime", color.YellowString("~"), color.GreenString("\u25B6"), color.RedString("\u25FC")}, table.RowConfig{AutoMerge: true})

    for _, c := range resp {
        pending := ""
        if c.Instances.Pending > 0 {
            pending = color.YellowString("%d", c.Instances.Pending)
        }

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
            pending,
            running,
            stopped,
        })
    }

    fmt.Println(tbl.Render())

    return nil
}

func (c *connectorCommand) getConnector(pc *fisk.ParseContext) error {
    appCtx, err := LoadOptions(c.opts)
    fisk.FatalIfError(err, "failed to load options")
    defer appCtx.Close()

    resp, err := appCtx.Client.GetConnector(c.id, c.opts.Timeout)
    fisk.FatalIfError(err, "failed to list instances for %s", c.id)

    if resp == nil {
        fmt.Printf("Connector %s not found\n", c.id)
        return nil
    }

    fmt.Println(renderConnector(*resp))

    return nil
}

func (c *connectorCommand) removeConnector(pc *fisk.ParseContext) error {
    appCtx, err := LoadOptions(c.opts)
    fisk.FatalIfError(err, "failed to load options")
    defer appCtx.Close()

    err = appCtx.Client.DeleteConnector(c.id, c.opts.Timeout)
    if err != nil {
        return fmt.Errorf("failed to stop connector: %w", err)
    }

    fmt.Println(color.GreenString("Connector %s deleted", c.id))

    return nil
}

func (c *connectorCommand) connectorStatus(pc *fisk.ParseContext) error {
    appCtx, err := LoadOptions(c.opts)
    fisk.FatalIfError(err, "failed to load options")
    defer appCtx.Close()

    instances, err := appCtx.Client.ListConnectorInstances(c.id, c.opts.Timeout)
    fisk.FatalIfError(err, "failed to list instances for %s", c.id)

    if len(instances) == 0 {
        fmt.Printf("No instances found for component %s\n", c.id)
        return nil
    }

    tbl := table.NewWriter()
    tbl.SetStyle(table.StyleRounded)
    tbl.SetColumnConfigs([]table.ColumnConfig{
        {Number: 1, Name: "ID", WidthMin: 8, WidthMax: 8},
        {Number: 2, Name: "", WidthMin: 1, WidthMax: 1, Align: text.AlignCenter, AlignHeader: text.AlignCenter},
        {Number: 3, Name: "Uptime"},
        {Number: 4, Name: "Message"},
    })
    tbl.SetTitle(fmt.Sprintf("Connector: %s", c.id))
    tbl.AppendHeader(table.Row{"ID", "", "Uptime", "Message"})

    for _, i := range instances {
        uptime := ""
        if i.Status == model.InstanceStatusRunning && i.Uptime != nil {
            uptime = *i.Uptime
        }

        msg := ""
        if i.Message != nil {
            msg = *i.Message
        }

        status := ""
        switch i.Status {
        case model.InstanceStatusRunning:
            status = color.GreenString("\u25B6")
        case model.InstanceStatusStopped:
            status = color.RedString("\u25FC")
        case model.InstanceStatusPending:
            status = color.YellowString("~")
        }

        tbl.AppendRow(table.Row{
            i.Id[len(i.Id)-8:],
            status,
            uptime,
            msg,
        })
    }

    fmt.Println(tbl.Render())

    return nil
}

func (c *connectorCommand) startConnector(pc *fisk.ParseContext) error {
    appCtx, err := LoadOptions(c.opts)
    fisk.FatalIfError(err, "failed to load options")
    defer appCtx.Close()

    opts := &model.ConnectorStartOptions{
        Pull:          c.pull,
        Replicas:      c.replicas,
        PlacementTags: c.placementTags,
        EnvVars:       c.envVars,
        Timeout:       c.startTimeout,
    }

    if c.pullUsernameSetByUser || c.pullPasswordSetByUser {
        opts.PullAuth = &model.ConnectorStartOptionsPullAuth{
            Enabled: true,
        }

        if c.pullUsernameSetByUser {
            opts.PullAuth.Username = &c.pullUsername
        }

        if c.pullPasswordSetByUser {
            opts.PullAuth.Password = &c.pullPassword
        }
    }

    instances, err := appCtx.Client.StartConnector(c.id, opts, c.opts.Timeout)
    fisk.FatalIfError(err, "start")

    fmt.Printf("Connector %s instances started: \n", c.id)
    for _, i := range instances {
        fmt.Printf("  %s [%s]\n", i.Id, i.Status)
    }

    return nil
}

func (c *connectorCommand) reloadConnector(context *fisk.ParseContext) error {
    appCtx, err := LoadOptions(c.opts)
    fisk.FatalIfError(err, "failed to load options")
    defer appCtx.Close()

    instances, err := appCtx.Client.ListConnectorInstances(c.id, c.opts.Timeout)
    fisk.FatalIfError(err, "failed to get connector instances")

    if len(instances) >= 0 {
        stoppedInstances, err := appCtx.Client.StopConnector(c.id, c.opts.Timeout)
        fisk.FatalIfError(err, "failed to reload connector")

        fmt.Printf("Instances stopped:\n")
        for _, i := range stoppedInstances {
            fmt.Printf("  %s\n", i.Id)
        }
    }

    instances, err = appCtx.Client.StartConnector(c.id, &model.ConnectorStartOptions{}, c.opts.Timeout)
    fisk.FatalIfError(err, "failed to reload connector")

    fmt.Printf("Instances started: \n")
    for _, i := range instances {
        fmt.Printf("  %s\n", i.Id)
    }

    return nil
}

func (c *connectorCommand) stopConnector(pc *fisk.ParseContext) error {
    appCtx, err := LoadOptions(c.opts)
    fisk.FatalIfError(err, "failed to load options")
    defer appCtx.Close()

    instances, err := appCtx.Client.StopConnector(c.id, c.opts.Timeout)
    if err != nil {
        return fmt.Errorf("failed to stop connector: %w", err)
    }

    fmt.Printf("Connector %s stopped: \n", c.id)
    for _, i := range instances {
        fmt.Printf("  %s [%s]\n", i.Id, i.Status)
    }

    return nil
}

func (c *connectorCommand) saveConnector(pc *fisk.ParseContext) error {
    appCtx, err := LoadOptions(c.opts)
    fisk.FatalIfError(err, "failed to load options")
    defer appCtx.Close()

    // -- check if the connector exists
    conn, err := appCtx.Client.GetConnector(c.id, c.opts.Timeout)
    fisk.FatalIfError(err, "failed to get connector %s: %v", c.id, err)
    exists := conn != nil

    var sp spec.ConnectorSpec
    if exists {
        sp = spec.ConnectorSpec{
            Description: conn.Description,
            RuntimeId:   conn.RuntimeId,
            Steps:       convert.ConvertStepsToSpec(conn.Steps),
        }
    } else {
        if !c.fileSetByUser {
            ssp, err := c.selectConnectorTemplate(appCtx.Client)
            fisk.FatalIfError(err, "could not select template: %v", err)
            sp = *ssp
        }
    }

    var changed bool
    var result *spec.ConnectorSpec
    if c.fileSetByUser {
        result, changed, err = fromFile(&sp, c.file)
        fisk.FatalIfError(err, "could not load connector spec from file: %v", err)
    } else {
        result, changed, err = fromEditor(&sp)
        fisk.FatalIfError(err, "could not edit connector spec: %v", err)
    }

    if !changed {
        fmt.Println(color.YellowString("No changes made"))
        return nil
    }

    var connector *model.Connector
    if !exists {
        connector, err = appCtx.Client.CreateConnector(c.id, result.Description, result.RuntimeId, convert.ConvertStepsFromSpec(result.Steps), c.opts.Timeout)
        if err != nil {
            color.Red("Could not save connector: %s", err)
            os.Exit(1)
        }

        fmt.Printf("Created connector %s\n", color.GreenString(c.id))
    } else {
        b, err := createMergePatch(sp, result)
        if err != nil {
            color.Red("Could not marshall connector patch: %s", err)
            os.Exit(1)
        }

        connector, err = appCtx.Client.PatchConnector(c.id, string(b), c.opts.Timeout)
        if err != nil {
            color.Red("Could not save connector: %s", err)
            os.Exit(1)
        }

        fmt.Printf("Updated connector %s\n", color.GreenString(c.id))
    }

    fmt.Println(renderConnector(*connector))
    fmt.Println()

    // ask the user if we want to reload the connector
    choice := ""
    _ = survey.AskOne(&survey.Select{
        Message: "Do you want to reload the connector now?",
        Options: []string{"Yes", "No"},
    }, &choice, survey.WithValidator(survey.Required))

    if choice == "Yes" {
        _ = c.reloadConnector(pc)
    }

    return nil
}

func (c *connectorCommand) copyConnector(context *fisk.ParseContext) error {
    appCtx, err := LoadOptions(c.opts)
    fisk.FatalIfError(err, "failed to load options")
    defer appCtx.Close()

    // -- check if the connector exists
    conn, err := appCtx.Client.GetConnector(c.id, c.opts.Timeout)
    fisk.FatalIfError(err, "failed to get connector %s: %v", c.id, err)
    exists := conn != nil

    if !exists {
        fmt.Printf("Connector %s not found\n", c.id)
        return nil
    }

    _, err = appCtx.Client.CreateConnector(c.targetId, conn.Description, conn.RuntimeId, convert.ConvertStepsFromSpec(convert.ConvertStepsToSpec(conn.Steps)), c.opts.Timeout)
    fisk.FatalIfError(err, "failed to create connector %s: %v", c.targetId, err)

    fmt.Printf("Created connector %s\n", color.GreenString(c.targetId))
    return nil
}

func fromFile(existing *spec.ConnectorSpec, file string) (*spec.ConnectorSpec, bool, error) {
    // -- check if the file exists
    if _, err := os.Stat(file); os.IsNotExist(err) {
        return nil, false, fmt.Errorf("ConnectFile %q not found", file)
    }

    // -- open the file
    f, err := os.Open(file)
    fisk.FatalIfError(err, "failed to open ConnectFile %q: %v", file, err)
    defer f.Close()

    // -- read the file
    var sp spec.Spec
    if err := yaml.NewDecoder(f).Decode(&sp); err != nil {
        return nil, false, fmt.Errorf("failed to decode ConnectFile %q: %w", file, err)
    }

    if sp.Type != spec.SpecTypeConnector {
        return nil, false, fmt.Errorf("file %q is not a connector spec file", file)
    }

    var csp spec.ConnectorSpec
    if err := mapstructure.Decode(sp.Spec, &csp); err != nil {
        return nil, false, fmt.Errorf("failed to decode connector spec: %w", err)
    }

    changed := false
    if existing == nil {
        changed = true
    } else {
        b1, _ := yaml.Marshal(existing)
        b2, _ := yaml.Marshal(csp)
        changed = !bytes.Equal(b1, b2)
    }

    return &csp, changed, nil
}

func fromEditor(existing *spec.ConnectorSpec) (*spec.ConnectorSpec, bool, error) {
    configYml, err := yaml.Marshal(existing)
    if err != nil {
        return nil, false, fmt.Errorf("could not marshal connector configuration: %s", err)
    }

    tmpFile, err := os.CreateTemp("", "*.yaml")
    if err != nil {
        return nil, false, fmt.Errorf("could not create temporary file: %s", err)
    }
    defer os.Remove(tmpFile.Name())

    _, err = fmt.Fprint(tmpFile, string(configYml))
    if err != nil {
        return nil, false, fmt.Errorf("could not write to temporary file: %s", err)
    }
    tmpFile.Close()

    err = editFile(tmpFile.Name())
    if err != nil {
        return nil, false, fmt.Errorf("could not edit file: %s", err)
    }

    modifiedConfig, err := os.ReadFile(tmpFile.Name())
    if err != nil {
        return nil, false, fmt.Errorf("could not read modified file: %s", err)
    }

    var payload spec.ConnectorSpec
    // Use yaml.Unmarshal to support loading both YAML and JSON as input.
    if err := yaml.Unmarshal(modifiedConfig, &payload); err != nil {
        return nil, false, fmt.Errorf("could not parse connector configuration: %s", err)
    }

    return &payload, !bytes.Equal(configYml, modifiedConfig), nil
}

func (c *connectorCommand) selectConnectorTemplate(cl client.Client) (*spec.ConnectorSpec, error) {
    rt, err := cl.GetRuntime(c.runtime, 5*time.Second)
    if err != nil {
        return nil, fmt.Errorf("could not get runtime: %s", err)
    }
    if rt == nil {
        return nil, fmt.Errorf("runtime %s not found", c.runtime)
    }

    sp := spec.ConnectorSpec{
        Description: "A summary of what this connector does",
        RuntimeId:   rt.Id,
    }

    templates := []struct {
        name   string
        modify func()
    }{
        {name: "Inlet - Reading from a Source to NATS", modify: func() {
            sp.Steps = spec.StepsSpec{
                Source:   &spec.SourceStepSpec{},
                Producer: &spec.ProducerStepSpec{},
            }
        }},
        {name: "Outlet - Writing from NATS to a Sink", modify: func() {
            sp.Steps = spec.StepsSpec{
                Consumer: &spec.ConsumerStepSpec{},
                Sink:     &spec.SinkStepSpec{},
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

    return &sp, nil
}

func createMergePatch(original any, modified any) ([]byte, error) {
    originalB, err := json.Marshal(original)
    if err != nil {
        return nil, fmt.Errorf("could not marshal original connector: %w", err)
    }
    modifiedB, err := json.Marshal(modified)
    if err != nil {
        return nil, fmt.Errorf("could not marshal modified connector: %w", err)
    }

    return jsonpatch.CreateMergePatch(originalB, modifiedB)
}
