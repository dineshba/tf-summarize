package writer

import (
	"io"

	"github.com/dineshba/tf-summarize/terraform_state"
)

type Writer interface {
	Write(writer io.Writer) error
}

func CreateWriter(tree, separateTree, drawable, mdEnabled, json, jsonSum bool, terraformState terraform_state.TerraformState) Writer {
	if tree {
		return NewTreeWriter(terraformState.ResourceChanges, drawable)
	}
	if separateTree {
		return NewSeparateTree(terraformState.AllChanges(), drawable)
	}
	if json {
		return NewJsonWriter(terraformState.ResourceChanges)
	}
	if jsonSum {
		return NewJsonSumWriter(terraformState.AllChanges())
	}
	return NewTableWriter(terraformState.AllChanges(), mdEnabled)
}
