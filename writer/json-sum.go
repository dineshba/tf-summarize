package writer

import (
	"fmt"
	"io"

	"github.com/dineshba/tf-summarize/terraformstate"
)

// JSONSumWriter writes a JSON summary of change counts by action type.
type JSONSumWriter struct {
	changes map[string]terraformstate.ResourceChanges
}

func (t JSONSumWriter) Write(writer io.Writer) error {
	result := make(map[string]int, len(t.changes))
	for k, v := range t.changes {
		result[k] = len(v)
	}
	s, _ := Marshal(map[string]map[string]int{"changes": result})
	_, err := fmt.Fprint(writer, string(s))
	return err
}

// NewJSONSumWriter returns a new JSONSumWriter.
func NewJSONSumWriter(changes map[string]terraformstate.ResourceChanges) Writer {
	return JSONSumWriter{changes: changes}
}
