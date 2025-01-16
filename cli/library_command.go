package cli

import (
	"fmt"
	"github.com/choria-io/fisk"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/synadia-io/connect/cli/render"
	"github.com/synadia-io/connect/library"
	"github.com/synadia-io/connect/model"

	"gopkg.in/yaml.v3"
	"os"
	"strings"
)

func init() {
	registerCommand("library", 0, configureLibraryCommand)
}

type libraryCommand struct {
	runtime   string
	version   string
	kind      string
	status    string
	component string
}

func configureLibraryCommand(parentCmd commandHost) {
	c := &libraryCommand{}

	componentCmd := parentCmd.Command("library", "Manage the library").Alias("l")

	kindOpts := []string{string(model.ComponentKindSource), string(model.ComponentKindSink), string(model.ComponentKindScanner)}

	componentCmd.Command("runtimes", "List the available runtimes").Alias("ls").Action(c.listRuntimes)

	runtimeCmd := componentCmd.Command("runtime", "Show information about a runtime").Action(c.getRuntime)
	runtimeCmd.Arg("id", "The id of the runtime to describe").Required().StringVar(&c.runtime)

	searchCmd := componentCmd.Command("search", "Search for components").Action(c.search)
	searchCmd.Flag("runtime", "The runtime id").Required().StringVar(&c.runtime)
	searchCmd.Flag("rev", "The version id or 'latest' to get the latest version").Default("latest").StringVar(&c.version)
	searchCmd.Flag("kind", "The kind of the component to search for").EnumVar(&c.kind, kindOpts...)
	searchCmd.Flag("status", "The status of the component to search for").EnumVar(&c.status, string(model.ComponentStatusActive), string(model.ComponentStatusPreview), string(model.ComponentStatusExperimental), string(model.ComponentStatusDeprecated))

	infoCmd := componentCmd.Command("show", "Show information about a component").Alias("get").Action(c.info)
	infoCmd.Flag("runtime", "The runtime id").Required().StringVar(&c.runtime)
	infoCmd.Flag("rev", "The version id or 'latest' to get the latest version").Default("latest").StringVar(&c.version)
	infoCmd.Arg("kind", "The kind of the component").EnumVar(&c.kind, kindOpts...)
	infoCmd.Arg("name", "The name of the component").StringVar(&c.component)
}

func (c *libraryCommand) listRuntimes(pc *fisk.ParseContext) error {
	w := table.NewWriter()
	w.AppendHeader(table.Row{"Id", "Latest", "Author"})
	w.SetStyle(table.StyleRounded)

	err := libraryClient().ListRuntimes(func(runtime *library.RuntimeInfo, hasMore bool) error {
		if runtime == nil {
			return nil
		}

		w.AppendRow(table.Row{runtime.Id, runtime.LatestVersionId, runtime.Author.Name})
		return nil
	})

	if err != nil {
		color.Red("Could not list runtimes: %s", err)
		os.Exit(1)
	}

	result := w.Render()
	fmt.Println(result)
	return nil
}

func (c *libraryCommand) getRuntime(pc *fisk.ParseContext) error {
	rt, err := libraryClient().GetRuntime(c.runtime)
	if err != nil {
		color.Red("Could not get runtime: %s", err)
		os.Exit(1)
	}

	rw := render.New("")
	rw.Println()
	rw.Println(rt.Description)
	rw.Println()

	rw.AddSectionTitle("Details")
	rw.AddRow("Id", rt.Id)

	rw.AddSectionTitle("Author")
	rw.AddRow("ConnectorId", rt.Author.Name)
	rw.AddRow("Email", rt.Author.Email)
	rw.AddRow("URL", rt.Author.Url)

	rw.AddSectionTitle("Versions")
	w := table.NewWriter()
	w.AppendHeader(table.Row{"Label", "Tags"})
	w.SetStyle(table.StyleRounded)
	versionFilter := library.VersionFilter{RuntimeId: rt.Id}

	err = libraryClient().ListVersions(versionFilter, func(version *library.VersionInfo, hasMore bool) error {
		if version == nil {
			return nil
		}

		tags := ""
		if version.VersionId == rt.LatestVersionId {
			tags = "latest"
		}

		w.AppendRow(table.Row{version.VersionId, tags})
		return nil
	})

	if err != nil {
		color.Red("Could not list versions: %s", err)
	}

	rw.Println(w.Render())
	return rw.Frender(os.Stdout)
}

