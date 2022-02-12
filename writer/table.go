package writer

import (
	"github.com/olekukonko/tablewriter"
	"io"
	"terraform-plan-summary/terraform-state"
)

type TableWriter struct {
	changes map[string]terraform_state.ResourceChanges
}

func (t TableWriter) Write(writer io.Writer) error {
	tableString := make([][]string, 0, 4)
	for change, changedResources := range t.changes {
		for _, changedResource := range changedResources {
			tableString = append(tableString, []string{change, changedResource.Address})
		}
	}

	table := tablewriter.NewWriter(writer)
	table.SetHeader([]string{"Change", "Name"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.AppendBulk(tableString)
	table.Render()

	return nil
}

func NewTableWriter(changes map[string]terraform_state.ResourceChanges) Writer {
	return TableWriter{changes: changes}
}
