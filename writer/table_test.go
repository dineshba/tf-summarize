package writer

import (
	"bytes"
	"testing"

	"github.com/dineshba/tf-summarize/terraformstate"
	"github.com/stretchr/testify/assert"
)

func TestTableWriter_Write_NoMarkdown(t *testing.T) {
	changes := map[string]terraformstate.ResourceChanges{
		"add": {
			{
				Address: "aws_instance.example1",
			},
		},
		"delete": {
			{
				Address: "aws_instance.example2",
			},
		},
	}

	outputChanges := map[string][]string{
		"update": {
			"output.example",
		},
	}

	tw := NewTableWriter(changes, outputChanges, false)
	var output bytes.Buffer
	err := tw.Write(&output)
	assert.NoError(t, err)

	expectedOutput := `+--------+-----------------------+
| CHANGE |       RESOURCE        |
+--------+-----------------------+
| add    | aws_instance.example1 |
+--------+-----------------------+
| delete | aws_instance.example2 |
+--------+-----------------------+
+--------+----------------+
| CHANGE |     OUTPUT     |
+--------+----------------+
| update | output.example |
+--------+----------------+
`

	assert.Equal(t, expectedOutput, output.String())
}

func TestTableWriter_Write_WithMarkdown(t *testing.T) {
	changes := map[string]terraformstate.ResourceChanges{
		"add": {
			{
				Address: "aws_instance.example1",
			},
		},
		"delete": {
			{
				Address: "aws_instance.example2",
			},
		},
	}

	outputChanges := map[string][]string{
		"update": {
			"output.example",
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

| CHANGE |      OUTPUT      |
|--------|------------------|
| update | ` + "`output.example`" + ` |
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
