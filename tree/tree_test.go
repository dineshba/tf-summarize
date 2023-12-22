package tree

import (
	"testing"

	"github.com/dineshba/tf-summarize/terraformstate"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
)

func TestCreateTreeForEmptyResourceChanges(t *testing.T) {
	assert.Equal(t, Trees{}, CreateTree(terraformstate.ResourceChanges{}))
}

func TestCreateTreeForOneResourceChanges(t *testing.T) {
	resourceChanges := terraformstate.ResourceChanges{
		&tfjson.ResourceChange{Address: "a"},
	}
	expected := Trees{
		{
			Name:  "a",
			Value: &tfjson.ResourceChange{Address: "a"},
		},
	}
	assert.Equal(t, expected, CreateTree(resourceChanges))
}

func TestCreateTreeWithQuotesInResources(t *testing.T) {
	resourceChanges := terraformstate.ResourceChanges{
		&tfjson.ResourceChange{Address: "a.b[\"c.d\"]"},
	}
	expected := Trees{
		{
			Name:  "a",
			Value: nil,
			Children: Trees{
				{
					Name:  "b[\"c.d\"]",
					Value: &tfjson.ResourceChange{Address: "a.b[\"c.d\"]"},
				},
			},
		},
	}
	actual := CreateTree(resourceChanges)
	assert.Equal(t, expected, actual)
}

func TestCreateTreeForOneResourceChangesMultiLevel(t *testing.T) {
	resourceChanges := terraformstate.ResourceChanges{
		&tfjson.ResourceChange{Address: "a.b.c"},
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
							Value: &tfjson.ResourceChange{Address: "a.b.c"},
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
		&tfjson.ResourceChange{Address: "a"},
		&tfjson.ResourceChange{Address: "b"},
	}
	expected := Trees{
		{
			Name:  "a",
			Value: &tfjson.ResourceChange{Address: "a"},
		},
		{
			Name:  "b",
			Value: &tfjson.ResourceChange{Address: "b"},
		},
	}

	assert.Equal(t, expected, CreateTree(resourceChanges))
}

func TestCreateTreeForTwoResourceChangesOverlap(t *testing.T) {
	resourceChanges := terraformstate.ResourceChanges{
		&tfjson.ResourceChange{Address: "a.b"},
		&tfjson.ResourceChange{Address: "a.c.x"},
		&tfjson.ResourceChange{Address: "a.c.y"},
		&tfjson.ResourceChange{Address: "d"},
	}
	expected := Trees{
		{
			Name: "a",
			Children: Trees{
				{
					Name:     "b",
					Value:    &tfjson.ResourceChange{Address: "a.b"},
					Children: nil,
				},
				{
					Name: "c",
					Children: Trees{
						{
							Name:     "x",
							Value:    &tfjson.ResourceChange{Address: "a.c.x"},
							Children: nil,
						},
						{
							Name:     "y",
							Value:    &tfjson.ResourceChange{Address: "a.c.y"},
							Children: nil,
						},
					},
				},
			},
		},
		{
			Name:  "d",
			Value: &tfjson.ResourceChange{Address: "d"},
		},
	}

	tree := CreateTree(resourceChanges)
	s := tree.String()
	assert.Equal(t, expected.String(), s)
}
