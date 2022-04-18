package parser

import (
	"encoding/json"
	"fmt"
	"terraform-plan-summary/terraform_state"
)

type JsonParser struct {
	data []byte
}

func (j JsonParser) Parse() (terraform_state.TerraformState, error) {
	ts := terraform_state.TerraformState{}
	err := json.Unmarshal(j.data, &ts)
	if err != nil {
		return terraform_state.TerraformState{}, fmt.Errorf("error when parsing input: %s", err.Error())
	}
	return ts, nil
}

func NewJsonParser(data []byte) Parser {
	return JsonParser{
		data: data,
	}
}
