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

func main() {
	tree := flag.Bool("tree", false, "tree format")
	separateTree := flag.Bool("separate-tree", false, "separate tree format")
	flag.Parse()

	var input []byte
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
