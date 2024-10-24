package writer

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/dineshba/tf-summarize/terraformstate"
	"github.com/dineshba/tf-summarize/tree"
	"github.com/nsf/jsondiff"
)

type JSONWriter struct {
	changes terraformstate.ResourceChanges
}

func (t JSONWriter) Write(writer io.Writer) error {
	trees := tree.CreateTree(t.changes)

	resultMap := make(map[string]interface{})
	for _, t := range trees {
		resultMap[t.Name] = treeValue(*t)
	}
	s, _ := Marshal(resultMap)
	_, err := fmt.Fprint(writer, string(s))
	return err
}

func treeValue(t tree.Tree) interface{} {
	resultMap := make(map[string]interface{})

	if t.Value != nil {
		_, suffix := terraformstate.GetColorPrefixAndSuffixText(t.Value)
		var diff interface{}
		if t.IsUpdate() || t.IsRecreate() {
			opts := jsondiff.DefaultJSONOptions()
			opts.SkipMatches = true

			before, _ := json.Marshal(t.Value.Change.Before)
			after, _ := json.Marshal(t.Value.Change.After)
			_, str := jsondiff.Compare(before, after, &opts)
			diff = make(map[string]interface{})
			err := json.Unmarshal([]byte(str), &diff)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: unmarshalling diff error: %s\n", err)
				diff = fmt.Sprintf("raw diff: %s", str)
			}
		} else {
			if t.IsAddition() || t.IsImport() {
				diff = t.Value.Change.After
			}
			if t.IsRemoval() {
				diff = t.Value.Change.Before
			}
		}

		resultMap[suffix] = diff
		return resultMap
	}
	for _, child := range t.Children {
		resultMap[child.Name] = treeValue(*child)
	}
	return resultMap
}

func NewJSONWriter(changes terraformstate.ResourceChanges) Writer {
	return JSONWriter{changes: changes}
}
