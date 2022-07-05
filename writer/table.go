package writer

import (
	"fmt"
	"io"
	"terraform-plan-summary/terraform_state"

	"github.com/olekukonko/tablewriter"
)

type TableWriter struct {
	mdEnabled bool
	changes   map[string]terraform_state.ResourceChanges
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
	table.SetHeader([]string{"Change", "Name"})
	table.SetAutoMergeCells(true)
	table.AppendBulk(tableString)

	if t.mdEnabled {
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
	} else {
		table.SetRowLine(true)
	}

	table.Render()

	return nil
}

func NewTableWriter(changes map[string]terraform_state.ResourceChanges, mdEnabled bool) Writer {
	return TableWriter{
		changes:   changes,
		mdEnabled: mdEnabled,
	}
}
