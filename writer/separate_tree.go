package writer

import (
	"fmt"
	"io"
	"strings"

	"github.com/dineshba/tf-summarize/terraformstate"
)

const SEPARATOR = "###################"

type SeparateTree struct {
	changes  map[string](terraformstate.ResourceChanges)
	drawable bool
}

func (s SeparateTree) Write(writer io.Writer) error {
	var err error
	for k, v := range s.changes {
		if len(v) > 0 {
			_, err = fmt.Fprintf(writer, "%s %s %s\n", SEPARATOR, strings.ToUpper(k), SEPARATOR)
			if err != nil {
				return fmt.Errorf("error writing to %s: %s", writer, err)
			}
			treeWriter := NewTreeWriter(v, s.drawable)
			err = treeWriter.Write(writer)
			if err != nil {
				return fmt.Errorf("error writing to %s: %s", writer, err)
			}
			_, err = fmt.Fprintf(writer, "\n\n")
			if err != nil {
				return fmt.Errorf("error writing to %s: %s", writer, err)
			}
		}
	}
	return nil
}

func NewSeparateTree(changes map[string]terraformstate.ResourceChanges, drawable bool) Writer {
	return SeparateTree{changes: changes, drawable: drawable}
}
