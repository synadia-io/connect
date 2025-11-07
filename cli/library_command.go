package cli

import (
	"fmt"

	"github.com/choria-io/fisk"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/synadia-io/connect/model"

	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

type libraryCommand struct {
	opts *Options

	runtime        string
	kind           string
	status         string
	component      string
	runtimeVersion string
}

func ConfigureLibraryCommand(parentCmd commandHost, opts *Options) {
	c := &libraryCommand{
		opts: opts,
	}

	componentCmd := parentCmd.Command("library", "Manage the library").Alias("l")

	kindOpts := []string{string(model.ComponentKindSource), string(model.ComponentKindSink), string(model.ComponentKindScanner)}
	statusOpts := []string{string(model.ComponentStatusStable), string(model.ComponentStatusPreview), string(model.ComponentStatusExperimental), string(model.ComponentStatusDeprecated)}

	componentCmd.Command("runtimes", "List the available runtimes").Action(c.listRuntimes)

	runtimeCmd := componentCmd.Command("runtime", "Show information about a runtime").Action(c.getRuntime)
	runtimeCmd.Arg("id", "The id of the runtime to describe").Required().StringVar(&c.runtime)
	runtimeCmd.Flag("runtime-version", "The runtime version").Required().StringVar(&c.runtimeVersion)

	versionsCmd := componentCmd.Command("versions", "List versions of a runtime").Action(c.listRuntimeVersions)
	versionsCmd.Arg("id", "The runtime id").Required().StringVar(&c.runtime)

	searchCmd := componentCmd.Command("list", "List components").Alias("ls").Action(c.search)
	searchCmd.Flag("runtime", "The runtime id").StringVar(&c.runtime)
	searchCmd.Flag("runtime-version", "The runtime version").StringVar(&c.runtimeVersion)
	searchCmd.Flag("kind", "The kind of components").EnumVar(&c.kind, kindOpts...)
	searchCmd.Flag("status", "The status of the components").EnumVar(&c.status, statusOpts...)

	infoCmd := componentCmd.Command("get", "Get component information").Alias("show").Action(c.info)
	infoCmd.Arg("runtime", "The runtime id").StringVar(&c.runtime)
	infoCmd.Arg("kind", "The kind of component").EnumVar(&c.kind, kindOpts...)
	infoCmd.Arg("name", "The name of the component").StringVar(&c.component)
	infoCmd.Flag("runtime-version", "The runtime version").StringVar(&c.runtimeVersion)
}

func (c *libraryCommand) listRuntimes(pc *fisk.ParseContext) error {
	appCtx, err := LoadOptions(c.opts)
	fisk.FatalIfError(err, "failed to load options")
	defer appCtx.Close()

	w := table.NewWriter()
	w.AppendHeader(table.Row{"Id", "Version", "Name", "Description", "Author"})
	w.SetStyle(table.StyleRounded)

	runtimes, err := appCtx.Client.ListRuntimes(c.opts.Timeout)
	if err != nil {
		color.Red("Could not list runtimes: %s", err)
		os.Exit(1)
	}

	for _, runtime := range runtimes {
		desc := ""
		if runtime.Description != nil {
			desc = *runtime.Description
		}
		version := runtime.Version
		w.AppendRow(table.Row{runtime.Id, version, runtime.Label, desc, runtime.Author.Name})
	}

	result := w.Render()
	fmt.Println(result)
	return nil
}

func (c *libraryCommand) getRuntime(pc *fisk.ParseContext) error {
	appCtx, err := LoadOptions(c.opts)
	fisk.FatalIfError(err, "failed to load options")
	defer appCtx.Close()

	rt, err := appCtx.Client.GetRuntime(c.runtime, &c.runtimeVersion, c.opts.Timeout)
	if err != nil {
		color.Red("Could not get runtime: %s", err)
		os.Exit(1)
	}

	if rt == nil {
		color.Red("Runtime not found")
		os.Exit(1)
	}

	fmt.Println(renderRuntime(*rt))

	return nil
}

func (c *libraryCommand) listRuntimeVersions(pc *fisk.ParseContext) error {
	appCtx, err := LoadOptions(c.opts)
	fisk.FatalIfError(err, "failed to load options")
	defer appCtx.Close()

	runtimes, err := appCtx.Client.ListRuntimes(c.opts.Timeout)
	if err != nil {
		color.Red("Could not list runtimes: %s", err)
		os.Exit(1)
	}

	var versions []string
	for _, runtime := range runtimes {
		if runtime.Id == c.runtime {
			versions = append(versions, runtime.Version)
		}
	}

	if len(versions) == 0 {
		color.Red("No versions found for runtime '%s'", c.runtime)
		os.Exit(1)
	}

	fmt.Printf("Available versions for runtime '%s':\n", c.runtime)
	w := table.NewWriter()
	w.AppendHeader(table.Row{"Version"})
	w.SetStyle(table.StyleRounded)

	for _, version := range versions {
		w.AppendRow(table.Row{version})
	}

	result := w.Render()
	fmt.Println(result)
	return nil
}

func (c *libraryCommand) search(pc *fisk.ParseContext) error {
	appCtx, err := LoadOptions(c.opts)
	fisk.FatalIfError(err, "failed to load options")
	defer appCtx.Close()

	w := table.NewWriter()
	w.AppendHeader(table.Row{"Name", "Kind", "Runtime", "Version", "Status"})
	w.SetStyle(table.StyleRounded)

	filter := &model.ComponentSearchFilter{}

	if c.runtime != "" {
		filter.RuntimeId = &c.runtime
	}

	if c.status != "" {
		st := model.ComponentStatus(c.status)
		filter.Status = &st
	}

	if c.kind != "" {
		k := model.ComponentKind(c.kind)
		filter.Kind = &k
	}

	if c.runtimeVersion != "" {
		filter.RuntimeVersion = &c.runtimeVersion
	}

	components, err := appCtx.Client.SearchComponents(filter, c.opts.Timeout)
	if err != nil {
		color.Red("Could not list components: %s", err)
		os.Exit(1)
	}

	for _, component := range components {
		w.AppendRow(table.Row{component.Name, component.Kind, component.RuntimeId, component.RuntimeVersion, component.Status})
	}

	result := w.Render()
	fmt.Println(result)
	return nil
}

func (c *libraryCommand) info(pc *fisk.ParseContext) error {
	appCtx, err := LoadOptions(c.opts)
	fisk.FatalIfError(err, "failed to load options")
	defer appCtx.Close()

	var runtimeVersion *string
	if c.runtimeVersion != "" {
		runtimeVersion = &c.runtimeVersion
	}

	component, err := appCtx.Client.GetComponent(c.runtime, model.ComponentKind(c.kind), c.component, runtimeVersion, c.opts.Timeout)
	if err != nil {
		color.Red("Could not get component: %s", err)
		os.Exit(1)
	}

	// Component info.
	w := table.NewWriter()
	w.SetStyle(table.StyleRounded)
	w.SetTitle("Component Description")
	w.AppendRow(table.Row{"Runtime", component.RuntimeId})
	w.AppendRow(table.Row{"Runtime Version", component.RuntimeVersion})
	w.AppendRow(table.Row{"Name", component.Name})
	w.AppendRow(table.Row{"Kind", component.Kind})
	w.AppendRow(table.Row{"Status", component.Status})

	if component.Description != nil {
		w.AppendRow(table.Row{"Description", text.WrapSoft(*component.Description, 75)})
	}

	result := w.Render()
	fmt.Println(result)

	for _, field := range component.Fields {
		printField(field, "")
	}

	return nil
}

func printField(field model.ComponentField, prefix string) {
	boldUnderline := color.New(color.Bold, color.Underline)
	bold := color.New(color.Bold)

	indent := func(s string, indent int) string {
		tab := strings.Repeat("\t", indent)
		indented := tab + strings.ReplaceAll(s, "\n", "\n"+tab)
		return strings.TrimSuffix(indented, tab)
	}

	name := fmt.Sprintf("%s%s", prefix, field.Name)
	_, _ = boldUnderline.Printf("%s\n", name)

	fmt.Printf("\t%s field", field.Type)
	if field.Default != nil {
		fmt.Printf(", defaults to %v\n", field.Default)
	} else {
		fmt.Print("\n")
	}

	if field.Description != nil {
		wrappedDescription := text.WrapSoft(*field.Description, 85)
		fmt.Printf("\n%s\n", indent(wrappedDescription, 1))
	}

	if len(field.Examples) > 0 {
		_, _ = bold.Printf("\n\t%s\n", "Examples")
		for _, example := range field.Examples {
			m := map[string]any{field.Name: example}
			yml, _ := yaml.Marshal(m)
			fmt.Printf("%s", indent(string(yml), 2))
		}
	}

	if len(field.Constraints) > 0 {
		_, _ = bold.Printf("\n\t%s\n", "Constraints")

		for _, constraint := range field.Constraints {
			if constraint.Regex != nil {
				fmt.Printf("\t\t%s: %s\n", "regex", *constraint.Regex)
			}
			if constraint.Range != nil {
				if constraint.Range.Lt != nil {
					fmt.Printf("\t\t%s: %v\n", "lesser than", *constraint.Range.Lt)
				}
				if constraint.Range.Lte != nil {
					fmt.Printf("\t\t%s: %v\n", "lesser than equal", constraint.Range.Lte)
				}
				if constraint.Range.Gt != nil {
					fmt.Printf("\t\t%s: %v\n", "greater than", constraint.Range.Gt)
				}
				if constraint.Range.Gte != nil {
					fmt.Printf("\t\t%s: %v\n", "greater than equal", constraint.Range.Gte)
				}
			}
			if len(constraint.Enum) > 0 {
				fmt.Printf("\t\t%s: %s\n", "enum", constraint.Enum)
			}
			if constraint.Preset != nil {
				fmt.Printf("\t\t%s: %s\n", "preset", *constraint.Preset)
			}
		}
	}

	// Separator
	fmt.Println()

	for _, innerField := range field.Fields {
		if innerField == nil {
			continue
		}

		genPrefix := fmt.Sprintf("%s%s.", prefix, field.Name)
		printField(*innerField, genPrefix)
	}
}
