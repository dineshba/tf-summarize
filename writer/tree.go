package writer

import (
	"fmt"
	"io"

	"github.com/dineshba/tf-summarize/terraformstate"
	"github.com/dineshba/tf-summarize/tree"
)

// TreeWriter writes resource changes in a tree format.
type TreeWriter struct {
	changes  terraformstate.ResourceChanges
	drawable bool
}

func (t TreeWriter) Write(writer io.Writer) error {
	trees := tree.CreateTree(t.changes)

	if t.drawable {
		drawableTree := trees.DrawableTree()
		_, err := fmt.Fprint(writer, drawableTree.String())
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

// NewTreeWriter returns a new TreeWriter.
func NewTreeWriter(changes terraformstate.ResourceChanges, drawable bool) Writer {
	return TreeWriter{changes: changes, drawable: drawable}
}

func printTree(writer io.Writer, tree *tree.Tree, prefixSpace string) error {
	var err error
	prefixSymbol := fmt.Sprintf("%s|---", prefixSpace)
	if tree.Value != nil {
		colorPrefix, suffix := terraformstate.GetColorPrefixAndSuffixText(tree.Value)
		_, err = fmt.Fprintf(writer, "%s%s%s%s%s\n", prefixSymbol, colorPrefix, tree.Name, suffix, terraformstate.ColorReset)
	} else {
		_, err = fmt.Fprintf(writer, "%s%s\n", prefixSymbol, tree.Name)
	}
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
