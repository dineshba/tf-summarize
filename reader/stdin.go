package reader

import (
	"os"
)

type StdinReader struct {
}

func (s StdinReader) Read() ([]byte, error) {
	return readFile(os.Stdin)
}

func NewStdinReader() Reader {
	return StdinReader{}
}
