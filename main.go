package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"terraform-plan-summary/reader"

	"github.com/olekukonko/tablewriter"
)

func main() {
	tree := flag.Bool("tree", false, "tree format")
	separateTree := flag.Bool("separate-tree", false, "separate tree format")
	flag.Parse()

	newReader, err := createReader(os.Stdin, os.Args)
	if err != nil {
		panic(fmt.Errorf("error creating input reader: %s", err.Error()))
	}

	input, err := newReader.Read()
	if err != nil {
		panic(fmt.Errorf("error reading from input: %s", err.Error()))
	}

	ts := terraformState{}
	err = json.Unmarshal(input, &ts)
	if err != nil {
		panic(fmt.Sprintf("Error when parsing input: %s", err))
	}

	ts.filterNoOpResources()
	allChanges := ts.AllChanges()

	trees := CreateTree(ts.ResourceChanges)
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
		trees := CreateTree(ts.ResourceChanges)

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

func createReader(stdin *os.File, args []string) (reader.Reader, error) {
	stat, _ := stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		return reader.NewStdinReader(), nil
	}
	if len(args) < 2 {
		return nil, fmt.Errorf("should either have stdin through pipe or first argument should be file")
	}
	fileName := os.Args[1]
	return reader.NewFileReader(fileName), nil
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
