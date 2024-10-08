package writer

import (
	"fmt"
	"io"

	"github.com/dineshba/tf-summarize/terraformstate"
	tfjson "github.com/hashicorp/terraform-json"
)

type Writer interface {
	Write(writer io.Writer) error
}

func CreateWriter(tree, separateTree, drawable, mdEnabled, json, html bool, plan tfjson.Plan) Writer {
	if tree {
		fmt.Printf("plan.ResourceChanges = %v", plan.ResourceChanges)

		// Alternatively, use %+v to print the entire struct details in the slice
		for i, change := range plan.ResourceChanges {
			fmt.Printf("ResourceChange %d: %+v\n", i+1, *change)
		}
		return NewTreeWriter(plan.ResourceChanges, drawable)
	}
	if separateTree {
		fmt.Println("terraformstate.GetAllResourceChanges(plan) = ", terraformstate.GetAllResourceChanges(plan))
		return NewSeparateTree(terraformstate.GetAllResourceChanges(plan), drawable)
	}
	if json {
		return NewJSONWriter(plan.ResourceChanges)
	}
	if html {
		return NewHTMLWriter(terraformstate.GetAllResourceChanges(plan), terraformstate.GetAllOutputChanges(plan))
	}

	return NewTableWriter(terraformstate.GetAllResourceChanges(plan), terraformstate.GetAllOutputChanges(plan), mdEnabled)
}
