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

	outputChanges := map[string][]string{
		"update": {
			"output.example",
			"output.long_resource_name.this[\"Custom/Resource Name\"]",
		},
	}

	tw := NewTableWriter(changes, outputChanges, false)
	var output bytes.Buffer
	err := tw.Write(&output)
	assert.NoError(t, err)

	expectedOutput := `+--------+--------------------------------------------------+
| CHANGE |                     RESOURCE                     |
+--------+--------------------------------------------------+
| add    | aws_instance.example1                            |
+--------+--------------------------------------------------+
| update | aws_instance.example3                            |
+        +--------------------------------------------------+
|        | aws_instance.example4.tag["Custom Instance Tag"] |
+--------+--------------------------------------------------+
| delete | aws_instance.example2                            |
+--------+--------------------------------------------------+
+--------+--------------------------------------------------------+
| CHANGE |                         OUTPUT                         |
+--------+--------------------------------------------------------+
| update | output.example                                         |
+        +--------------------------------------------------------+
|        | output.long_resource_name.this["Custom/Resource Name"] |
+--------+--------------------------------------------------------+
`

	assert.Equal(t, expectedOutput, output.String())
}

func TestTableWriter_Write_WithMarkdown(t *testing.T) {
	changes := createMockChanges()

	outputChanges := map[string][]string{
		"update": {
			"output.example",
			"output.long_resource_name.this[\"Custom/Resource Name\"]",
		},
	}

	tw := NewTableWriter(changes, outputChanges, true)
	var output bytes.Buffer
	err := tw.Write(&output)
	assert.NoError(t, err)

	expectedOutput := `| CHANGE |        RESOURCE         |
|--------|-------------------------|
| add    | ` + "`aws_instance.example1`" + ` |
| delete | ` + "`aws_instance.example2`" + ` |

| CHANGE |                          OUTPUT                          |
|--------|----------------------------------------------------------|
| update | ` + "`output.example`" + `                                         |
|        | ` + "`output.long_resource_name.this[\"Custom/Resource Name\"]`" + ` |
`

	assert.Equal(t, expectedOutput, output.String())
}

func TestTableWriter_NoChanges(t *testing.T) {
	changes := map[string]terraformstate.ResourceChanges{}
	outputChanges := map[string][]string{}

	tw := NewTableWriter(changes, outputChanges, false)
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
