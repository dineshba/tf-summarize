package writer

import (
	"embed"
	"regexp"
)

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

// Function to remove ANSI escape sequences
func removeANSI(input string) string {
	re := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	return re.ReplaceAllString(input, "")
}
