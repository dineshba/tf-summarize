package terraformstate

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResourceChangeColor(t *testing.T) {
	ExpectedColors := map[string]string{
		"create": ColorGreen,
		"delete": ColorRed,
		"update": ColorYellow,
	}

	for action, expectedColor := range ExpectedColors {
		create := ResourceChange{Change: Change{Actions: []string{action}}}
		color, _ := create.ColorPrefixAndSuffixText()

		assert.Equal(t, color, expected_color)
	}
	create := ResourceChange{Change: Change{Actions: []string{"create", "delete"}}}
	color, _ := create.ColorPrefixAndSuffixText()
	assert.Equal(t, color, ColorMagenta)

	create = ResourceChange{Change: Change{Actions: []string{"delete", "create"}}}
	color, _ = create.ColorPrefixAndSuffixText()
	assert.Equal(t, color, ColorMagenta)
}

func TestResourceChangeSuffix(t *testing.T) {
	ExpectedSuffix := map[string]string{
		"create": "(+)",
		"delete": "(-)",
		"update": "(~)",
	}

	for action, expectedSuffix := range ExpectedSuffix {
		create := ResourceChange{Change: Change{Actions: []string{action}}}
		_, suffix := create.ColorPrefixAndSuffixText()

		assert.Equal(t, suffix, expected_suffix)
	}
	create := ResourceChange{Change: Change{Actions: []string{"create", "delete"}}}
	_, suffix := create.ColorPrefixAndSuffixText()
	assert.Equal(t, suffix, "(+/-)")

	create = ResourceChange{Change: Change{Actions: []string{"delete", "create"}}}
	_, suffix = create.ColorPrefixAndSuffixText()
	assert.Equal(t, suffix, "(-/+)")
}

func TestResourceChangeColorAndSuffixImport(t *testing.T) {
	importing := ResourceChange{Change: Change{Importing: Importing{ID: "id"}}}
	color, suffix := importing.ColorPrefixAndSuffixText()

	assert.Equal(t, color, ColorCyan)
	assert.Equal(t, suffix, "(i)")
}
