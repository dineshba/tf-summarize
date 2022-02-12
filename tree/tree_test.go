package tree

import (
	"terraform-plan-summary/terraform_state"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTreeForEmptyResourceChanges(t *testing.T) {
	assert.Equal(t, Trees{}, CreateTree(terraform_state.ResourceChanges{}))
}

func TestCreateTreeForOneResourceChanges(t *testing.T) {
	resourceChanges := terraform_state.ResourceChanges{
		terraform_state.ResourceChange{Address: "a"},
	}
	expected := Trees{
		{
			Name:  "a",
			Value: &terraform_state.ResourceChange{Address: "a"},
		},
	}
	assert.Equal(t, expected, CreateTree(resourceChanges))
}

func TestCreateTreeForOneResourceChangesMultiLevel(t *testing.T) {
	resourceChanges := terraform_state.ResourceChanges{
		terraform_state.ResourceChange{Address: "a.b.c"},
	}
	expected := Trees{
		{
			Name:  "a",
			Value: nil,
			Children: Trees{
				{
					Name:  "b",
					Value: nil,
					Children: Trees{
						{
							Name:  "c",
							Value: &terraform_state.ResourceChange{Address: "a.b.c"},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expected, CreateTree(resourceChanges))
}

func TestCreateTreeForTwoResourceChangesNoOverlap(t *testing.T) {
	resourceChanges := terraform_state.ResourceChanges{
		terraform_state.ResourceChange{Address: "a"},
		terraform_state.ResourceChange{Address: "b"},
	}
	expected := Trees{
		{
			Name:  "a",
			Value: &terraform_state.ResourceChange{Address: "a"},
		},
		{
			Name:  "b",
			Value: &terraform_state.ResourceChange{Address: "b"},
		},
	}

	assert.Equal(t, expected, CreateTree(resourceChanges))
}

func TestCreateTreeForTwoResourceChangesOverlap(t *testing.T) {
	resourceChanges := terraform_state.ResourceChanges{
		terraform_state.ResourceChange{Address: "a.b"},
		terraform_state.ResourceChange{Address: "a.c.x"},
		terraform_state.ResourceChange{Address: "a.c.y"},
		terraform_state.ResourceChange{Address: "d"},
	}
	expected := Trees{
		{
			Name: "a",
			Children: Trees{
				{
					Name:     "b",
					Value:    &terraform_state.ResourceChange{Address: "a.b"},
					Children: nil,
				},
				{
					Name: "c",
					Children: Trees{
						{
							Name:     "x",
							Value:    &terraform_state.ResourceChange{Address: "a.c.x"},
							Children: nil,
						},
						{
							Name:     "y",
							Value:    &terraform_state.ResourceChange{Address: "a.c.y"},
							Children: nil,
						},
					},
				},
			},
		},
		{
			Name:  "d",
			Value: &terraform_state.ResourceChange{Address: "d"},
		},
	}

	tree := CreateTree(resourceChanges)
	s := tree.String()
	assert.Equal(t, expected.String(), s)
}
