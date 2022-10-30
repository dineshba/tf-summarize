package writer

import (
	"encoding/json"
	"fmt"
	"io"
	"terraform-plan-summary/terraform_state"
	"terraform-plan-summary/tree"

	"github.com/nsf/jsondiff"
)

type JsonWriter struct {
	changes terraform_state.ResourceChanges
}

func (t JsonWriter) Write(writer io.Writer) error {
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
		_, suffix := t.Value.ColorPrefixAndSuffixText()
		var diff interface{}
		if t.IsUpdate() || t.IsRecreate() {
			opts := jsondiff.DefaultJSONOptions()
			opts.SkipMatches = true
			_, str := jsondiff.Compare(t.Value.Change.Before, t.Value.Change.After, &opts)
			diff = make(map[string]interface{})
			_ = json.Unmarshal([]byte(str), &diff)
		} else {
			if t.IsAddition() {
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

func NewJsonWriter(changes terraform_state.ResourceChanges) Writer {
	return JsonWriter{changes: changes}
}
