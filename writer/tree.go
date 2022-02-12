package writer

import (
	"fmt"
	"io"
	"terraform-plan-summary/terraform-state"
	"terraform-plan-summary/tree"
)

const colorReset = "\033[0m"
const colorRed = "\033[31m"
const colorGreen = "\033[32m"
const colorYellow = "\033[33m"

type TreeWriter struct {
	changes  terraform_state.ResourceChanges
	drawable bool
}

func (t TreeWriter) Write(writer io.Writer) error {
	trees := tree.CreateTree(t.changes)

	if t.drawable {
		drawableTree := trees.DrawableTree()
		_, err := fmt.Fprintf(writer, drawableTree.String())
		return err
	}

	for _, t := range trees {
		err := printTree(writer, t, "")
		if err != nil {
			return fmt.Errorf("error writing data to %s: %s", writer, err.Error())
		}
	}
	return nil
}

func NewTreeWriter(changes terraform_state.ResourceChanges, drawable bool) Writer {
	return TreeWriter{changes: changes, drawable: drawable}
}

func printTree(writer io.Writer, tree *tree.Tree, prefix string) error {
	var err error
	if tree.Value != nil {
		actions := tree.Value.Change.Actions
		colorPrefix := ""
		suffix := ""
		if len(actions) == 1 {
			if actions[0] == "create" {
				colorPrefix = colorGreen
				suffix = "(+)"
			} else if actions[0] == "delete" {
				colorPrefix = colorRed
				suffix = "(-)"
			} else {
				colorPrefix = colorYellow
				suffix = "(~)"
			}
		}
		if len(actions) == 2 {
			colorPrefix = colorRed
			suffix = "(+/-)"
		}
		_, err = fmt.Fprintf(writer, "%s%s%s%s%s\n", prefix, colorPrefix, tree.Name, suffix, colorReset)

	} else {
		_, err = fmt.Fprintf(writer, "%s%s\n", prefix, tree.Name)
	}
	if err != nil {
		return fmt.Errorf("error writing data to %s: %s", writer, err.Error())
	}

	for _, c := range tree.Children {
		err = printTree(writer, c, fmt.Sprintf("%s----", prefix))
		if err != nil {
			return fmt.Errorf("error writing data to %s: %s", writer, err.Error())
		}
	}
	return nil
}
