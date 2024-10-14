package parser

import (
	"encoding/json"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
)

func TestJSONParser_Parse_Success(t *testing.T) {
	mockPlan := tfjson.Plan{
		FormatVersion:    "0.1",
		TerraformVersion: "1.0.0",
		Variables:        map[string]*tfjson.PlanVariable{"example": {}},
	}
	data, err := json.Marshal(mockPlan)
	assert.NoError(t, err)

	parser := NewJSONParser(data)
	parsedPlan, err := parser.Parse()

	assert.NoError(t, err)
	assert.Equal(t, mockPlan.FormatVersion, parsedPlan.FormatVersion)
	assert.Equal(t, mockPlan.TerraformVersion, parsedPlan.TerraformVersion)
	assert.Equal(t, mockPlan.Variables, parsedPlan.Variables)
}

func TestJSONParser_Parse_InvalidJSON(t *testing.T) {
	// invalid JSON data
	invalidData := []byte(`{"invalid-json"}`)

	parser := NewJSONParser(invalidData)
	_, err := parser.Parse()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error when parsing input")
}

func TestNewJSONParser(t *testing.T) {
	data := []byte(`{"plan": "mock"}`)

	parser := NewJSONParser(data)
	jp, ok := parser.(JSONParser)
	assert.True(t, ok, "expected a JSONParser instance")
	assert.Equal(t, data, jp.data, "expected data to match the input")
}
