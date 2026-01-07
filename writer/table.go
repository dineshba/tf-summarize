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
var tableColor = map[string]string{"import": "\033[34;1m", "add": "\033[32;1m", "update": "\033[33;1m", "recreate": "\033[36;1m", "delete": "\033[31;1m"}

func (t TableWriter) Write(writer io.Writer) error {
	tableString := make([][]string, 0, 4)

	for _, change := range tableOrder {
		changedResources := t.changes[change]

		for _, changedResource := range changedResources {
			if t.mdEnabled {
				tableString = append(tableString, []string{change, fmt.Sprintf("`%s`", changedResource.Address)})
			} else {
				tableString = append(tableString, []string{tableColor[change] + change + "\033[0m", changedResource.Address})
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

func NewTableWriter(changes map[string]terraformstate.ResourceChanges, moves map[string]terraformstate.ResourceChanges, outputChanges map[string][]string, mdEnabled bool) Writer {
	return TableWriter{
		changes:       changes,
		moves:         moves,
		mdEnabled:     mdEnabled,
		outputChanges: outputChanges,
	}
}
