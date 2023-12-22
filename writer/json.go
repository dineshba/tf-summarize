package writer

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/dineshba/tf-summarize/terraformstate"
	"github.com/dineshba/tf-summarize/tree"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/nsf/jsondiff"
)

type JSONWriter struct {
	changes []*tfjson.ResourceChange
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

			before := t.Value.Change.Before.([]byte)
			after := t.Value.Change.After.([]byte)
			_, str := jsondiff.Compare(before, after, &opts)
			diff = make(map[string]interface{})
			_ = json.Unmarshal([]byte(str), &diff)
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

func NewJSONWriter(changes []*tfjson.ResourceChange) Writer {
	return JSONWriter{changes: changes}
}
