package writer

import (
	"fmt"
	"io"

	"github.com/dineshba/tf-summarize/terraformstate"
	"github.com/olekukonko/tablewriter"
)

// TableWriter writes resource changes in a table format.
type TableWriter struct {
	mdEnabled     bool
	changes       map[string]terraformstate.ResourceChanges
	outputChanges map[string][]string
}

var tableOrder = []string{"import", "add", "update", "recreate", "delete", "moved"}

var actionColors = map[string]tablewriter.Colors{
	"import":   {tablewriter.FgCyanColor},
	"add":      {tablewriter.FgGreenColor},
	"update":   {tablewriter.FgYellowColor},
	"recreate": {tablewriter.FgMagentaColor},
	"delete":   {tablewriter.FgRedColor},
	"moved":    {tablewriter.FgCyanColor},
}

func (t TableWriter) Write(writer io.Writer) error {
	type row struct {
		cells  []string
		action string
	}

	rows := make([]row, 0, 4)
	for _, change := range tableOrder {
		changedResources := t.changes[change]
		resourceCount := len(changedResources)

		for _, changedResource := range changedResources {
			changeLabel := fmt.Sprintf("%s (%d)", change, resourceCount)
			if change == "moved" {
				if t.mdEnabled {
					rows = append(rows, row{[]string{changeLabel, fmt.Sprintf("`%s` to `%s`", changedResource.PreviousAddress, changedResource.Address)}, change})
				} else {
					rows = append(rows, row{[]string{changeLabel, fmt.Sprintf("%s to %s", changedResource.PreviousAddress, changedResource.Address)}, change})
				}
			} else {
				if t.mdEnabled {
					rows = append(rows, row{[]string{changeLabel, fmt.Sprintf("`%s`", changedResource.Address)}, change})
				} else {
					rows = append(rows, row{[]string{changeLabel, changedResource.Address}, change})
				}
			}
		}
	}

	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Change", "Resource"})
	table.SetAutoMergeCells(true)
	table.SetAutoWrapText(false)

	if t.mdEnabled {
		for _, r := range rows {
			table.Append(r.cells)
		}
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
	} else {
		for _, r := range rows {
			colors := []tablewriter.Colors{actionColors[r.action], actionColors[r.action]}
			table.Rich(r.cells, colors)
		}
		table.SetRowLine(true)
		table.SetColumnSeparator("│")
		table.SetRowSeparator("─")
		table.SetCenterSeparator("┼")
	}

	table.Render()

	if hasOutputChanges(t.outputChanges) {
		outputRows := make([]row, 0, 4)
		for _, change := range tableOrder {
			changedOutputs := t.outputChanges[change]
			outputCount := len(changedOutputs)
			for _, changedOutput := range changedOutputs {
				if t.mdEnabled {
					outputRows = append(outputRows, row{[]string{fmt.Sprintf("%s (%d)", change, outputCount), fmt.Sprintf("`%s`", changedOutput)}, change})
				} else {
					outputRows = append(outputRows, row{[]string{fmt.Sprintf("%s (%d)", change, outputCount), changedOutput}, change})
				}
			}
		}

		table = tablewriter.NewWriter(writer)
		table.SetHeader([]string{"Change", "Output"})
		table.SetAutoMergeCells(true)
		table.SetAutoWrapText(false)

		if t.mdEnabled {
			_, _ = fmt.Fprint(writer, tablewriter.NEWLINE)
			for _, r := range outputRows {
				table.Append(r.cells)
			}
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("|")
		} else {
			for _, r := range outputRows {
				colors := []tablewriter.Colors{actionColors[r.action], actionColors[r.action]}
				table.Rich(r.cells, colors)
			}
			table.SetRowLine(true)
			table.SetColumnSeparator("│")
			table.SetRowSeparator("─")
			table.SetCenterSeparator("┼")
		}

		table.Render()
	}

	return nil
}

// NewTableWriter returns a new TableWriter.
func NewTableWriter(changes map[string]terraformstate.ResourceChanges, outputChanges map[string][]string, mdEnabled bool) Writer {
	return TableWriter{
		changes:       changes,
		mdEnabled:     mdEnabled,
		outputChanges: outputChanges,
	}
}
