package terraformstate

import (
	"testing"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/stretchr/testify/assert"
)

func TestResourceChangeColor(t *testing.T) {
	ExpectedColors := map[tfjson.Action]string{
		tfjson.ActionCreate: ColorGreen,
		tfjson.ActionDelete: ColorRed,
		tfjson.ActionUpdate: ColorYellow,
	}

	for action, expectedColor := range ExpectedColors {
		create := &tfjson.ResourceChange{Change: &tfjson.Change{Actions: []tfjson.Action{action}}}
		color, _ := GetColorPrefixAndSuffixText(create)

		assert.Equal(t, color, expectedColor)
	}

	CreateDelete := &tfjson.ResourceChange{Change: &tfjson.Change{Actions: []tfjson.Action{tfjson.ActionCreate, tfjson.ActionDelete}}}
	color, _ := GetColorPrefixAndSuffixText(CreateDelete)
	assert.Equal(t, color, ColorMagenta)

	DeleteCreate := &tfjson.ResourceChange{Change: &tfjson.Change{Actions: []tfjson.Action{tfjson.ActionDelete, tfjson.ActionCreate}}}
	color, _ = GetColorPrefixAndSuffixText(DeleteCreate)
	assert.Equal(t, color, ColorMagenta)
}

func TestGetAllResourceChanges(t *testing.T) {
	resourceChanges := ResourceChanges{
		&tfjson.ResourceChange{Address: "create2", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}}},
		&tfjson.ResourceChange{Address: "create1", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}}},
		&tfjson.ResourceChange{Address: "delete2", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionDelete}}},
		&tfjson.ResourceChange{Address: "delete1", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionDelete}}},
		&tfjson.ResourceChange{Address: "update2", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionUpdate}}},
		&tfjson.ResourceChange{Address: "update1", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionUpdate}}},
		&tfjson.ResourceChange{Address: "import2", Change: &tfjson.Change{Importing: &tfjson.Importing{ID: "id1"}}},
		&tfjson.ResourceChange{Address: "import1", Change: &tfjson.Change{Importing: &tfjson.Importing{ID: "id2"}}},
		&tfjson.ResourceChange{Address: "move1", PreviousAddress: "move", Change: &tfjson.Change{Actions: tfjson.Actions{}}},
		&tfjson.ResourceChange{Address: "recreate2", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionDelete, tfjson.ActionCreate}}},
		&tfjson.ResourceChange{Address: "recreate1", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionDelete, tfjson.ActionCreate}}},
	}
	plan := tfjson.Plan{ResourceChanges: resourceChanges}

	result := GetAllResourceChanges(plan)

	expectedResourceChanges := map[string]ResourceChanges{
		"add": {
			&tfjson.ResourceChange{Address: "create1", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}}},
			&tfjson.ResourceChange{Address: "create2", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}}},
		},
		"delete": {
			&tfjson.ResourceChange{Address: "delete1", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionDelete}}},
			&tfjson.ResourceChange{Address: "delete2", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionDelete}}},
		},
		"update": {
			&tfjson.ResourceChange{Address: "update1", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionUpdate}}},
			&tfjson.ResourceChange{Address: "update2", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionUpdate}}},
		},
		"recreate": {
			&tfjson.ResourceChange{Address: "recreate1", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionDelete, tfjson.ActionCreate}}},
			&tfjson.ResourceChange{Address: "recreate2", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionDelete, tfjson.ActionCreate}}},
		},
		"import": {
			&tfjson.ResourceChange{Address: "import1", Change: &tfjson.Change{Importing: &tfjson.Importing{ID: "id2"}}},
			&tfjson.ResourceChange{Address: "import2", Change: &tfjson.Change{Importing: &tfjson.Importing{ID: "id1"}}},
		},
	}

	assert.Equal(t, expectedResourceChanges, result)
}

func TestGetAllResourceMoves(t *testing.T) {
	resourceChanges := ResourceChanges{
		&tfjson.ResourceChange{Address: "create2", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}}},
		&tfjson.ResourceChange{Address: "create1", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}}},
		&tfjson.ResourceChange{Address: "delete2", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionDelete}}},
		&tfjson.ResourceChange{Address: "delete1", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionDelete}}},
		&tfjson.ResourceChange{Address: "update2", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionUpdate}}},
		&tfjson.ResourceChange{Address: "update1", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionUpdate}}},
		&tfjson.ResourceChange{Address: "import2", Change: &tfjson.Change{Importing: &tfjson.Importing{ID: "id1"}}},
		&tfjson.ResourceChange{Address: "import1", Change: &tfjson.Change{Importing: &tfjson.Importing{ID: "id2"}}},
		&tfjson.ResourceChange{Address: "move1", PreviousAddress: "move", Change: &tfjson.Change{Actions: tfjson.Actions{}}},
		&tfjson.ResourceChange{Address: "recreate2", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionDelete, tfjson.ActionCreate}}},
		&tfjson.ResourceChange{Address: "recreate1", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionDelete, tfjson.ActionCreate}}},
	}
	plan := tfjson.Plan{ResourceChanges: resourceChanges}

	result := GetAllResourceMoves(plan)

	expectedResourceMoves := map[string]ResourceChanges{
		"moved": {
			&tfjson.ResourceChange{Address: "move1", PreviousAddress: "move", Change: &tfjson.Change{Actions: tfjson.Actions{}}},
		},
	}

	assert.Equal(t, expectedResourceMoves, result)
}

