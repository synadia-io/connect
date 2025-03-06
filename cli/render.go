package cli

import (
    "fmt"
    "github.com/choria-io/fisk"
    "github.com/jedib0t/go-pretty/v6/text"
    "github.com/synadia-io/connect/model"

    "github.com/fatih/color"
    "github.com/jedib0t/go-pretty/v6/table"
    "gopkg.in/yaml.v3"
)

func renderBoolean(b *bool) string {
    if b == nil {
        return ""
    }
    if *b {
        return color.GreenString("yes")
    }

    return color.RedString("no")
}

func renderConnector(c model.Connector) string {
    tbl := table.NewWriter()
    tbl.SetStyle(table.StyleRounded)
    tbl.SetColumnConfigs([]table.ColumnConfig{
        {Number: 1, Name: "ID", WidthMin: 15, WidthMax: 15},
        {Number: 2, Name: "", WidthMin: 1, WidthMax: 60},
    })
    tbl.SetTitle(fmt.Sprintf("Connector: %s", c.ConnectorId))
    tbl.AppendRow(table.Row{"Description", text.WrapSoft(c.Description, 50)})

    tbl.AppendRow(table.Row{"Runtime", c.RuntimeId})

    b, err := yaml.Marshal(c.Steps)
    fisk.FatalIfError(err, "failed to render steps")

    return fmt.Sprintf("%s\n\n%s", tbl.Render(), string(b))
}

func renderRuntime(rt model.Runtime) string {
    tbl := table.NewWriter()
    tbl.SetStyle(table.StyleRounded)
    tbl.SetColumnConfigs([]table.ColumnConfig{
        {Number: 1, Name: "ID", WidthMin: 15, WidthMax: 15},
        {Number: 2, Name: "", WidthMin: 1, WidthMax: 60},
    })
    tbl.SetTitle(rt.Label)
    tbl.AppendRow(table.Row{"Id", rt.Id})

    if rt.Description != nil {
        tbl.AppendRow(table.Row{"Description", text.WrapSoft(*rt.Description, 50)})
    }
    tbl.AppendRow(table.Row{"Author Name", rt.Author.Name})

    if rt.Author.Email != nil {
        tbl.AppendRow(table.Row{"Author Email", *rt.Author.Email})
    }

    if rt.Author.Url != nil {
        tbl.AppendRow(table.Row{"Author URL", *rt.Author.Url})
    }
    tbl.AppendRow(table.Row{"Image", rt.Image})
    metrics := rt.Metrics != nil
    tbl.AppendRow(table.Row{"Metrics", renderBoolean(&metrics)})

    return tbl.Render()
}
