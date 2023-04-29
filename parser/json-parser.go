package parser

import (
	"encoding/json"
	"fmt"

	"github.com/dineshba/tf-summarize/terraformstate"
)

type JSONParser struct {
	data []byte
}

func (j JSONParser) Parse() (terraformstate.TerraformState, error) {
	ts := terraformstate.TerraformState{}
	err := json.Unmarshal(j.data, &ts)
	if err != nil {
		return terraformstate.TerraformState{}, fmt.Errorf("error when parsing input: %s", err.Error())
	}
	return ts, nil
}

func NewJSONParser(data []byte) Parser {
	return JSONParser{
		data: data,
	}
}
