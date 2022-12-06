package writer

import (
	"fmt"
	"io"

	"github.com/dineshba/tf-summarize/terraformstate"
)

type JsonSumWriter struct {
	changes map[string]terraformstate.ResourceChanges
}

func (t JsonSumWriter) Write(writer io.Writer) error {
	result := make(map[string]int, len(t.changes))
	for k, v := range t.changes {
		result[k] = len(v)
	}
	s, _ := Marshal(map[string]map[string]int{"changes": result})
	_, err := fmt.Fprint(writer, string(s))
	return err
}

func NewJsonSumWriter(changes map[string]terraformstate.ResourceChanges) Writer {
	return JsonSumWriter{changes: changes}
}
