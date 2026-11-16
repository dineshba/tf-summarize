package writer

import (
	"fmt"
	"io"

	"github.com/dineshba/tf-summarize/terraformstate"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
)

// TableWriter writes resource changes in a table format.
type TableWriter struct {
	mdEnabled     bool
	changes       map[string]terraformstate.ResourceChanges
	outputChanges map[string][]string
}

var tableOrder = []string{"import", "add", "update", "recreate", "delete", "moved"}

func rowCellConfig() tw.CellConfig {
	return tw.CellConfig{
		Merging: tw.CellMerging{Mode: tw.MergeBoth},
		Formatting: tw.CellFormatting{
			AutoWrap: tw.WrapNone,
		},
	}
}

func newChangeTable(writer io.Writer, mdEnabled bool) *tablewriter.Table {
	cfg := tablewriter.Config{
		Row:    rowCellConfig(),
		Header: rowCellConfig(),
	}

	if mdEnabled {
		return tablewriter.NewTable(writer,
			tablewriter.WithConfig(cfg),
			tablewriter.WithRenderer(renderer.NewMarkdown()),
		)
	}

	return tablewriter.NewTable(writer,
		tablewriter.WithConfig(cfg),
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{
			Symbols: tw.NewSymbols(tw.StyleASCII),
			Settings: tw.Settings{
				Separators: tw.Separators{
					BetweenRows: tw.On,
				},
			},
		})),
	)
}

func (t TableWriter) Write(writer io.Writer) error {
	tableString := make([][]string, 0, 4)

	for _, change := range tableOrder {
		changedResources := t.changes[change]
		resourceCount := len(changedResources)

		for _, changedResource := range changedResources {
			changeLabel := fmt.Sprintf("%s (%d)", change, resourceCount)
			if change == "moved" {
				if t.mdEnabled {
					tableString = append(tableString, []string{changeLabel, fmt.Sprintf("`%s` to `%s`", changedResource.PreviousAddress, changedResource.Address)})
				} else {
					tableString = append(tableString, []string{changeLabel, fmt.Sprintf("%s to %s", changedResource.PreviousAddress, changedResource.Address)})
				}
			} else {
				if t.mdEnabled {
					tableString = append(tableString, []string{changeLabel, fmt.Sprintf("`%s`", changedResource.Address)})
				} else {
					tableString = append(tableString, []string{changeLabel, changedResource.Address})
				}
			}
		}
	}

	table := newChangeTable(writer, t.mdEnabled)
	table.Header([]string{"Change", "Resource"})
	_ = table.Bulk(tableString)
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
		if t.mdEnabled {
			// Without a line break separating each table, a single malformed markdown table is printed.
			_, _ = fmt.Fprint(writer, "\n")
		}
		table = newChangeTable(writer, t.mdEnabled)
		table.Header([]string{"Change", "Output"})
		_ = table.Bulk(tableString)
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