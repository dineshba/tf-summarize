// binary_parser_test.go
package parser

import (
	"encoding/json"
	"testing"

	"github.com/dineshba/tf-summarize/parser/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	tfjson "github.com/hashicorp/terraform-json"
)

// MockCommandExecutor mocks the command execution.
type MockCommandExecutor struct {
	Output     []byte
	Err        error
	ActualName string
	ActualArgs []string
}

func (e *MockCommandExecutor) CombinedOutput(name string, args ...string) ([]byte, error) {
	e.ActualName = name
	e.ActualArgs = args
	return e.Output, e.Err
}

func TestBinaryParser_Parse_Success(t *testing.T) {
	// Initialize GoMock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Prepare mock output
	mockPlan := tfjson.Plan{
		FormatVersion:    "0.1",
		TerraformVersion: "1.0.0",
	}
	output, err := json.Marshal(mockPlan)
	assert.NoError(t, err, "Failed to marshal mock plan")

	// Create mock executor
	mockExecutor := mocks.NewMockCommandExecutor(ctrl)
	tfbinary := "terraform"

	// Set up expected call
	mockExecutor.
		EXPECT().
		CombinedOutput(tfbinary, "show", "-json", "mock-file").
		Return(output, nil)

	// Create parser with mock executor
	parser := BinaryParser{
		fileName: "mock-file",
		executor: mockExecutor,
	}

	// Call Parse
	parsedPlan, err := parser.Parse()

	// assertions
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, mockPlan, parsedPlan)
}
func TestBinaryParser_Parse_EmptyOutput(t *testing.T) {
	// Initialize GoMock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock executor
	mockExecutor := mocks.NewMockCommandExecutor(ctrl)
	tfbinary := "terraform"

	// Set up expectation CombinedOutput will return empty output
	mockExecutor.
		EXPECT().
		CombinedOutput(tfbinary, "show", "-json", "mock-file").
		Return([]byte(""), nil)

	// Create parser with mock executor
	parser := BinaryParser{
		fileName: "mock-file",
		executor: mockExecutor,
	}

	// Call Parse
	_, err := parser.Parse()

	// assertions
	assert.Error(t, err, "Expected error due to empty output")
	assert.Contains(t, err.Error(), "error when parsing input")
}

func TestBinaryParser_Parse_MissingRequiredFields(t *testing.T) {
	// Initialize GoMock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Prepare mock output with missing required fields
	incompleteMockOutput := map[string]interface{}{
		// format_version and terraform_version are missing
		"variables": map[string]interface{}{},
	}
	output, err := json.Marshal(incompleteMockOutput)
	assert.NoError(t, err, "Failed to marshal incomplete plan")

	// Create mock executor
	mockExecutor := mocks.NewMockCommandExecutor(ctrl)
	tfbinary := "terraform"

	// Set up expectation: CombinedOutput should return incomplete output
	mockExecutor.
		EXPECT().
		CombinedOutput(tfbinary, "show", "-json", "mock-file").
		Return(output, nil)

	// Create parser with mock executor
	parser := BinaryParser{
		fileName: "mock-file",
		executor: mockExecutor,
	}

	// Call Parse
	parsedPlan, err := parser.Parse()

	// assertions
	assert.Error(t, err, "Expected an error due to missing required fields")
	assert.Contains(t, err.Error(), "format version is missing")
	assert.Equal(t, "", parsedPlan.TerraformVersion)
}

func TestBinaryParser_Parse_InvalidJSONOutput(t *testing.T) {
	// Initialize GoMock controller
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Prepare mock output with invalid json output
	invalidJSON := []byte(`{"invalid-json"}`)
	mockExecutor := mocks.NewMockCommandExecutor(ctrl)
	tfbinary := "terraform"

	// Set up expectation
	mockExecutor.
		EXPECT().
		CombinedOutput(tfbinary, "show", "-json", "mock-file").
		Return(invalidJSON, nil)

	// Create parser with mock executor
	parser := BinaryParser{
		fileName: "mock-file",
		executor: mockExecutor,
	}

	// Call Parse
	_, err := parser.Parse()

	// Assertions
	assert.Error(t, err, "Expected error due to invalid JSON output")
	assert.Contains(t, err.Error(), "error when parsing input")
}
