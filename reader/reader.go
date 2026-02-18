// Package reader provides abstractions for reading Terraform plan JSON input
// from files or STDIN.
package reader

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

// Reader reads Terraform plan JSON input and returns the raw bytes.
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

// CreateReader returns a Reader based on the provided args. If a single
// filename argument is given, it returns a FileReader; otherwise it returns a
// StdinReader.
func CreateReader(args []string) (Reader, error) {
	if len(args) > 1 {
		return nil, fmt.Errorf("expected input via a single filename argument or via STDIN; received multiple arguments: %s", strings.Join(args, ", "))
	}

	if len(args) == 1 {
		return NewFileReader(args[0]), nil
	}

	return NewStdinReader(), nil
}
