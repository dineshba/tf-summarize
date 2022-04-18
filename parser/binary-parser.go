package parser

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"terraform-plan-summary/terraform_state"
)

type BinaryParser struct {
	fileName string
}

func (j BinaryParser) Parse() (terraform_state.TerraformState, error) {
	cmd := exec.Command("terraform", "show", "-json", j.fileName)
	output, err := cmd.Output()
	if err != nil {
		return terraform_state.TerraformState{}, fmt.Errorf("error when running terraform show -json %s: %s", j.fileName, err.Error())
	}
	ts := terraform_state.TerraformState{}
	err = json.Unmarshal(output, &ts)
	if err != nil {
		return terraform_state.TerraformState{}, fmt.Errorf("error when parsing input: %s", err.Error())
	}
	return ts, nil
}

func NewBinaryParser(fileName string) Parser {
	return BinaryParser{
		fileName: fileName,
	}
}
