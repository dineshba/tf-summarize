package parser

import (
	"encoding/json"
	"fmt"

	tfjson "github.com/hashicorp/terraform-json"
)

type JSONParser struct {
	data []byte
}

func (j JSONParser) Parse() (tfjson.Plan, error) {
	ts := tfjson.Plan{}
	err := json.Unmarshal(j.data, &ts)
	if err != nil {
		return tfjson.Plan{}, fmt.Errorf("error when parsing input: %s", err.Error())
	}
	return ts, nil
}

func NewJSONParser(data []byte) Parser {
	return JSONParser{
		data: data,
	}
}
