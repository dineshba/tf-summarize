package writer

import (
	"io"
	"terraform-plan-summary/terraform_state"
)

type Writer interface {
	Write(writer io.Writer) error
}

func CreateWriter(tree, separateTree, drawable, mdEnabled, json bool, terraformState terraform_state.TerraformState) Writer {
	if tree {
		return NewTreeWriter(terraformState.ResourceChanges, drawable)
	}
	if separateTree {
		return NewSeparateTree(terraformState.AllChanges(), drawable)
	}
	if json {
		return NewJsonWriter(terraformState.ResourceChanges)
	}
	return NewTableWriter(terraformState.AllChanges(), mdEnabled)
}
