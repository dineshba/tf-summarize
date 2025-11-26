package writer

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/dineshba/tf-summarize/terraformstate"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/dineshba/tf-summarize/testdata/mocks"

	. "github.com/hashicorp/terraform-json"
	tfjson "github.com/hashicorp/terraform-json"
)

// Mock tree writer
type mockTreeWriter struct{}

func (m *mockTreeWriter) Write(writer io.Writer) error {
	return nil
}

func NewMockTreeWriter(changes []*tfjson.ResourceChange, drawable bool) Writer {
	return &mockTreeWriter{}
}

// Mock tree writer to simulate an error during treeWriter.Write
type mockTreeWriterWithError struct{}

func (m *mockTreeWriterWithError) Write(writer io.Writer) error {
	return errors.New("tree writer error")
}

func NewMockTreeWriterWithError(changes []*tfjson.ResourceChange, drawable bool) Writer {
	return &mockTreeWriterWithError{}
}

// Helper function to create changes
func createMockChanges() map[string]terraformstate.ResourceChanges {
	return map[string]terraformstate.ResourceChanges{
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
}

func TestSeparateTree_Write(t *testing.T) {
	mockChanges := createMockChanges()

	t.Run("Drawable True", func(t *testing.T) {
		tw := NewSeparateTree(mockChanges, true)
		var buf bytes.Buffer
		err := tw.Write(&buf)

		assert.NoError(t, err)

		expectedAdd := `################### ADD ###################
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
`
		expectedDelete := `################### DELETE ###################
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

		actualOutput := removeANSI(buf.String())
		assert.Contains(t, actualOutput, expectedAdd)
		assert.Contains(t, actualOutput, expectedDelete)

	})

	t.Run("Drawable False", func(t *testing.T) {
		tw := NewSeparateTree(mockChanges, false)
		var buf bytes.Buffer
		err := tw.Write(&buf)

		assert.NoError(t, err)

		expectedAdd := `################### ADD ###################
|---aws_instance
|	|---example1(+)
`
		expectedDelete := `################### DELETE ###################
|---aws_instance
|	|---example2(-)
`

		actualOutput := removeANSI(buf.String())
		assert.Contains(t, actualOutput, expectedAdd)
		assert.Contains(t, actualOutput, expectedDelete)
	})

	t.Run("Error Handling", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		s := NewSeparateTree(mockChanges, false)

		t.Run("Write Error", func(t *testing.T) {
			mockWriter := mocks.NewMockWriter(ctrl)
			mockWriter.EXPECT().Write(gomock.Any()).Return(0, errors.New("write error")).Times(1)

			err := s.Write(mockWriter)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "write error")
		})

		t.Run("Tree Writer Error", func(t *testing.T) {
			originalFunc := NewTreeWriterFunc
			defer func() { NewTreeWriterFunc = originalFunc }()

			// Replace NewTreeWriter with the mock that returns an error
			NewTreeWriterFunc = NewMockTreeWriterWithError

			mockWriter := mocks.NewMockWriter(ctrl)
			mockWriter.EXPECT().Write(gomock.Any()).Return(0, nil).AnyTimes()

			err := s.Write(mockWriter)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "tree writer error")
		})

		t.Run("Final Newline Write Error", func(t *testing.T) {
			originalFunc := NewTreeWriterFunc
			defer func() { NewTreeWriterFunc = originalFunc }()

			NewTreeWriterFunc = NewMockTreeWriter

			mockWriter := mocks.NewMockWriter(ctrl)
			mockWriter.EXPECT().Write(gomock.Any()).Return(0, nil).Times(1)
			mockWriter.EXPECT().Write(gomock.Any()).Return(0, errors.New("write error")).Times(1)

			err := s.Write(mockWriter)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "write error")
		})
	})
}
