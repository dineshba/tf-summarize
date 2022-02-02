package main

import (
	"fmt"
	"strings"
)

type Tree struct {
	name     string
	level    int
	value    *ResourceChange
	children Trees
}

func (t Tree) String() string {
	return fmt.Sprintf("{name: %s, children: %+v}", t.name, t.children)
}

type Trees []*Tree

func (t Trees) String() string {
	result := ""
	for _, tree := range t {
		result = fmt.Sprintf("%s,{name: %s, children: %+v}", result, tree.name, tree.children)
	}
	return strings.TrimPrefix(result, ",")
}

func CreateTree(resources ResourceChanges) Trees {
	result := &Tree{name: ".", children: Trees{}, level: 0}
	for _, r := range resources {
		levels := strings.Split(r.Address, ".")
		createTreeMultiLevel(r, levels, result)
	}
	return result.children
}

func createTreeMultiLevel(r ResourceChange, levels []string, currentTree *Tree) {
	parentTree := currentTree
	for i, name := range levels {
		matchedTree := getTree(name, parentTree.children)
		if matchedTree == nil {
			var resourceChange *ResourceChange
			if i+1 == len(levels) {
				resourceChange = &r
			}
			newTree := &Tree{
				name:  name,
				value: resourceChange,
			}
			parentTree.children = append(parentTree.children,
				newTree)
			parentTree = newTree
		} else {
			parentTree = matchedTree
		}
	}
}

func getTree(name string, siblings Trees) *Tree {
	for _, s := range siblings {
		if s.name == name {
			return s
		}
	}
	return nil
}
