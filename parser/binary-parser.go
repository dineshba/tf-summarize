package parser

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	tfjson "github.com/hashicorp/terraform-json"
)

// BinaryParser parses a Terraform binary plan file by invoking the terraform CLI.
type BinaryParser struct {
	fileName string
}

// Parse runs terraform show -json on the binary plan file and returns the parsed Plan.
func (j BinaryParser) Parse() (tfjson.Plan, error) {
	tfbinary := "terraform"
	cmdArgs := []string{"show", "-json", j.fileName}
	if tfoverride, ok := os.LookupEnv("TF_BINARY"); ok {
		if tfoverride == "terragrunt" {
			cmdArgs = append(cmdArgs, "--terragrunt-log-disable")
		}
		tfbinary = tfoverride
	}
	cmd := exec.Command(tfbinary, cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return tfjson.Plan{}, fmt.Errorf(
			"error when running '%s show -json %s': \n%s\n\nMake sure you are running in %s directory and %s init is done",
			tfbinary, j.fileName, output, tfbinary, tfbinary)
	}
	plan := tfjson.Plan{}
	err = json.Unmarshal(output, &plan)
	if err != nil {
		return tfjson.Plan{}, fmt.Errorf("error when parsing input: %s", err.Error())
	}
	return plan, nil
}

// NewBinaryParser returns a new BinaryParser for the given file.
func NewBinaryParser(fileName string) Parser {
	return BinaryParser{
		fileName: fileName,
	}
}
