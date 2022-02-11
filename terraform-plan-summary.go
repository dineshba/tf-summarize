package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/olekukonko/tablewriter"
)

type ResourceChange struct {
	Address       string `json:"address"`
	ModuleAddress string `json:"module_address"`
	Mode          string `json:"mode"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	ProviderName  string `json:"provider_name"`
	Change        struct {
		Actions []string `json:"actions"`
	} `json:"change"`
	ActionReason string `json:"action_reason,omitempty"`
}

type ResourceChanges []ResourceChange

type terraformState struct {
	ResourceChanges ResourceChanges `json:"resource_changes"`
}

func main() {

	var input []byte

	tree := flag.Bool("tree", false, "tree format")
	separateTree := flag.Bool("separate-tree", false, "separate tree format")
	flag.Parse()

	// check if there is something to read on STDIN
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		f := os.Stdin
		r := bufio.NewReader(f)
		line, err := r.ReadBytes('\n')
		for err == nil {
			input = append(input, line...)
			line, err = r.ReadBytes('\n')
		}
		if err != io.EOF {
			fmt.Println(err)
			return
		}
	} else {
		if len(os.Args) < 2 {
			panic("Should either have stdin through pipe or first argument should be file")
		}

		fileName := os.Args[1]

		data, err := os.ReadFile(fileName)
		if err != nil {
			panic(fmt.Sprintf("Error when reading from file %s: %s", fileName, err))
		}
		input = data
	}

	ts := terraformState{}
	err := json.Unmarshal(input, &ts)
	if err != nil {
		panic(fmt.Sprintf("Error when parsing input: %s", err))
	}

	resources := filterNoOpResources(ts.ResourceChanges)

	addedResources := addedResources(resources)
	deletedResources := deletedResources(resources)
	updatedResources := updatedResources(resources)
	recreatedResources := recreatedResources(resources)

	allChanges := map[string]ResourceChanges{
		"added":     addedResources,
		"deleted":   deletedResources,
		"updated":   updatedResources,
		"recreated": recreatedResources,
	}

	trees := CreateTree(resources)
	drawableTree := trees.DrawableTree()
	fmt.Println(drawableTree)

	if *separateTree {
		for k, v := range allChanges {
			trees := CreateTree(v)
			if len(v) > 0 {
				fmt.Println(k)
				for _, tree := range trees {
					printTree(tree, "")
				}
				fmt.Println("------------------------------")
			}
		}
		return
	}

	if *tree {
		trees := CreateTree(resources)

		for _, tree := range trees {
			printTree(tree, "")
		}
		return
	}
	tableString := make([][]string, 0, 4)
	for change, changedResources := range allChanges {
		for _, changedResource := range changedResources {
			tableString = append(tableString, []string{change, changedResource.Address})
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Change", "Name"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.AppendBulk(tableString)
	table.Render()

}

var colorReset = "\033[0m"
var colorRed = "\033[31m"
var colorGreen = "\033[32m"
var colorYellow = "\033[33m"

func printTree(tree *Tree, prefix string) {
	if tree.value != nil {
		actions := tree.value.Change.Actions
		colorPrefix := ""
		suffix := ""
		if len(actions) == 1 {
			if actions[0] == "create" {
				colorPrefix = colorGreen
				suffix = "(+)"
			} else if actions[0] == "delete" {
				colorPrefix = colorRed
				suffix = "(-)"
			} else {
				colorPrefix = colorYellow
				suffix = "(~)"
			}
		}
		if len(actions) == 2 {
			colorPrefix = colorRed
			suffix = "(+/-)"
		}
		fmt.Printf("%s%s%s%s%s\n", prefix, colorPrefix, tree.name, suffix, colorReset)
	} else {
		fmt.Printf("%s%s\n", prefix, tree.name)
	}

	for _, c := range tree.children {
		printTree(c, fmt.Sprintf("%s----", prefix))
	}
}

func addedResources(resources ResourceChanges) ResourceChanges {
	return filterResources(resources, "create")
}

func updatedResources(resources ResourceChanges) ResourceChanges {
	return filterResources(resources, "update")
}

func recreatedResources(resources ResourceChanges) ResourceChanges {
	acc := make(ResourceChanges, 0)
	for _, r := range resources {
		if len(r.Change.Actions) == 2 { // if Change is two, it will be create, delete
			acc = append(acc, r)
		}
	}
	return acc
}

func deletedResources(resources ResourceChanges) ResourceChanges {
	return filterResources(resources, "delete")
}

func filterNoOpResources(resources ResourceChanges) ResourceChanges {
	acc := make(ResourceChanges, 0)
	for _, r := range resources {
		if len(r.Change.Actions) == 1 && r.Change.Actions[0] != "no-op" {
			acc = append(acc, r)
		}
	}
	return acc
}

func filterResources(resources ResourceChanges, action string) ResourceChanges {
	acc := make(ResourceChanges, 0)
	for _, r := range resources {
		if len(r.Change.Actions) == 1 && r.Change.Actions[0] == action {
			acc = append(acc, r)
		}
	}
	return acc
}
