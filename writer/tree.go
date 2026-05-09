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

	for i, t := range trees {
		isLast := i == len(trees)-1
		err := printTree(writer, t, "", isLast)
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

func printTree(writer io.Writer, tree *tree.Tree, prefix string, isLast bool) error {
	var err error
	connector := "├── "
	if isLast {
		connector = "└── "
	}

	if tree.Value != nil {
		colorPrefix, suffix := terraformstate.GetColorPrefixAndSuffixText(tree.Value)
		if suffix != "" {
			suffix = " " + suffix
		}
		_, err = fmt.Fprintf(writer, "%s%s%s%s%s\n", prefix+connector, colorPrefix, tree.Name, suffix, terraformstate.ColorReset)
	} else {
		_, err = fmt.Fprintf(writer, "%s%s\n", prefix+connector, tree.Name)
	}
	if err != nil {
		return fmt.Errorf("error writing data to %s: %s", writer, err.Error())
	}

	childPrefix := prefix + "│   "
	if isLast {
		childPrefix = prefix + "    "
	}

	for i, c := range tree.Children {
		childIsLast := i == len(tree.Children)-1
		err = printTree(writer, c, childPrefix, childIsLast)
		if err != nil {
			return fmt.Errorf("error writing data to %s: %s", writer, err.Error())
		}
	}
	return nil
}
