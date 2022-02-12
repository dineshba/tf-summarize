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
	drawable := flag.Bool("draw", false, "drawable tree format")
	flag.Parse()

	args := flag.Args()
	err := validateFlags(*tree, *separateTree, *drawable, args)
	logIfErrorAndExit("invalid input flags: %s", err)

	newReader, err := reader.CreateReader(os.Stdin, args)
	logIfErrorAndExit("error creating input reader: %s", err)

	input, err := newReader.Read()
	logIfErrorAndExit("error reading from input: %s", err)

	terraformState, err := terraform_state.Parse(input)
	logIfErrorAndExit("%s", err)

	terraformState.FilterNoOpResources()

	newWriter := writer.CreateWriter(*tree, *separateTree, *drawable, terraformState)
	err = newWriter.Write(os.Stdout)
	logIfErrorAndExit("error writing: %s", err)
}

func logIfErrorAndExit(format string, err error) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", format), err.Error())
		os.Exit(1)
	}
}

func validateFlags(tree, separateTree, drawable bool, args []string) error {
	if tree && separateTree {
		return fmt.Errorf("both -tree and -seperate-tree should not be provided")
	}
	if !tree && !separateTree && drawable {
		return fmt.Errorf("drawable should be provided with -tree or -seperate-tree")
	}
	if len(args) > 1 {
		return fmt.Errorf("only one argument is allowed which is filename, but got %v", args)
	}
	return nil
}
