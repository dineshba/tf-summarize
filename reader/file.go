package reader

import (
	"fmt"
	"os"
)

type FileReader struct {
	fileName string
}

func (f FileReader) Read() ([]byte, error) {
	file, err := os.Open(f.fileName)
	if err != nil {
		return nil, fmt.Errorf("error when opening file %s: %s", f.fileName, err)
	}
	defer file.Close()

	return readFile(file)
}

func NewFileReader(name string) Reader {
	return FileReader{fileName: name}
}
