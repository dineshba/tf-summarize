package main

import (
	"flag"
	"fmt"
	"os"
	"terraform-plan-summary/reader"
	"terraform-plan-summary/terraform_state"
	"terraform-plan-summary/writer"
)

func main() {
	tree := flag.Bool("tree", false, "tree format")
	separateTree := flag.Bool("separate-tree", false, "separate tree format")
	drawable := flag.Bool("drawable", false, "drawable tree format")
	flag.Parse()
	err := validateFlags(*tree, *separateTree, *drawable)
	if err != nil {
		panic(fmt.Errorf("invalid input flags: %s", err.Error()))
	}

	newReader, err := reader.CreateReader(os.Stdin, os.Args)
	if err != nil {
		panic(fmt.Errorf("error creating input reader: %s", err.Error()))
	}

	input, err := newReader.Read()
	if err != nil {
		panic(fmt.Errorf("error reading from input: %s", err.Error()))
	}

	terraformState, err := terraform_state.Parse(input)
	if err != nil {
		panic(fmt.Errorf("%s", err.Error()))
	}

	terraformState.FilterNoOpResources()

	newWriter := writer.CreateWriter(*tree, *separateTree, *drawable, terraformState)
	err = newWriter.Write(os.Stdout)
	if err != nil {
		panic(fmt.Errorf("error writing: %s", err.Error()))
	}
}

func validateFlags(tree, separateTree, drawable bool) error {
	if tree && separateTree {
		return fmt.Errorf("both -tree and -seperate-tree should not be provided")
	}
	if !tree && !separateTree && drawable {
		return fmt.Errorf("drawable should be provided with -tree or -seperate-tree")
	}
	return nil
}
