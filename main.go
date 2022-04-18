package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"terraform-plan-summary/reader"
	"terraform-plan-summary/terraform_state"
	"terraform-plan-summary/writer"
)

func main() {
	tree := flag.Bool("tree", false, "[Optional] print changes in tree format")
	separateTree := flag.Bool("separate-tree", false, "[Optional] print changes in tree format for add/delete/change/recreate changes")
	drawable := flag.Bool("draw", false, "[Optional, used only with -tree or -separate-tree] draw trees instead of plain tree")
	outputFileName := flag.String("out", "", "[Optional] write output to file")

	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "\nUsage of %s [args] [tf-plan.json]\n\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	args := flag.Args()
	err := validateFlags(*tree, *separateTree, *drawable, args)
	logIfErrorAndExit("invalid input flags: %s", err, flag.Usage)

	newReader, err := reader.CreateReader(os.Stdin, args)
	logIfErrorAndExit("error creating input reader: %s", err, flag.Usage)

	input, err := newReader.Read()
	logIfErrorAndExit("error reading from input: %s", err, func() {})

	terraformState, err := terraform_state.Parse(input)
	logIfErrorAndExit("%s", err, func() {})

	terraformState.FilterNoOpResources()

	newWriter := writer.CreateWriter(*tree, *separateTree, *drawable, terraformState)

	var outputFile io.Writer = os.Stdout

	if *outputFileName != "" {
		file, err := os.OpenFile(*outputFileName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0600)
		logIfErrorAndExit("error opening file: %v", err, func() {})
		defer func() {
			if err := file.Close(); err != nil {
				logIfErrorAndExit("Error closing file: %s\n", err, func() {})
			}
		}()
		outputFile = file
	}

	err = newWriter.Write(outputFile)
	logIfErrorAndExit("error writing: %s", err, func() {})

	if err == nil && *outputFileName != "" {
		_, _ = fmt.Fprintf(os.Stderr, "Written plan summary to %s\n", *outputFileName)
	}
}

func logIfErrorAndExit(format string, err error, callback func()) {
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, fmt.Sprintf("%s\n", format), err.Error())
		callback()
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
