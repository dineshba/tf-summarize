package writer

import (
	"fmt"
	"io"
	"terraform-plan-summary/terraform_state"
	"terraform-plan-summary/tree"
)

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
		colorPrefix, suffix := tree.Value.ColorPrefixAndSuffixText()
		_, err = fmt.Fprintf(writer, "%s%s%s%s%s\n", prefix, colorPrefix, tree.Name, suffix, terraform_state.ColorReset)
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