func TestGetAllOutputChanges(t *testing.T) {

	outputChanges := map[string]*tfjson.Change{
		"create2": {Actions: tfjson.Actions{tfjson.ActionCreate}},
		"create1": {Actions: tfjson.Actions{tfjson.ActionCreate}},
		"delete2": {Actions: tfjson.Actions{tfjson.ActionDelete}},
		"delete1": {Actions: tfjson.Actions{tfjson.ActionDelete}},
		"update2": {Actions: tfjson.Actions{tfjson.ActionUpdate}},
		"update1": {Actions: tfjson.Actions{tfjson.ActionUpdate}},
	}

	plan := tfjson.Plan{OutputChanges: outputChanges}

	result := GetAllOutputChanges(plan)

	expectedResourceChanges := map[string][]string{
		"add":    {"create1", "create2"},
		"delete": {"delete1", "delete2"},
		"update": {"update1", "update2"},
	}

	assert.Equal(t, expectedResourceChanges, result)
}

func TestResourceChangeSuffix(t *testing.T) {
	ExpectedSuffix := map[tfjson.Action]string{
		tfjson.ActionCreate: "(+)",
		tfjson.ActionDelete: "(-)",
		tfjson.ActionUpdate: "(~)",
	}

	for action, expectedSuffix := range ExpectedSuffix {
		create := &tfjson.ResourceChange{Change: &tfjson.Change{Actions: []tfjson.Action{action}}}
		_, suffix := GetColorPrefixAndSuffixText(create)

		assert.Equal(t, suffix, expectedSuffix)
	}
	CreateDelete := &tfjson.ResourceChange{Change: &tfjson.Change{Actions: []tfjson.Action{tfjson.ActionCreate, tfjson.ActionDelete}}}
	_, suffix := GetColorPrefixAndSuffixText(CreateDelete)
	assert.Equal(t, suffix, "(+/-)")

	DeleteCreate := &tfjson.ResourceChange{Change: &tfjson.Change{Actions: []tfjson.Action{tfjson.ActionDelete, tfjson.ActionCreate}}}
	_, suffix = GetColorPrefixAndSuffixText(DeleteCreate)
	assert.Equal(t, suffix, "(-/+)")
}

func TestResourceChangeColorAndSuffixImport(t *testing.T) {
	importing := &tfjson.ResourceChange{Change: &tfjson.Change{Importing: &tfjson.Importing{ID: "id"}}}
	color, suffix := GetColorPrefixAndSuffixText(importing)

	assert.Equal(t, color, ColorCyan)
	assert.Equal(t, suffix, "(i)")
}

func TestFilterNoOpResources(t *testing.T) {
	identityImport := &tfjson.ResourceChange{Address: "no-op5", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionNoop}, Importing: &tfjson.Importing{Identity: struct{ Account string }{Account: "account ID"}}}}
	resourceChanges := ResourceChanges{
		&tfjson.ResourceChange{Address: "no-op1", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionNoop}}},
		&tfjson.ResourceChange{Address: "no-op3", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionNoop}, Importing: nil}},
		&tfjson.ResourceChange{Address: "no-op2", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionNoop}, Importing: &tfjson.Importing{ID: ""}}},
		&tfjson.ResourceChange{Address: "no-op4", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionNoop}, Importing: &tfjson.Importing{Identity: nil}}},
		&tfjson.ResourceChange{Address: "create", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}}},
		identityImport,
	}
	plan := tfjson.Plan{ResourceChanges: resourceChanges}

	FilterNoOpResources(&plan)

	expectedResourceChangesAfterFiltering := ResourceChanges{
		&tfjson.ResourceChange{Address: "create", Change: &tfjson.Change{Actions: tfjson.Actions{tfjson.ActionCreate}}},
		identityImport,
	}
	assert.Equal(t, expectedResourceChangesAfterFiltering, plan.ResourceChanges)
}
