// Package parser provides parsers for Terraform plan output in JSON and binary formats.
package parser

import (
	"strings"

	"github.com/dineshba/tf-summarize/reader"
	tfjson "github.com/hashicorp/terraform-json"
)

// Parser parses Terraform plan data into a structured Plan.
type Parser interface {
	Parse() (tfjson.Plan, error)
}

// CreateParser returns a Parser appropriate for the given file type.
func CreateParser(data []byte, fileName string) (Parser, error) {
	if fileName != reader.StdinFileName && !strings.HasSuffix(fileName, ".json") {
		return NewBinaryParser(fileName), nil
	}
	return NewJSONParser(data), nil
}
