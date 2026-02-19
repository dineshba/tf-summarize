package parser

import (
	"encoding/json"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
)

// JSONParser parses a JSON-encoded Terraform plan.
type JSONParser struct {
	data []byte
}

// Parse unmarshals the JSON data into a Plan.
func (j JSONParser) Parse() (tfjson.Plan, error) {
	plan := tfjson.Plan{}
	err := json.Unmarshal(j.data, &plan)
	if err != nil {
		return tfjson.Plan{}, fmt.Errorf("error when parsing input: %s", err.Error())
	}
	return plan, nil
}

// NewJSONParser returns a new JSONParser for the given data.
func NewJSONParser(data []byte) Parser {
	return JSONParser{
		data: data,
	}
}
