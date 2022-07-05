package reader

import (
	"os"
)

const StdinFileName = "stdin"

type StdinReader struct {
}

func (s StdinReader) Name() string {
	return StdinFileName
}

func (s StdinReader) Read() ([]byte, error) {
	return readFile(os.Stdin)
}

func NewStdinReader() Reader {
	return StdinReader{}
}
