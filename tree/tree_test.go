package tree

import (
	"testing"

	"github.com/dineshba/tf-summarize/terraformstate"
	"github.com/stretchr/testify/assert"
)

func TestCreateTreeForEmptyResourceChanges(t *testing.T) {
	assert.Equal(t, Trees{}, CreateTree(terraformstate.ResourceChanges{}))
}

func TestCreateTreeForOneResourceChanges(t *testing.T) {
	resourceChanges := terraformstate.ResourceChanges{
		terraformstate.ResourceChange{Address: "a"},
	}
	expected := Trees{
		{
			Name:  "a",
			Value: &terraformstate.ResourceChange{Address: "a"},
		},
	}
	assert.Equal(t, expected, CreateTree(resourceChanges))
}

func TestCreateTreeWithQuotesInResources(t *testing.T) {
	resourceChanges := terraformstate.ResourceChanges{
		terraformstate.ResourceChange{Address: "a.b[\"c.d\"]"},
	}
	expected := Trees{
		{
			Name:  "a",
			Value: nil,
			Children: Trees{
				{
					Name:  "b[\"c.d\"]",
					Value: &terraformstate.ResourceChange{Address: "a.b[\"c.d\"]"},
				},
			},
		},
	}
	actual := CreateTree(resourceChanges)
	assert.Equal(t, expected, actual)
}

func TestCreateTreeForOneResourceChangesMultiLevel(t *testing.T) {
	resourceChanges := terraformstate.ResourceChanges{
		terraformstate.ResourceChange{Address: "a.b.c"},
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
							Value: &terraformstate.ResourceChange{Address: "a.b.c"},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expected, CreateTree(resourceChanges))
}

func TestCreateTreeForTwoResourceChangesNoOverlap(t *testing.T) {
	resourceChanges := terraformstate.ResourceChanges{
		terraformstate.ResourceChange{Address: "a"},
		terraformstate.ResourceChange{Address: "b"},
	}
	expected := Trees{
		{
			Name:  "a",
			Value: &terraformstate.ResourceChange{Address: "a"},
		},
		{
			Name:  "b",
			Value: &terraformstate.ResourceChange{Address: "b"},
		},
	}

	assert.Equal(t, expected, CreateTree(resourceChanges))
}

func TestCreateTreeForTwoResourceChangesOverlap(t *testing.T) {
	resourceChanges := terraformstate.ResourceChanges{
		terraformstate.ResourceChange{Address: "a.b"},
		terraformstate.ResourceChange{Address: "a.c.x"},
		terraformstate.ResourceChange{Address: "a.c.y"},
		terraformstate.ResourceChange{Address: "d"},
	}
	expected := Trees{
		{
			Name: "a",
			Children: Trees{
				{
					Name:     "b",
					Value:    &terraformstate.ResourceChange{Address: "a.b"},
					Children: nil,
				},
				{
					Name: "c",
					Children: Trees{
						{
							Name:     "x",
							Value:    &terraformstate.ResourceChange{Address: "a.c.x"},
							Children: nil,
						},
						{
							Name:     "y",
							Value:    &terraformstate.ResourceChange{Address: "a.c.y"},
							Children: nil,
						},
					},
				},
			},
		},
		{
			Name:  "d",
			Value: &terraformstate.ResourceChange{Address: "d"},
		},
	}

	tree := CreateTree(resourceChanges)
	s := tree.String()
	assert.Equal(t, expected.String(), s)
}
