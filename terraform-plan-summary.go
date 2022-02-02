package main

import (
	"bufio"
	"encoding/json"
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

	// check if there is somethinig to read on STDIN
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

	tableString := make([][]string, 0, 4)
	for change, changedresources := range allChanges {
		for _, changedresource := range changedresources {
			tableString = append(tableString, []string{change, changedresource.Address})
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Change", "Name"})
	table.SetAutoMergeCells(true)
	table.SetRowLine(true)
	table.AppendBulk(tableString)
	table.Render()

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
