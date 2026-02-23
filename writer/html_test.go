package writer

import (
	"bytes"
	"testing"

	"github.com/dineshba/tf-summarize/terraformstate"
	tfjson "github.com/hashicorp/terraform-json"
)

func TestHTMLWriter(t *testing.T) {
	resourceChanges := map[string]terraformstate.ResourceChanges{
		"update": {
			{
				Address: "aws_instance.example",
				Name:    "example",
				Change: &tfjson.Change{
					Before:  map[string]interface{}{"name": "old_instance"},
					After:   map[string]interface{}{"name": "new_instance"},
					Actions: tfjson.Actions{tfjson.ActionCreate},
				},
			},
		},
	}
	movedResources := map[string]terraformstate.ResourceChanges{
		"moved": {
			{
				Address:         "aws_instance.foo",
				PreviousAddress: "aws_instance.bar",
				Name:            "foo",
				Change: &tfjson.Change{
					Actions: tfjson.Actions{},
				},
			},
		},
	}
	outputChanges := map[string][]string{
		"output_key": {"output_value"},
	}

	htmlWriter := NewHTMLWriter(resourceChanges, movedResources, outputChanges)
	var buf bytes.Buffer

	err := htmlWriter.Write(&buf)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedOutput := `<table>
  <tr>
    <th>CHANGE</th>
    <th>RESOURCE</th>
  </tr>
  <tr>
    <td>update (1)</td>
    <td>
      <ul>
        <li><code>aws_instance.example</code></li>
      </ul>
    </td>
  </tr>
  <tr>
    <td>moved</td>
    <td>
      <ul>
        <li><code>aws_instance.bar</code> to <code>aws_instance.foo</code></li>
      </ul>
    </td>
  </tr>
</table>
<table>
  <tr>
    <th>CHANGE</th>
    <th>OUTPUT</th>
  </tr>
  <tr>
    <td>output_key (1)</td>
    <td>
      <ul>
        <li><code>output_value</code></li>
      </ul>
    </td>
  </tr>
</table>
`
	if buf.String() != expectedOutput {
		t.Errorf("expected %s, got %s", expectedOutput, buf.String())
	}

}
