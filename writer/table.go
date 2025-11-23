package writer

import (
	"fmt"
	"io"

	"github.com/dineshba/tf-summarize/terraformstate"
	"github.com/olekukonko/tablewriter"
)

type TableWriter struct {
	mdEnabled     bool
	changes       map[string]terraformstate.ResourceChanges
	outputChanges map[string][]string
}

var tableOrder = []string{"import", "move", "add", "update", "recreate", "delete"}

func (t TableWriter) Write(writer io.Writer) error {
	tableString := make([][]string, 0, 4)
	for _, change := range tableOrder {
		changedResources := t.changes[change]
		for _, changedResource := range changedResources {
			if t.mdEnabled {
				if change == "moved" {
					tableString = append(tableString, []string{change, fmt.Sprintf("`%s` to `%s`", changedResource.PreviousAddress, changedResource.Address)})
				} else {
					tableString = append(tableString, []string{change, fmt.Sprintf("`%s`", changedResource.Address)})
				}
			} else {
				tableString = append(tableString, []string{change, changedResource.Address})
			}
		}
	}

	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Change", "Resource"})
	table.SetAutoMergeCells(true)
	table.SetAutoWrapText(false)
	table.AppendBulk(tableString)

	if t.mdEnabled {
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
	} else {
		table.SetRowLine(true)
	}

	table.Render()

	// Disable the Output Summary if there are no outputs to display
	if hasOutputChanges(t.outputChanges) {
		tableString = make([][]string, 0, 4)
		for _, change := range tableOrder {
			changedOutputs := t.outputChanges[change]
			for _, changedOutput := range changedOutputs {
				if t.mdEnabled {
					tableString = append(tableString, []string{change, fmt.Sprintf("`%s`", changedOutput)})
				} else {
					tableString = append(tableString, []string{change, changedOutput})
				}
			}
		}
		table = tablewriter.NewWriter(writer)
		table.SetHeader([]string{"Change", "Output"})
		table.SetAutoMergeCells(true)
		table.SetAutoWrapText(false)
		table.AppendBulk(tableString)

		if t.mdEnabled {
			// Without a line break separating each table, a single malformed markdown table is printed.
			// Printing an empty newline ensures distinct, separate tables are rendered.
			fmt.Fprint(writer, tablewriter.NEWLINE)

			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("|")
		} else {
			table.SetRowLine(true)
		}

		table.Render()
	}

	return nil
}

func NewTableWriter(changes map[string]terraformstate.ResourceChanges, outputChanges map[string][]string, mdEnabled bool) Writer {
	return TableWriter{
		changes:       changes,
		mdEnabled:     mdEnabled,
		outputChanges: outputChanges,
	}
}
