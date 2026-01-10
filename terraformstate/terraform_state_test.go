package terraformstate

import (
	"testing"

	. "github.com/hashicorp/terraform-json"
	tfjson "github.com/hashicorp/terraform-json"

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

	CreateDelete := &ResourceChange{Change: &Change{Actions: []Action{ActionCreate, ActionDelete}}}
	color, _ := GetColorPrefixAndSuffixText(CreateDelete)
	assert.Equal(t, color, ColorMagenta)

	DeleteCreate := &ResourceChange{Change: &Change{Actions: []Action{ActionDelete, ActionCreate}}}
	color, _ = GetColorPrefixAndSuffixText(DeleteCreate)
	assert.Equal(t, color, ColorMagenta)
}

func TestGetAllResourceChanges(t *testing.T) {
	resourceChanges := ResourceChanges{
		&ResourceChange{Address: "create2", Change: &Change{Actions: Actions{ActionCreate}}},
		&ResourceChange{Address: "create1", Change: &Change{Actions: Actions{ActionCreate}}},
		&ResourceChange{Address: "delete2", Change: &Change{Actions: Actions{ActionDelete}}},
		&ResourceChange{Address: "delete1", Change: &Change{Actions: Actions{ActionDelete}}},
		&ResourceChange{Address: "update2", Change: &Change{Actions: Actions{ActionUpdate}}},
		&ResourceChange{Address: "update1", Change: &Change{Actions: Actions{ActionUpdate}}},
		&ResourceChange{Address: "import2", Change: &Change{Importing: &Importing{ID: "id1"}}},
		&ResourceChange{Address: "import1", Change: &Change{Importing: &Importing{ID: "id2"}}},
		&ResourceChange{Address: "move1", PreviousAddress: "move", Change: &Change{Actions: Actions{}}},
		&ResourceChange{Address: "recreate2", Change: &Change{Actions: Actions{ActionDelete, ActionCreate}}},
		&ResourceChange{Address: "recreate1", Change: &Change{Actions: Actions{ActionDelete, ActionCreate}}},
	}
	plan := tfjson.Plan{ResourceChanges: resourceChanges}

	result := GetAllResourceChanges(plan)

	expectedResourceChanges := map[string]ResourceChanges{
		"add": {
			&ResourceChange{Address: "create1", Change: &Change{Actions: Actions{ActionCreate}}},
			&ResourceChange{Address: "create2", Change: &Change{Actions: Actions{ActionCreate}}},
		},
		"delete": {
			&ResourceChange{Address: "delete1", Change: &Change{Actions: Actions{ActionDelete}}},
			&ResourceChange{Address: "delete2", Change: &Change{Actions: Actions{ActionDelete}}},
		},
		"update": {
			&ResourceChange{Address: "update1", Change: &Change{Actions: Actions{ActionUpdate}}},
			&ResourceChange{Address: "update2", Change: &Change{Actions: Actions{ActionUpdate}}},
		},
		"recreate": {
			&ResourceChange{Address: "recreate1", Change: &Change{Actions: Actions{ActionDelete, ActionCreate}}},
			&ResourceChange{Address: "recreate2", Change: &Change{Actions: Actions{ActionDelete, ActionCreate}}},
		},
		"import": {
			&ResourceChange{Address: "import1", Change: &Change{Importing: &Importing{ID: "id2"}}},
			&ResourceChange{Address: "import2", Change: &Change{Importing: &Importing{ID: "id1"}}},
		},
	}

	assert.Equal(t, expectedResourceChanges, result)
}

