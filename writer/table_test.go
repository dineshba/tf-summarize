package writer

import (
	"bytes"
	"testing"

	"github.com/dineshba/tf-summarize/terraformstate"
	"github.com/stretchr/testify/assert"

	. "github.com/hashicorp/terraform-json"
)

func TestTableWriter_Write_NoMarkdown(t *testing.T) {
	changes := createMockChanges()

	changes["update"] = terraformstate.ResourceChanges{
		{
			Address: "aws_instance.example3",
			Change:  &Change{Actions: Actions{ActionUpdate}},
		},
		{
			Address: "aws_instance.example4.tag[\"Custom Instance Tag\"]",
			Change:  &Change{Actions: Actions{ActionUpdate}},
		},
	}

	movedResources := map[string]terraformstate.ResourceChanges{
		"moved": {
			{
				Address:         "aws_instance.new",
				PreviousAddress: "aws_instance.old",
				Change:          &Change{Actions: Actions{}},
			},
		},
	}

	outputChanges := map[string][]string{
		"update": {
			"output.example",
			"output.long_resource_name.this[\"Custom/Resource Name\"]",
		},
	}

	tw := NewTableWriter(changes, movedResources, outputChanges, false)
	var output bytes.Buffer
	err := tw.Write(&output)
	assert.NoError(t, err)

	expectedOutput := `+------------+--------------------------------------------------+
|   CHANGE   |                     RESOURCE                     |
+------------+--------------------------------------------------+
| add (1)    | aws_instance.example1                            |
+------------+--------------------------------------------------+
| update (2) | aws_instance.example3                            |
+            +--------------------------------------------------+
|            | aws_instance.example4.tag["Custom Instance Tag"] |
+------------+--------------------------------------------------+
| delete (1) | aws_instance.example2                            |
+------------+--------------------------------------------------+
| moved      | aws_instance.old to aws_instance.new             |
+------------+--------------------------------------------------+
+------------+--------------------------------------------------------+
|   CHANGE   |                         OUTPUT                         |
+------------+--------------------------------------------------------+
| update (2) | output.example                                         |
+            +--------------------------------------------------------+
|            | output.long_resource_name.this["Custom/Resource Name"] |
+------------+--------------------------------------------------------+
`

	assert.Equal(t, expectedOutput, output.String())
}

func TestTableWriter_Write_WithMarkdown(t *testing.T) {
	changes := createMockChanges()

	movedResources := map[string]terraformstate.ResourceChanges{
		"moved": {
			{
				Address:         "aws_instance.new",
				PreviousAddress: "aws_instance.old",
				Change:          &Change{Actions: Actions{}},
			},
		},
	}

	outputChanges := map[string][]string{
		"update": {
			"output.example",
			"output.long_resource_name.this[\"Custom/Resource Name\"]",
		},
	}

	tw := NewTableWriter(changes, movedResources, outputChanges, true)
	var output bytes.Buffer
	err := tw.Write(&output)
	assert.NoError(t, err)

	expectedOutput := `|   CHANGE   |                 RESOURCE                 |
|------------|------------------------------------------|
| add (1)    | ` + "`aws_instance.example1`" + `                  |
| delete (1) | ` + "`aws_instance.example2`" + `                  |
| moved      | ` + "`aws_instance.old` to `aws_instance.new`" + ` |

|   CHANGE   |                          OUTPUT                          |
|------------|----------------------------------------------------------|
| update (2) | ` + "`output.example`" + `                                         |
|            | ` + "`output.long_resource_name.this[\"Custom/Resource Name\"]`" + ` |
`

	assert.Equal(t, expectedOutput, output.String())
}

func TestTableWriter_NoChanges(t *testing.T) {
	changes := map[string]terraformstate.ResourceChanges{}
	movedResources := map[string]terraformstate.ResourceChanges{}
	outputChanges := map[string][]string{}

	tw := NewTableWriter(changes, movedResources, outputChanges, false)
	var output bytes.Buffer
	err := tw.Write(&output)
	assert.NoError(t, err)

	expectedOutput := `+--------+----------+
| CHANGE | RESOURCE |
+--------+----------+
+--------+----------+
`
	assert.Equal(t, expectedOutput, output.String())
}
