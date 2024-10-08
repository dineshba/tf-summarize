package writer

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/dineshba/tf-summarize/terraformstate"
	"github.com/stretchr/testify/assert"

	. "github.com/hashicorp/terraform-json"
	tfjson "github.com/hashicorp/terraform-json"
)

func TestSeparateTree_Write_DrawableTrue(t *testing.T) {

	changes := map[string]terraformstate.ResourceChanges{
		"add": {
			{
				Address: "aws_instance.example1",
				Change:  &Change{Actions: Actions{ActionCreate}},
			},
		},
		"delete": {
			{
				Address: "aws_instance.example2",
				Change:  &Change{Actions: Actions{ActionDelete}},
			},
		},
	}

	tw := NewSeparateTree(changes, true)
	var buf bytes.Buffer
	err := tw.Write(&buf)

	assert.NoError(t, err)
	expectedOutput := `################### ADD ###################
      ╭─╮      
      │.│      
      ╰┬╯      
       │       
╭──────┴─────╮ 
│aws_instance│ 
╰──────┬─────╯ 
       │       
 ╭─────┴─────╮ 
 │example1(+)│ 
 ╰───────────╯ 


################### DELETE ###################
      ╭─╮      
      │.│      
      ╰┬╯      
       │       
╭──────┴─────╮ 
│aws_instance│ 
╰──────┬─────╯ 
       │       
 ╭─────┴─────╮ 
 │example2(-)│ 
 ╰───────────╯ 


`

	assert.Equal(t, expectedOutput, removeANSI(buf.String()))
}

func TestSeparateTree_Write_NonDrawable(t *testing.T) {

	changes := map[string]terraformstate.ResourceChanges{
		"add": {
			{
				Address: "aws_instance.example1",
				Change:  &Change{Actions: Actions{ActionCreate}},
			},
		},
		"delete": {
			{
				Address: "aws_instance.example2",
				Change:  &Change{Actions: Actions{ActionDelete}},
			},
		},
	}

	tw := NewSeparateTree(changes, false)
	var buf bytes.Buffer
	err := tw.Write(&buf)

	assert.NoError(t, err)

	expectedOutput := `################### ADD ###################
|---aws_instance
|	|---example1(+)


################### DELETE ###################
|---aws_instance
|	|---example2(-)


`

	assert.Equal(t, expectedOutput, removeANSI(buf.String()))
}

// Mock writer that returns an error after a certain number of writes
type controlledErrorWriter struct {
	failAfter int // Number of successful writes before error
	writes    int // Counter for number of writes
}

func (w *controlledErrorWriter) Write(p []byte) (n int, err error) {
	w.writes++
	if w.writes > w.failAfter {
		return 0, errors.New("write error")
	}
	return len(p), nil
}

// Mock tree writer to simulate an error during treeWriter.Write
type mockTreeWriterWithError struct{}

func (m *mockTreeWriterWithError) Write(writer io.Writer) error {
	return errors.New("tree writer error")
}

// Function to create a mock TreeWriter with an error
func NewMockTreeWriterWithError(changes []*tfjson.ResourceChange, drawable bool) Writer {
	return &mockTreeWriterWithError{}
}

type mockTreeWriter struct{}

func (m *mockTreeWriter) Write(writer io.Writer) error {
	return nil
}

func NewMockTreeWriter(changes []*tfjson.ResourceChange, drawable bool) Writer {
	return &mockTreeWriter{}
}

func TestSeparateTree_Write_Error(t *testing.T) {
	changes := map[string]terraformstate.ResourceChanges{
		"add": {
			{
				Address: "aws_instance.example1",
				Change:  &Change{Actions: Actions{ActionCreate}},
			},
		},
		"delete": {
			{
				Address: "aws_instance.example2",
				Change:  &Change{Actions: Actions{ActionDelete}},
			},
		},
	}

	s := NewSeparateTree(changes, false)

	t.Run("Test fmt.Fprintf error for section header", func(t *testing.T) {
		writer := &controlledErrorWriter{failAfter: 0}
		err := s.Write(writer)

		if err == nil || !strings.Contains(err.Error(), "write error") {
			t.Errorf("expected write error for section header, got %v", err)
		}
	})

	t.Run("Test treeWriter.Write returns error", func(t *testing.T) {
		originalFunc := NewTreeWriterFunc
		defer func() { NewTreeWriterFunc = originalFunc }() // Ensure restoration after the test

		// Replace NewTreeWriter with the mock
		NewTreeWriterFunc = NewMockTreeWriterWithError

		writer := &strings.Builder{}
		err := s.Write(writer)

		if err == nil || !strings.Contains(err.Error(), "tree writer error") {
			t.Errorf("expected tree writer error, got %v", err)
		}
	})

	t.Run("Test fmt.Fprintf error on final newline", func(t *testing.T) {

		// Backup the original NewTreeWriterFunc to restore after this test
		originalFunc := NewTreeWriterFunc
		defer func() { NewTreeWriterFunc = originalFunc }() // Ensure restoration after the test

		// Replace NewTreeWriter with the mock
		NewTreeWriterFunc = NewMockTreeWriter

		// Use controlledErrorWriter to simulate write error after 1 write
		writer := &controlledErrorWriter{failAfter: 1}
		err := s.Write(writer)

		if err == nil || !strings.Contains(err.Error(), "write error") {
			t.Errorf("expected write error on newline, got %v", err)
		}
	})
}
