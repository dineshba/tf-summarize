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

	newReader, err := createReader(os.Stdin, os.Args)
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

	writer := createWriter(*tree, *separateTree, *drawable, terraformState)
	err = writer.Write(os.Stdout)
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

func createWriter(tree, separateTree, drawable bool, terraformState terraform_state.TerraformState) writer.Writer {
	if tree {
		return writer.NewTreeWriter(terraformState.ResourceChanges, drawable)
	}
	if separateTree {
		return writer.NewSeparateTree(terraformState.AllChanges(), drawable)
	}
	return writer.NewTableWriter(terraformState.AllChanges())

	//if separateTree {
	//	for k, v := range allChanges {
	//		trees := tree2.CreateTree(v)
	//		if len(v) > 0 {
	//			fmt.Println(k)
	//
	//			if drawable {
	//				drawableTree := trees.DrawableTree()
	//				fmt.Println(drawableTree)
	//				continue
	//			}
	//
	//			for _, tree := range trees {
	//				printTree(tree, "")
	//			}
	//			fmt.Println("------------------------------")
	//		}
	//	}
	//	return
	//}
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
