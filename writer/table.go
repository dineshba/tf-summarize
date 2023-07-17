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

func (t TableWriter) Write(writer io.Writer) error {
	tableString := make([][]string, 0, 4)
	for change, changedResources := range t.changes {
		for _, changedResource := range changedResources {
			if t.mdEnabled {
				tableString = append(tableString, []string{change, fmt.Sprintf("`%s`", changedResource.Address)})
			} else {
				tableString = append(tableString, []string{change, changedResource.Address})
			}
		}
	}

	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Change", "Resource"})
	table.SetAutoMergeCells(true)
	table.AppendBulk(tableString)

	if t.mdEnabled {
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
	} else {
		table.SetRowLine(true)
	}

	table.Render()

	// Disable the Output Summary if there are no outputs to display
	if len(t.outputChanges["add"]) > 0 || len(t.outputChanges["delete"]) > 0 || len(t.outputChanges["update"]) > 0 {
		tableString = make([][]string, 0, 4)
		for change, changedOutputs := range t.outputChanges {
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
		table.AppendBulk(tableString)

		if t.mdEnabled {
			// Adding a println to break up the tables in md mode
			fmt.Println()
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
