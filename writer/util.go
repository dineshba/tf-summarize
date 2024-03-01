package writer

import "embed"

// Embed the templates directory in the compiled binary.
//
//go:embed templates
var templates embed.FS

func hasOutputChanges(opChanges map[string][]string) bool {
	hasChanges := false

	for _, v := range opChanges {
		if len(v) > 0 {
			hasChanges = true
			break
		}
	}

	return hasChanges
}
