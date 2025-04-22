package writer

import (
	"bytes"
	"testing"
	"testing/fstest"

	"github.com/dineshba/tf-summarize/terraformstate"
	. "github.com/hashicorp/terraform-json"
)

func TestHTMLWriterWithMockFileSystem(t *testing.T) {
	origFS := cfs
	cfs = fstest.MapFS{
		"templates/resourceChanges.html": &fstest.MapFile{
			Data: []byte(`<table>
  <tr>
    <th>CHANGE</th>
    <th>RESOURCE</th>
  </tr>{{ range $change, $resources := .ResourceChanges }}{{ $length := len $resources }}{{ if gt $length 0 }}
  <tr>
    <td>{{ $change }} ({{ len $resources }})</td>
    <td>
      <ul>{{ range $i, $r := $resources }}
        <li><code>{{ $r.Address }}</code></li>{{ end }}
      </ul>
    </td>
  </tr>{{ end }}{{ end }}
</table>
`),
		},
		"templates/outputChanges.html": &fstest.MapFile{
			Data: []byte(`<table>
  <tr>
    <th>CHANGE</th>
    <th>OUTPUT</th>
  </tr>{{ range $change, $outputs := .OutputChanges }}{{ $length := len $outputs }}{{ if gt $length 0 }}
  <tr>
    <td>{{ $change }}</td>
    <td>
      <ul>{{ range $i, $o := $outputs }}
        <li><code>{{ $o }}</code></li>{{ end }}
      </ul>
    </td>
  </tr>{{ end }}{{ end }}
</table>
`),
		},
	}
	t.Cleanup(func() {
		cfs = origFS
	})

	resourceChanges := map[string]terraformstate.ResourceChanges{
		"module.test": {
			{
				Address: "aws_instance.example",
				Name:    "example",
				Change: &Change{
					Before:  map[string]interface{}{"name": "old_instance"},
					After:   map[string]interface{}{"name": "new_instance"},
					Actions: Actions{ActionCreate},
				},
			},
		},
	}
	outputChanges := map[string][]string{
		"output_key": {"output_value"},
	}

	htmlWriter := NewHTMLWriter(resourceChanges, outputChanges)
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
    <td>module.test</td>
    <td>
      <ul>
        <li><code>aws_instance.example</code></li>
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
    <td>output_key</td>
    <td>
      <ul>
        <li><code>output_value</code></li>
      </ul>
    </td>
  </tr>
</table>
`
	if buf.String() != expectedOutput {
		t.Errorf("expected %q, got %q", expectedOutput, buf.String())
	}

}
