// Package writer provides output formatters for Terraform plan summaries.
// The JSON pretty print code is inspired from https://github.com/hokaccha/go-prettyjson.
package writer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

// Formatter controls JSON pretty-printing with color support.
type Formatter struct {
	AddColor        *color.Color
	RemoveColor     *color.Color
	UpdateColor     *color.Color
	RecreateColor   *color.Color
	ImportColor     *color.Color
	StringMaxLength int
	Indent          int
	Newline         string
}

// NewFormatter returns a new formatter with following default values.
func NewFormatter() *Formatter {
	return &Formatter{
		AddColor:        color.New(color.FgGreen, color.Bold),
		RemoveColor:     color.New(color.FgRed, color.Bold),
		UpdateColor:     color.New(color.FgYellow, color.Bold),
		RecreateColor:   color.New(color.FgMagenta, color.Bold),
		ImportColor:     color.New(color.FgCyan, color.Bold),
		StringMaxLength: 0,
		Indent:          2,
		Newline:         "\n",
	}
}

// Marshal marshals and formats JSON data.
func (f *Formatter) Marshal(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	return f.Format(data)
}

// Format formats JSON string.
func (f *Formatter) Format(data []byte) ([]byte, error) {
	var v interface{}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&v); err != nil {
		return nil, err
	}

	return []byte(f.pretty(v, 1)), nil
}

func (f *Formatter) pretty(v interface{}, depth int) string {
	switch val := v.(type) {
	case string:
		return f.processString(val)
	case float64:
		return strconv.FormatFloat(val, 'f', -1, 64)
	case json.Number:
		return string(val)
	case bool:
		return strconv.FormatBool(val)
	case nil:
		return "null"
	case map[string]interface{}:
		return f.processMap(val, depth)
	case []interface{}:
		return f.processArray(val, depth)
	}

	return ""
}

func (f *Formatter) processString(s string) string {
	r := []rune(s)
	if f.StringMaxLength != 0 && len(r) >= f.StringMaxLength {
		s = string(r[0:f.StringMaxLength]) + "..."
	}

	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(s)
	s = buf.String()
	s = strings.TrimSuffix(s, "\n")

	return s
}

func (f *Formatter) processMap(m map[string]interface{}, depth int) string {
	if len(m) == 0 {
		return "{}"
	}

	currentIndent := f.generateIndent(depth - 1)
	nextIndent := f.generateIndent(depth)
	rows := []string{}
	keys := []string{}

	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		val := m[key]
		buf := &bytes.Buffer{}
		encoder := json.NewEncoder(buf)
		encoder.SetEscapeHTML(false)
		_ = encoder.Encode(key)
		k := strings.TrimSuffix(buf.String(), "\n")
		k = strings.ReplaceAll(k, "[\\\"", "[")
		k = strings.ReplaceAll(k, "\\\"]", "]")
		v := f.pretty(val, depth+1)

		// Add color based on the key value
		switch key {
		case "(+)":
			v = f.AddColor.SprintFunc()(v)
		case "(-)":
			v = f.RemoveColor.SprintFunc()(v)
		case "(~)":
			v = f.UpdateColor.SprintFunc()(v)
		case "(+/-)":
			v = f.RecreateColor.SprintFunc()(v)
		}

		valueIndent := " "
		if f.Newline == "" {
			valueIndent = ""
		}
		row := fmt.Sprintf("%s%s:%s%s", nextIndent, k, valueIndent, v)
		rows = append(rows, row)
	}

	return fmt.Sprintf("{%s%s%s%s}", f.Newline, strings.Join(rows, ","+f.Newline), f.Newline, currentIndent)
}

func (f *Formatter) processArray(a []interface{}, depth int) string {
	if len(a) == 0 {
		return "[]"
	}

	currentIndent := f.generateIndent(depth - 1)
	nextIndent := f.generateIndent(depth)
	rows := []string{}

	for _, val := range a {
		c := f.pretty(val, depth+1)
		row := nextIndent + c
		rows = append(rows, row)
	}
	return fmt.Sprintf("[%s%s%s%s]", f.Newline, strings.Join(rows, ","+f.Newline), f.Newline, currentIndent)
}

func (f *Formatter) generateIndent(depth int) string {
	return strings.Repeat(" ", f.Indent*depth)
}

// Marshal JSON data with default options.
func Marshal(v interface{}) ([]byte, error) {
	return NewFormatter().Marshal(v)
}

// Format JSON string with default options.
func Format(data []byte) ([]byte, error) {
	return NewFormatter().Format(data)
}
