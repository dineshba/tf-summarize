package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTreeForEmptyResourceChanges(t *testing.T) {
	assert.Equal(t, Trees{}, CreateTree(ResourceChanges{}))
}

func TestCreateTreeForOneResourceChanges(t *testing.T) {
	resourceChanges := ResourceChanges{
		ResourceChange{Address: "a"},
	}
	expected := Trees{
		{
			name:  "a",
			value: &ResourceChange{Address: "a"},
		},
	}
	assert.Equal(t, expected, CreateTree(resourceChanges))
}

func TestCreateTreeForOneResourceChangesMultiLevel(t *testing.T) {
	resourceChanges := ResourceChanges{
		ResourceChange{Address: "a.b.c"},
	}
	expected := Trees{
		{
			name:  "a",
			value: nil,
			children: Trees{
				{
					name:  "b",
					value: nil,
					children: Trees{
						{
							name:  "c",
							value: &ResourceChange{Address: "a.b.c"},
						},
					},
				},
			},
		},
	}

	assert.Equal(t, expected, CreateTree(resourceChanges))
}

func TestCreateTreeForTwoResourceChangesNoOverlap(t *testing.T) {
	resourceChanges := ResourceChanges{
		ResourceChange{Address: "a"},
		ResourceChange{Address: "b"},
	}
	expected := Trees{
		{
			name:  "a",
			value: &ResourceChange{Address: "a"},
		},
		{
			name:  "b",
			value: &ResourceChange{Address: "b"},
		},
	}

	assert.Equal(t, expected, CreateTree(resourceChanges))
}

func TestCreateTreeForTwoResourceChangesOverlap(t *testing.T) {
	resourceChanges := ResourceChanges{
		ResourceChange{Address: "a.b"},
		ResourceChange{Address: "a.c.x"},
		ResourceChange{Address: "a.c.y"},
		ResourceChange{Address: "d"},
	}
	expected := Trees{
		{
			name:  "a",
			children: Trees{
				{
					name:     "b",
					value:    &ResourceChange{Address: "a.b"},
					children: nil,
				},
				{
					name:     "c",
					children: Trees{
						{
							name:     "x",
							value:    &ResourceChange{Address: "a.c.x"},
							children: nil,
						},
						{
							name:     "y",
							value:    &ResourceChange{Address: "a.c.y"},
							children: nil,
						},
					},
				},
			},
		},
		{
			name:  "d",
			value: &ResourceChange{Address: "d"},
		},
	}

	tree := CreateTree(resourceChanges)
	s := tree.String()
	assert.Equal(t, expected.String(), s)
}
