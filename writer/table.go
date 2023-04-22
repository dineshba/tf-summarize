package writer

import (
	"fmt"
	"github.com/dineshba/tf-summarize/terraform_state"
	"io"

	"github.com/olekukonko/tablewriter"
)

type TableWriter struct {
	mdEnabled      bool
	changes        map[string]terraform_state.ResourceChanges
	output_changes map[string][]string
	outputs        bool
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
	if len(t.output_changes) > 0 {
		tableString = make([][]string, 0, 4)
		for change, changedOutputs := range t.output_changes {
			for _, changedOutput := range changedOutputs {
				if t.mdEnabled {
					tableString = append(tableString, []string{change, fmt.Sprintf("`%s`", changedOutput)})
				} else {
					tableString = append(tableString, []string{change, changedOutput})
				}
			}
		}
		fmt.Println("\n")
		table = tablewriter.NewWriter(writer)
		table.SetHeader([]string{"Change", "Output"})
		table.SetAutoMergeCells(true)
		table.AppendBulk(tableString)

		if t.mdEnabled {
			table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
			table.SetCenterSeparator("|")
		} else {
			table.SetRowLine(true)
		}

		table.Render()
	}

	return nil
}

func NewTableWriter(changes map[string]terraform_state.ResourceChanges, output_changes map[string][]string, outputs, mdEnabled bool) Writer {
	// Disable the Output by setting output_changes to an empty map
	var output = make(map[string][]string)
	if outputs == true {
		output = output_changes
	}
	return TableWriter{
		changes:        changes,
		mdEnabled:      mdEnabled,
		output_changes: output,
		outputs:        outputs,
	}
}
