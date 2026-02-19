package writer

import (
	"io"

	"github.com/dineshba/tf-summarize/terraformstate"
	tfjson "github.com/hashicorp/terraform-json"
)

// Writer writes formatted Terraform plan output.
type Writer interface {
	Write(writer io.Writer) error
}

// CreateWriter returns a Writer based on the provided output format flags.
func CreateWriter(tree, separateTree, drawable, mdEnabled, json, html bool, jsonSum bool, plan tfjson.Plan) Writer {
	if tree {
		return NewTreeWriter(plan.ResourceChanges, drawable)
	}
	if separateTree {
		return NewSeparateTree(terraformstate.GetAllResourceChanges(plan), drawable)
	}
	if json {
		return NewJSONWriter(plan.ResourceChanges)
	}
	if html {
		return NewHTMLWriter(terraformstate.GetAllResourceChanges(plan), terraformstate.GetAllResourceMoves(plan), terraformstate.GetAllOutputChanges(plan))
	}
	if jsonSum {
		return NewJSONSumWriter(terraformstate.GetAllResourceChanges(plan))
	}

	return NewTableWriter(terraformstate.GetAllResourceChanges(plan), terraformstate.GetAllResourceMoves(plan), terraformstate.GetAllOutputChanges(plan), mdEnabled)
}