func TestGetAllResourceMoves(t *testing.T) {
	resourceChanges := ResourceChanges{
		&ResourceChange{Address: "create2", Change: &Change{Actions: Actions{ActionCreate}}},
		&ResourceChange{Address: "create1", Change: &Change{Actions: Actions{ActionCreate}}},
		&ResourceChange{Address: "delete2", Change: &Change{Actions: Actions{ActionDelete}}},
		&ResourceChange{Address: "delete1", Change: &Change{Actions: Actions{ActionDelete}}},
		&ResourceChange{Address: "update2", Change: &Change{Actions: Actions{ActionUpdate}}},
		&ResourceChange{Address: "update1", Change: &Change{Actions: Actions{ActionUpdate}}},
		&ResourceChange{Address: "import2", Change: &Change{Importing: &Importing{ID: "id1"}}},
		&ResourceChange{Address: "import1", Change: &Change{Importing: &Importing{ID: "id2"}}},
		&ResourceChange{Address: "move1", PreviousAddress: "move", Change: &Change{Actions: Actions{}}},
		&ResourceChange{Address: "recreate2", Change: &Change{Actions: Actions{ActionDelete, ActionCreate}}},
		&ResourceChange{Address: "recreate1", Change: &Change{Actions: Actions{ActionDelete, ActionCreate}}},
	}
	plan := tfjson.Plan{ResourceChanges: resourceChanges}

	result := GetAllResourceMoves(plan)

	expectedResourceMoves := map[string]ResourceChanges{
		"moved": {
			&ResourceChange{Address: "move1", PreviousAddress: "move", Change: &Change{Actions: Actions{}}},
		},
	}

	assert.Equal(t, expectedResourceMoves, result)
}

func TestGetAllOutputChanges(t *testing.T) {

	outputChanges := map[string]*Change{
		"create2": {Actions: Actions{ActionCreate}},
		"create1": {Actions: Actions{ActionCreate}},
		"delete2": {Actions: Actions{ActionDelete}},
		"delete1": {Actions: Actions{ActionDelete}},
		"update2": {Actions: Actions{ActionUpdate}},
		"update1": {Actions: Actions{ActionUpdate}},
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
	CreateDelete := &ResourceChange{Change: &Change{Actions: []Action{ActionCreate, ActionDelete}}}
	_, suffix := GetColorPrefixAndSuffixText(CreateDelete)
	assert.Equal(t, suffix, "(+/-)")

	DeleteCreate := &ResourceChange{Change: &Change{Actions: []Action{ActionDelete, ActionCreate}}}
	_, suffix = GetColorPrefixAndSuffixText(DeleteCreate)
	assert.Equal(t, suffix, "(-/+)")
}

func TestResourceChangeColorAndSuffixImport(t *testing.T) {
	importing := &ResourceChange{Change: &Change{Importing: &Importing{ID: "id"}}}
	color, suffix := GetColorPrefixAndSuffixText(importing)

	assert.Equal(t, color, ColorCyan)
	assert.Equal(t, suffix, "(i)")
}

func TestFilterNoOpResources(t *testing.T) {
	identityImport := &ResourceChange{Address: "no-op5", Change: &Change{Actions: Actions{ActionNoop}, Importing: &Importing{Identity: struct{ Account string }{Account: "account ID"}}}}
	resourceChanges := ResourceChanges{
		&ResourceChange{Address: "no-op1", Change: &Change{Actions: Actions{ActionNoop}}},
		&ResourceChange{Address: "no-op3", Change: &Change{Actions: Actions{ActionNoop}, Importing: nil}},
		&ResourceChange{Address: "no-op2", Change: &Change{Actions: Actions{ActionNoop}, Importing: &Importing{ID: ""}}},
		&ResourceChange{Address: "no-op4", Change: &Change{Actions: Actions{ActionNoop}, Importing: &Importing{Identity: nil}}},
		&ResourceChange{Address: "create", Change: &Change{Actions: Actions{ActionCreate}}},
		identityImport,
	}
	plan := tfjson.Plan{ResourceChanges: resourceChanges}

	FilterNoOpResources(&plan)

	expectedResourceChangesAfterFiltering := ResourceChanges{
		&ResourceChange{Address: "create", Change: &Change{Actions: Actions{ActionCreate}}},
		identityImport,
	}
	assert.Equal(t, expectedResourceChangesAfterFiltering, plan.ResourceChanges)
}
