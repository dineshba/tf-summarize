package writer

import (
	"bytes"
	"testing"

	"github.com/dineshba/tf-summarize/terraformstate"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
)

func TestCreateWriter_Table(t *testing.T) {
	plan := tfjson.Plan{
		ResourceChanges: []*tfjson.ResourceChange{
			{
				Address: "aws_instance.example",
				Type:    "aws_instance",
				Name:    "example",
				Change:  &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}},
			},
		},
	}
	w := CreateWriter(false, false, false, false, false, false, false, plan)
	assert.IsType(t, TableWriter{}, w)
}

func TestCreateWriter_TableMd(t *testing.T) {
	plan := tfjson.Plan{}
	w := CreateWriter(false, false, false, true, false, false, false, plan)
	assert.IsType(t, TableWriter{}, w)
}

func TestCreateWriter_Tree(t *testing.T) {
	plan := tfjson.Plan{}
	w := CreateWriter(true, false, false, false, false, false, false, plan)
	assert.IsType(t, TreeWriter{}, w)
}

func TestCreateWriter_SeparateTree(t *testing.T) {
	plan := tfjson.Plan{
		ResourceChanges: []*tfjson.ResourceChange{
			{
				Address: "aws_instance.example",
				Type:    "aws_instance",
				Name:    "example",
				Change:  &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}},
			},
		},
	}
	w := CreateWriter(false, true, false, false, false, false, false, plan)
	assert.IsType(t, SeparateTree{}, w)
}

func TestCreateWriter_JSON(t *testing.T) {
	plan := tfjson.Plan{}
	w := CreateWriter(false, false, false, false, true, false, false, plan)
	assert.IsType(t, JSONWriter{}, w)
}

func TestCreateWriter_HTML(t *testing.T) {
	plan := tfjson.Plan{
		ResourceChanges: []*tfjson.ResourceChange{
			{
				Address: "aws_instance.example",
				Type:    "aws_instance",
				Name:    "example",
				Change:  &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionDelete}},
			},
		},
	}
	w := CreateWriter(false, false, false, false, false, true, false, plan)
	assert.IsType(t, HTMLWriter{}, w)
}

func TestCreateWriter_JSONSum(t *testing.T) {
	plan := tfjson.Plan{
		ResourceChanges: []*tfjson.ResourceChange{
			{
				Address: "aws_instance.example",
				Type:    "aws_instance",
				Name:    "example",
				Change:  &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionUpdate}},
			},
		},
	}
	w := CreateWriter(false, false, false, false, false, false, true, plan)
	assert.IsType(t, JSONSumWriter{}, w)
}

func TestHTMLWriter_NoOutputChanges(t *testing.T) {
	changes := map[string]terraformstate.ResourceChanges{
		"add": {
			{
				Address: "aws_instance.example",
				Name:    "example",
				Change:  &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}},
			},
		},
	}
	outputChanges := map[string][]string{}

	w := NewHTMLWriter(changes, outputChanges)
	var buf bytes.Buffer
	err := w.Write(&buf)
	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "aws_instance.example")
	assert.NotContains(t, buf.String(), "OUTPUT")
}

func TestHTMLWriter_EmptyChanges(t *testing.T) {
	changes := map[string]terraformstate.ResourceChanges{}
	outputChanges := map[string][]string{}

	w := NewHTMLWriter(changes, outputChanges)
	var buf bytes.Buffer
	err := w.Write(&buf)
	assert.NoError(t, err)
}
