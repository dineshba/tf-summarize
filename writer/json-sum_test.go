package writer

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/dineshba/tf-summarize/terraformstate"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
)

func TestJSONSumWriter(t *testing.T) {
	changes := map[string]terraformstate.ResourceChanges{
		"add": {
			{Address: "aws_instance.a", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}}},
			{Address: "aws_instance.b", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}}},
		},
		"delete": {
			{Address: "aws_instance.c", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionDelete}}},
		},
		"update":   {},
		"recreate": {},
		"import":   {},
		"moved": {
			{
				Address:         "aws_instance.new",
				PreviousAddress: "aws_instance.old",
				Change:          &tfjson.Change{Actions: tfjson.Actions{}},
			},
		},
	}

	w := NewJSONSumWriter(changes)
	var buf bytes.Buffer
	err := w.Write(&buf)
	assert.NoError(t, err)

	var result map[string]map[string]int
	err = json.Unmarshal(buf.Bytes(), &result)
	assert.NoError(t, err)

	expected := map[string]map[string]int{
		"changes": {
			"add":      2,
			"delete":   1,
			"update":   0,
			"recreate": 0,
			"import":   0,
			"moved":    1,
		},
	}
	assert.Equal(t, expected, result)
}
