package state

import (
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
)

func Summarize(tree, drawable, md bool, outputFileName string, stateValue stateV4) error {
	resources := make([]string, 0, len(stateValue.Resources))

	for _, resource := range stateValue.Resources {
		resources = append(resources, fmt.Sprintf("%s.%s.%s", resource.Module, resource.Type, resource.Name))
	}

	if tree {
		tree := CreateTree(resources)

		if drawable {
			fmt.Printf("%v\n", tree.DrawableTree())
			return nil
		}
		for _, t := range tree {
			err := printTree(os.Stdout, t, "")
			if err != nil {
				return fmt.Errorf("error writing data to %s: %s", "stdout", err.Error())
			}
		}
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Resources"})
	format := "%s"
	if md {
		format = "`%s`"
	}
	for _, resource := range resources {
		table.Append([]string{fmt.Sprintf(format, resource)})
	}

	if md {
		// Adding a println to break up the tables in md mode
		fmt.Println()
		table.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
		table.SetCenterSeparator("|")
	} else {
		table.SetRowLine(true)
	}

	table.Render()
	return nil
}

func printTree(writer io.Writer, tree *Tree, prefixSpace string) error {
	var err error
	prefixSymbol := fmt.Sprintf("%s|---", prefixSpace)
	// if tree.Value != nil {
	// 	colorPrefix, suffix := tree.Value.ColorPrefixAndSuffixText()
	// 	_, err = fmt.Fprintf(writer, "%s%s%s%s%s\n", prefixSymbol, colorPrefix, tree.Name, suffix, terraformstate.ColorReset)
	// } else {
	_, err = fmt.Fprintf(writer, "%s%s\n", prefixSymbol, tree.Name)
	// }
	if err != nil {
		return fmt.Errorf("error writing data to %s: %s", writer, err.Error())
	}

	for _, c := range tree.Children {
		separator := "|"
		err = printTree(writer, c, fmt.Sprintf("%s%s\t", prefixSpace, separator))
		if err != nil {
			return fmt.Errorf("error writing data to %s: %s", writer, err.Error())
		}
	}
	return nil
}
