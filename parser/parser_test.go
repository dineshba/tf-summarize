package parser

import (
	"testing"

	"github.com/dineshba/tf-summarize/reader"
	"github.com/stretchr/testify/assert"
)

// Test for JSON file input
func TestCreateParser_JSONFile(t *testing.T) {
	data := []byte(`{"plan": "mock"}`)
	fileName := "example.json"

	p, err := CreateParser(data, fileName)

	assert.NoError(t, err)

	_, ok := p.(JSONParser)
	assert.True(t, ok, "expected a JSON parser to be returned")
}

// Test for stdin input
func TestCreateParser_Stdin(t *testing.T) {
	data := []byte(`{"plan": "mock"}`)
	fileName := reader.StdinFileName

	p, err := CreateParser(data, fileName)

	assert.NoError(t, err)

	_, ok := p.(JSONParser)
	assert.True(t, ok, "expected a JSON parser for stdin input")
}

// Test for binary file input
func TestCreateParser_BinaryFile(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02} // Mock binary data
	fileName := "example.binary"

	p, err := CreateParser(data, fileName)

	assert.NoError(t, err)

	_, ok := p.(BinaryParser)
	assert.True(t, ok, "expected a Binary parser to be returned")
}

// Test for non-JSON file name (like .txt or other extensions)
func TestCreateParser_InvalidFileName(t *testing.T) {
	data := []byte(`irrelevant data`)
	fileName := "example.txt"

	p, err := CreateParser(data, fileName)

	assert.NoError(t, err)

	_, ok := p.(BinaryParser)
	assert.True(t, ok, "expected a Binary parser for non-JSON file extension")
}
