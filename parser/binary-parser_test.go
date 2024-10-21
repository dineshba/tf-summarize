package parser

import (
	"encoding/json"
	"testing"

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
	// Prepare mock output
	mockPlan := tfjson.Plan{
		FormatVersion:    "0.1",
		TerraformVersion: "1.0.0",
	}
	output, err := json.Marshal(mockPlan)
	if err != nil {
		t.Fatalf("Failed to marshal mock plan: %v", err)
	}

	// Create parser with mock executor
	parser := BinaryParser{
		fileName: "mock-file",
		executor: &MockCommandExecutor{
			Output: output,
			Err:    nil,
		},
	}

	// Call Parse
	parsedPlan, err := parser.Parse()

	//assertions
	assert.NoError(t, err, "Expected no error")
	assert.Equal(t, mockPlan, parsedPlan, "Parsed plan does not match expected plan")
}

func TestBinaryParser_Parse_EmptyOutput(t *testing.T) {
	// Simulate the command returning empty output
	parser := BinaryParser{
		fileName: "mock-file",
		executor: &MockCommandExecutor{
			Output: []byte(""), // Empty output
			Err:    nil,
		},
	}

	// Call Parse
	_, err := parser.Parse()

	//assertions
	assert.Error(t, err, "Expected error due to empty output")
	assert.Contains(t, err.Error(), "error when parsing input", "Unexpected error message")
}
func TestBinaryParser_Parse_MissingRequiredFields(t *testing.T) {
	// Prepare mock output with missing required fields
	incompleteMockOutput := map[string]interface{}{
		// Removing "format_version" and "terraform_version"
		"variables": map[string]interface{}{},
	}
	output, err := json.Marshal(incompleteMockOutput)
	if err != nil {
		t.Fatalf("Failed to marshal incomplete plan: %v", err)
	}

	// Create parser with mock executor
	parser := BinaryParser{
		fileName: "mock-file",
		executor: &MockCommandExecutor{
			Output: output,
			Err:    nil,
		},
	}

	// Call Parse
	parsedPlan, err := parser.Parse()

	// Assertions
	assert.Error(t, err, "Expected an error due to missing required fields")
	assert.Contains(t, err.Error(), "format version is missing", "Unexpected error message")
	assert.Equal(t, "", parsedPlan.TerraformVersion, "Expected empty TerraformVersion")
}

func TestBinaryParser_Parse_InvalidJSONOutput(t *testing.T) {
	// Simulate the command returning invalid JSON
	invalidJSON := []byte(`{"invalid-json"}`)
	parser := BinaryParser{
		fileName: "mock-file",
		executor: &MockCommandExecutor{
			Output: invalidJSON,
			Err:    nil,
		},
	}

	// Call Parse
	_, err := parser.Parse()

	//assertions
	assert.Error(t, err, "Expected error due to invalid JSON output")
	assert.Contains(t, err.Error(), "error when parsing input", "Unexpected error message")
}
