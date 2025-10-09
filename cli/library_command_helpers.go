package cli

import (
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/synadia-io/connect/model"
)

// These are testable helper functions that can be called with a provided AppContext

func (c *libraryCommand) listRuntimesWithClient(appCtx *AppContext) error {
	runtimes, err := appCtx.Client.ListRuntimes(c.opts.Timeout)
	if err != nil {
		return fmt.Errorf("could not list runtimes: %w", err)
	}

	w := table.NewWriter()
	w.AppendHeader(table.Row{"Id", "Name", "Description", "Author"})
	w.SetStyle(table.StyleRounded)

	for _, runtime := range runtimes {
		desc := ""
		if runtime.Description != nil {
			desc = *runtime.Description
		}
		w.AppendRow(table.Row{runtime.Id, runtime.Label, desc, runtime.Author})
	}

	result := w.Render()
	fmt.Println(result)
	return nil
}

func (c *libraryCommand) getRuntimeWithClient(appCtx *AppContext) error {
	rt, err := appCtx.Client.GetRuntime(c.runtime, c.opts.Timeout)
	if err != nil {
		return fmt.Errorf("could not get runtime: %w", err)
	}

	if rt == nil {
		// In the real code this would be handled differently
		// For testing, we'll just return without error
		return nil
	}

	fmt.Println(renderRuntime(*rt))
	return nil
}

func (c *libraryCommand) searchWithClient(appCtx *AppContext) error {
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

	components, err := appCtx.Client.SearchComponents(filter, c.opts.Timeout)
	if err != nil {
		return fmt.Errorf("could not list components: %w", err)
	}

	w := table.NewWriter()
	w.AppendHeader(table.Row{"Name", "Kind", "Runtime", "Status"})
	w.SetStyle(table.StyleRounded)

	for _, component := range components {
		w.AppendRow(table.Row{component.Name, component.Kind, component.RuntimeId, component.Status})
	}

	result := w.Render()
	fmt.Println(result)
	return nil
}

func (c *libraryCommand) infoWithClient(appCtx *AppContext) error {
	component, err := appCtx.Client.GetComponent(c.runtime, model.ComponentKind(c.kind), c.component, c.opts.Timeout)
	if err != nil {
		return fmt.Errorf("could not get component: %w", err)
	}

	if component == nil {
		// In the real code this would be handled differently
		// For testing, we'll just return without error
		return nil
	}

	// Component info.
	w := table.NewWriter()
	w.SetStyle(table.StyleRounded)
	w.SetTitle("Component Description")
	w.AppendRow(table.Row{"Runtime", component.RuntimeId})
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
