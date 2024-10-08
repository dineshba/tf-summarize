package writer

import (
	"io"
	"path"
	"text/template"

	"github.com/dineshba/tf-summarize/terraformstate"
)

// HTMLWriter is a Writer that writes HTML.
type HTMLWriter struct {
	ResourceChanges map[string]terraformstate.ResourceChanges
	OutputChanges   map[string][]string
}

var cfs = getFS()

// Write outputs the HTML summary to the io.Writer it's passed.
func (t HTMLWriter) Write(writer io.Writer) error {
	templatesDir := "templates"
	rcTmpl := "resourceChanges.html"
	tmpl, err := template.New(rcTmpl).ParseFS(cfs, path.Join(templatesDir, rcTmpl))
	if err != nil {
		return err
	}

	err = tmpl.Execute(writer, t)
	if err != nil {
		return err
	}

	if !hasOutputChanges(t.OutputChanges) {
		return nil
	}

	ocTmpl := "outputChanges.html"
	outputTmpl, err := template.New(ocTmpl).ParseFS(cfs, path.Join(templatesDir, ocTmpl))
	if err != nil {
		return err
	}

	return outputTmpl.Execute(writer, t)
}

// NewHTMLWriter returns a new HTMLWriter with the configuration it's passed.
func NewHTMLWriter(changes map[string]terraformstate.ResourceChanges, outputChanges map[string][]string) Writer {
	return HTMLWriter{
		ResourceChanges: changes,
		OutputChanges:   outputChanges,
	}
}
