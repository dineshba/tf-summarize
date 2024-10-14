package writer

import (
	"bytes"
	"errors"
	"testing"

	"github.com/dineshba/tf-summarize/terraformstate"
	. "github.com/hashicorp/terraform-json"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
)

func TestTreeWriter_Write_DrawableTrue(t *testing.T) {

	changes := terraformstate.ResourceChanges{
		&tfjson.ResourceChange{Address: "module.test.azapi_resource.logical_network", Change: &Change{Actions: Actions{ActionNoop}}},
	}

	tw := NewTreeWriter(changes, true)

	var buf bytes.Buffer
	err := tw.Write(&buf)

	assert.NoError(t, err)

	expectedOutput := `       ╭─╮       
       │.│       
       ╰┬╯       
        │        
    ╭───┴──╮     
    │module│     
    ╰───┬──╯     
        │        
     ╭──┴─╮      
     │test│      
     ╰──┬─╯      
        │        
╭───────┴──────╮ 
│azapi_resource│ 
╰───────┬──────╯ 
        │        
╭───────┴───────╮
│logical_network│
╰───────────────╯
`

	assert.Equal(t, expectedOutput, removeANSI(buf.String()))
}

func TestTreeWriter_Write_NonDrawable(t *testing.T) {

	changes := terraformstate.ResourceChanges{
		&tfjson.ResourceChange{Address: "module.test.azapi_resource.logical_network", Change: &Change{Actions: Actions{ActionNoop}}},
	}

	tw := NewTreeWriter(changes, false)

	var buf bytes.Buffer
	err := tw.Write(&buf)

	assert.NoError(t, err)
	assert.Contains(t, buf.String(), "|---module")
	assert.Contains(t, buf.String(), "|---test")
	assert.Contains(t, buf.String(), "|---azapi_resource")
	assert.Contains(t, buf.String(), "|---logical_network")

	expectedOutput := `|---module
|	|---test
|	|	|---azapi_resource
|	|	|	|---logical_network
`

	assert.Equal(t, expectedOutput, removeANSI(buf.String()))
}

func TestTreeWriter_Write_NonDrawable_PrintTreeError(t *testing.T) {

	changes := terraformstate.ResourceChanges{
		&tfjson.ResourceChange{Address: "module.test.azapi_resource.logical_network", Change: &Change{Actions: Actions{ActionNoop}}},
	}
	tw := NewTreeWriter(changes, false)
	faultyWriter := &errorWriter{}
	err := tw.Write(faultyWriter)
	assert.Error(t, err)
}

func TestTreeWriter_Write_PrintTreeError(t *testing.T) {
	changes := terraformstate.ResourceChanges{
		&tfjson.ResourceChange{Address: "module.test.azapi_resource.logical_network", Change: &Change{Actions: Actions{ActionNoop}}},
	}
	tw := NewTreeWriter(changes, true)
	faultyWriter := &errorWriter{}
	err := tw.Write(faultyWriter)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "simulated write error")
}

func TestTreeWriter_Write_EmptyChanges(t *testing.T) {
	changes := terraformstate.ResourceChanges{} // Empty changes
	tw := NewTreeWriter(changes, false)
	var buf bytes.Buffer
	err := tw.Write(&buf)

	// Verify output and no errors (it should handle empty cases gracefully)
	assert.NoError(t, err)
	assert.Equal(t, "", buf.String())
}

// Custom faulty writer to simulate write errors
type errorWriter struct{}

func (ew *errorWriter) Write(p []byte) (n int, err error) {
	return 0, errors.New("simulated write error")
}
