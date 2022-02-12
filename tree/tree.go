package tree

import (
	"fmt"
	"github.com/m1gwings/treedrawer/tree"
	"strings"
	"terraform-plan-summary/terraform_state"
)

type Tree struct {
	Name     string
	level    int
	Value    *terraform_state.ResourceChange
	Children Trees
}

func (t Tree) String() string {
	return fmt.Sprintf("{name: %s, children: %+v}", t.Name, t.Children)
}

type Trees []*Tree

func (t Trees) DrawableTree() *tree.Tree {
	newTree := tree.NewTree(tree.NodeString("."))
	for _, t1 := range t {
		t1.AddChild(newTree)
	}
	return newTree
}

func (t *Tree) AddChild(parent *tree.Tree) {
	isLeafNode := len(t.Children) == 0

	var childNode tree.NodeString
	if isLeafNode {
		_, suffix := t.Value.ColorPrefixAndSuffixText()
		childNode = tree.NodeString(fmt.Sprintf("%s%s", t.Name, suffix))
	} else {
		childNode = tree.NodeString(t.Name)
	}

	currentChildIndex := len(parent.Children())
	parent.AddChild(childNode)
	currentTree, err := parent.Child(currentChildIndex)
	for _, c := range t.Children {
		if err != nil {
			panic(err)
		}
		c.AddChild(currentTree)
	}
}

func (t Trees) String() string {
	result := ""
	for _, tree := range t {
		result = fmt.Sprintf("%s,{name: %s, children: %+v}", result, tree.Name, tree.Children)
	}
	return strings.TrimPrefix(result, ",")
}

func CreateTree(resources terraform_state.ResourceChanges) Trees {
	result := &Tree{Name: ".", Children: Trees{}, level: 0}
	for _, r := range resources {
		levels := strings.Split(r.Address, ".")
		createTreeMultiLevel(r, levels, result)
	}
	return result.Children
}

func createTreeMultiLevel(r terraform_state.ResourceChange, levels []string, currentTree *Tree) {
	parentTree := currentTree
	for i, name := range levels {
		matchedTree := getTree(name, parentTree.Children)
		if matchedTree == nil {
			var resourceChange *terraform_state.ResourceChange
			if i+1 == len(levels) {
				resourceChange = &r
			}
			newTree := &Tree{
				Name:  name,
				Value: resourceChange,
			}
			parentTree.Children = append(parentTree.Children,
				newTree)
			parentTree = newTree
		} else {
			parentTree = matchedTree
		}
	}
}

func getTree(name string, siblings Trees) *Tree {
	for _, s := range siblings {
		if s.Name == name {
			return s
		}
	}
	return nil
}
