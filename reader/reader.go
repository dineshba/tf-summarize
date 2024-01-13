package reader

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

type Reader interface {
	Read() ([]byte, error)
	Name() string
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
		return nil, fmt.Errorf("error reading input: %s", err.Error())
	}
	if len(input) == 0 {
		return nil, errors.New("no input data; expected input via a non-empty file or via STDIN")
	}
	return input, nil
}

func CreateReader(args []string) (Reader, error) {
	if len(args) > 1 {
		return nil, fmt.Errorf("expected input via a single filename argument or via STDIN; received multiple arguments: %s", strings.Join(args, ", "))
	}

	if len(args) == 1 {
		return NewFileReader(args[0]), nil
	}

	return NewStdinReader(), nil
}
