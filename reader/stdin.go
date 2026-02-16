package reader

import (
	"os"
)

// StdinFileName is the name reported by StdinReader.
const StdinFileName = "stdin"

// StdinReader reads Terraform plan JSON from standard input.
type StdinReader struct {
}

// Name returns "stdin".
func (s StdinReader) Name() string {
	return StdinFileName
}

func (s StdinReader) Read() ([]byte, error) {
	return readFile(os.Stdin)
}

// NewStdinReader returns a Reader that reads from STDIN.
func NewStdinReader() Reader {
	return StdinReader{}
}
