package terraformstate

import (
	"testing"

	. "github.com/hashicorp/terraform-json"

	"github.com/stretchr/testify/assert"
)

func TestResourceChangeColor(t *testing.T) {
	ExpectedColors := map[Action]string{
		ActionCreate: ColorGreen,
		ActionDelete: ColorRed,
		ActionUpdate: ColorYellow,
	}

	for action, expectedColor := range ExpectedColors {
		create := &ResourceChange{Change: &Change{Actions: []Action{action}}}
		color, _ := GetColorPrefixAndSuffixText(create)

		assert.Equal(t, color, expectedColor)
	}
	create_delete := &ResourceChange{Change: &Change{Actions: []Action{ActionCreate, ActionDelete}}}
	color, _ := GetColorPrefixAndSuffixText(create_delete)
	assert.Equal(t, color, ColorMagenta)

	delete_create := &ResourceChange{Change: &Change{Actions: []Action{ActionDelete, ActionCreate}}}
	color, _ = GetColorPrefixAndSuffixText(delete_create)
	assert.Equal(t, color, ColorMagenta)
}

func TestResourceChangeSuffix(t *testing.T) {
	ExpectedSuffix := map[Action]string{
		ActionCreate: "(+)",
		ActionDelete: "(-)",
		ActionUpdate: "(~)",
	}

	for action, expectedSuffix := range ExpectedSuffix {
		create := &ResourceChange{Change: &Change{Actions: []Action{action}}}
		_, suffix := GetColorPrefixAndSuffixText(create)

		assert.Equal(t, suffix, expectedSuffix)
	}
	create_delete := &ResourceChange{Change: &Change{Actions: []Action{ActionCreate, ActionDelete}}}
	_, suffix := GetColorPrefixAndSuffixText(create_delete)
	assert.Equal(t, suffix, "(+/-)")

	delete_create := &ResourceChange{Change: &Change{Actions: []Action{ActionDelete, ActionCreate}}}
	_, suffix = GetColorPrefixAndSuffixText(delete_create)
	assert.Equal(t, suffix, "(-/+)")
}

func TestResourceChangeColorAndSuffixImport(t *testing.T) {
	importing := &ResourceChange{Change: &Change{Importing: &Importing{ID: "id"}}}
	color, suffix := GetColorPrefixAndSuffixText(importing)

	assert.Equal(t, color, ColorCyan)
	assert.Equal(t, suffix, "(i)")
}
