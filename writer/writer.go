package writer

import (
	"io"

	"github.com/dineshba/tf-summarize/terraformstate"
)

type Writer interface {
	Write(writer io.Writer) error
}

func CreateWriter(tree, separateTree, drawable, mdEnabled, json bool, terraformState terraformstate.TerraformState) Writer {

	if tree {
		return NewTreeWriter(terraformState.ResourceChanges, drawable)
	}
	if separateTree {
		return NewSeparateTree(terraformState.AllResourceChanges(), drawable)
	}
	if json {
		return NewJSONWriter(terraformState.ResourceChanges)
	}

	return NewTableWriter(terraformState.AllResourceChanges(), terraformState.AllOutputChanges(), mdEnabled)
}