func (c *libraryCommand) search(pc *fisk.ParseContext) error {
	w := table.NewWriter()
	w.AppendHeader(table.Row{"Kind", "ConnectorId", "Status"})
	w.SetStyle(table.StyleRounded)

	filter := library.ComponentFilter{
		RuntimeId: c.runtime,
		VersionId: c.version,
		Status:    model.ComponentStatus(c.status),
		Kind:      model.ComponentKind(c.kind),
	}

	err := libraryClient().ListComponents(filter, func(component *library.ComponentInfo, hasMore bool) error {
		if component == nil {
			return nil
		}

		w.AppendRow(table.Row{component.Kind, component.Name, component.Status})
		return nil
	})

	if err != nil {
		color.Red("Could not list components: %s", err)
		os.Exit(1)
	}

	result := w.Render()
	fmt.Println(result)
	return nil
}

func (c *libraryCommand) info(pc *fisk.ParseContext) error {
	component, err := libraryClient().GetComponent(c.runtime, c.version, model.ComponentKind(c.kind), c.component)
	if err != nil {
		color.Red("Could not get component: %s", err)
		os.Exit(1)
	}

	// Component info.
	w := table.NewWriter()
	w.SetStyle(table.StyleRounded)
	w.SetTitle("Component Description")
	w.AppendRow(table.Row{"Runtime", component.RuntimeId})
	w.AppendRow(table.Row{"Version", component.VersionId})
	w.AppendRow(table.Row{"Name", component.Name})
	w.AppendRow(table.Row{"Kind", component.Kind})
	w.AppendRow(table.Row{"Status", component.Status})
	w.AppendRow(table.Row{"Description", text.WrapSoft(component.Description, 75)})
	result := w.Render()
	fmt.Println(result)

	for _, field := range component.Fields {
		printField(field, "")
	}

	return nil
}

func printField(field model.Field, prefix string) {
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
		fmt.Printf(", defaults to %s\n", field.Default)
	} else {
		fmt.Print("\n")
	}

	if field.Description != "" {
		wrappedDescription := text.WrapSoft(field.Description, 85)
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
			if constraint.Regex != "" {
				fmt.Printf("\t\t%s: %s\n", "regex", constraint.Regex)
			}
			if constraint.Range != nil {
				if constraint.Range.LesserThan != 0 {
					fmt.Printf("\t\t%s: %v\n", "lesser than", constraint.Range.LesserThan)
				}
				if constraint.Range.LesserThanEqual != 0 {
					fmt.Printf("\t\t%s: %v\n", "lesser than equal", constraint.Range.LesserThanEqual)
				}
				if constraint.Range.GreaterThan != 0 {
					fmt.Printf("\t\t%s: %v\n", "greater than", constraint.Range.GreaterThan)
				}
				if constraint.Range.GreaterThanEqual != 0 {
					fmt.Printf("\t\t%s: %v\n", "greater than equal", constraint.Range.GreaterThanEqual)
				}
			}
			if len(constraint.Enum) > 0 {
				fmt.Printf("\t\t%s: %s\n", "enum", constraint.Enum)
			}
			if constraint.Preset != "" {
				fmt.Printf("\t\t%s: %s\n", "preset", constraint.Preset)
			}
		}
	}

	// Separator
	fmt.Println()

	for _, innerField := range field.Fields {
		genPrefix := fmt.Sprintf("%s%s.", prefix, field.Name)
		printField(innerField, genPrefix)
	}
}
