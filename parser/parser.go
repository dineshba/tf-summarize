package parser

import (
	"strings"

	"github.com/dineshba/tf-summarize/reader"
	tfjson "github.com/hashicorp/terraform-json"
)

type Parser interface {
	Parse() (tfjson.Plan, error)
}

func CreateParser(data []byte, fileName string) (Parser, error) {
	if fileName != reader.StdinFileName && !strings.HasSuffix(fileName, ".json") {
		return NewBinaryParser(fileName), nil
	}
	return NewJSONParser(data), nil
}
