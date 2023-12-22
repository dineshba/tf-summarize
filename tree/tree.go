package tree

import (
	"fmt"
	"strings"

	"github.com/dineshba/tf-summarize/terraformstate"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/m1gwings/treedrawer/tree"
)

type Tree struct {
	Name     string
	level    int
	Value    *tfjson.ResourceChange
	Children Trees
}

func (t Tree) String() string {
	return fmt.Sprintf("{name: %s, children: %+v}", t.Name, t.Children)
}

func (t Tree) IsAddition() bool {
	return t.Value.Change.Actions[0] == "create"
}

func (t Tree) IsRemoval() bool {
	return t.Value.Change.Actions[0] == "delete"
}

func (t Tree) IsUpdate() bool {
	return t.Value.Change.Actions[0] == "update"
}

func (t Tree) IsRecreate() bool {
	return len(t.Value.Change.Actions) == 2
}

func (t Tree) IsImport() bool {
	return t.Value.Change.Importing.ID != ""
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
		_, suffix := terraformstate.GetColorPrefixAndSuffixText(t.Value)
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

func CreateTree(changes terraformstate.ResourceChanges) Trees {
	result := &Tree{Name: ".", Children: Trees{}, level: 0}
	for _, r := range changes {
		change := *r
		levels := splitResources(change.Address)
		createTreeMultiLevel(change, levels, result)
	}
	return result.Children
}

func splitResources(address string) []string {
	acc := make([]string, 0)
	var resource strings.Builder
	for i := 0; i < len(address); i++ {
		currentIndex := string(address[i])

		if currentIndex == "[" {
			lastIndex := strings.Index(address[i:], "]")
			resource.WriteString(address[i : i+lastIndex+1])
			i = i + lastIndex
			continue
		}

		if currentIndex == "." {
			acc = append(acc, resource.String())
			resource = strings.Builder{}
			continue
		}
		resource.Write([]byte{address[i]})
	}
	acc = append(acc, resource.String())
	return acc
}

func createTreeMultiLevel(change tfjson.ResourceChange, levels []string, currentTree *Tree) {
	parentTree := currentTree
	for i, name := range levels {
		matchedTree := getTree(name, parentTree.Children)
		if matchedTree == nil {
			var resourceChange *tfjson.ResourceChange
			if i+1 == len(levels) {
				resourceChange = &change
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
