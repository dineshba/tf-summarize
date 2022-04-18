package reader

import (
	"bufio"
	"fmt"
	"io"
	"os"
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

func CreateReader(stdin *os.File, args []string) (Reader, error) {
	stat, _ := stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return NewStdinReader(), nil
	}
	if len(args) < 1 {
		return nil, fmt.Errorf("should have either stdin input through pipe or first argument should be file")
	}
	fileName := args[0]
	return NewFileReader(fileName), nil
}
