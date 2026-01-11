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
	moves         map[string]terraformstate.ResourceChanges
	outputChanges map[string][]string
}

var tableOrder = []string{"import", "add", "update", "recreate", "delete"}

func (t TableWriter) Write(writer io.Writer) error {
	tableString := make([][]string, 0, 4)

	for _, change := range tableOrder {
		changedResources := t.changes[change]
		resourceCount := len(changedResources)

		for _, changedResource := range changedResources {
			if t.mdEnabled {
				tableString = append(tableString, []string{fmt.Sprintf("%s (%d)", change, resourceCount), fmt.Sprintf("`%s`", changedResource.Address)})
			} else {
				tableString = append(tableString, []string{fmt.Sprintf("%s (%d)", change, resourceCount), changedResource.Address})
			}
		}
	}

	for move, movedResources := range t.moves {
		for _, movedResource := range movedResources {
			if t.mdEnabled {
				tableString = append(tableString, []string{move, fmt.Sprintf("`%s` to `%s`", movedResource.PreviousAddress, movedResource.Address)})
			} else {
				tableString = append(tableString, []string{move, fmt.Sprintf("%s to %s", movedResource.PreviousAddress, movedResource.Address)})
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
			outputCount := len(changedOutputs)
			for _, changedOutput := range changedOutputs {
				if t.mdEnabled {
					tableString = append(tableString, []string{fmt.Sprintf("%s (%d)", change, outputCount), fmt.Sprintf("`%s`", changedOutput)})
				} else {
					tableString = append(tableString, []string{fmt.Sprintf("%s (%d)", change, outputCount), changedOutput})
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

func NewTableWriter(changes map[string]terraformstate.ResourceChanges, moves map[string]terraformstate.ResourceChanges, outputChanges map[string][]string, mdEnabled bool) Writer {
	return TableWriter{
		changes:       changes,
		moves:         moves,
		mdEnabled:     mdEnabled,
		outputChanges: outputChanges,
	}
}
