package reader

import (
	"bufio"
	"fmt"
	"io"
)

type Reader interface {
	Read() ([]byte, error)
}

func readFile(f io.Reader) ([]byte, error) {
	var input []byte
	r := bufio.NewReader(f)
	var err error
	var line []byte
	for err == nil {
		line, err = r.ReadBytes('\n')
		input = append(input, line...)
	}
	if err != io.EOF {
		return nil, fmt.Errorf("error reading file: %s", err.Error())
	}
	return input, nil
}
