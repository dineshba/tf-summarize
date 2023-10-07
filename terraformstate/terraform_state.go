package terraformstate

import (
	"encoding/json"
	"fmt"
)

const ColorReset = "\033[0m"
const ColorRed = "\033[31m"
const ColorGreen = "\033[32m"
const ColorMagenta = "\033[35m"
const ColorYellow = "\033[33m"
const ColorCyan = "\033[36m"

type ResourceChange struct {
	Address       string `json:"address"`
	ModuleAddress string `json:"module_address"`
	Mode          string `json:"mode"`
	Type          string `json:"type"`
	Name          string `json:"name"`
	ProviderName  string `json:"provider_name"`
	Change        struct {
		Actions   []string        `json:"actions"`
		Before    json.RawMessage `json:"before,omitempty"`
		After     json.RawMessage `json:"after,omitempty"`
		Importing struct {
			ID string `json:"id"`
		} `json:"importing"`
	} `json:"change"`
	ActionReason string `json:"action_reason,omitempty"`
}

type OutputValues struct {
	Actions []string        `json:"actions"`
	Before  json.RawMessage `json:"before"`
	After   json.RawMessage `json:"after"`
}

func (rc ResourceChange) ColorPrefixAndSuffixText() (string, string) {
	var colorPrefix, suffix string
	actions := rc.Change.Actions
	if len(actions) == 1 && actions[0] != "no-op" {
		if actions[0] == "create" {
			colorPrefix = ColorGreen
			suffix = "(+)"
		} else if actions[0] == "delete" {
			colorPrefix = ColorRed
			suffix = "(-)"
		} else {
			colorPrefix = ColorYellow
			suffix = "(~)"
		}
	} else if rc.Change.Importing.ID != "" {
		colorPrefix = ColorCyan
		suffix = "(i)"
	} else {
		colorPrefix = ColorMagenta
		suffix = "(+/-)"
	}
	return colorPrefix, suffix
}

type ResourceChanges []ResourceChange

type TerraformState struct {
	ResourceChanges ResourceChanges         `json:"resource_changes"`
	OutputChanges   map[string]OutputValues `json:"output_changes"`
}

func Parse(input []byte) (TerraformState, error) {
	ts := TerraformState{}
	err := json.Unmarshal(input, &ts)
	if err != nil {
		return TerraformState{}, fmt.Errorf("error when parsing input: %s", err.Error())
	}
	return ts, nil
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

func importedResources(resources ResourceChanges) ResourceChanges {
	acc := make(ResourceChanges, 0)
	for _, r := range resources {
		id := r.Change.Importing.ID
		if id != "" {
			acc = append(acc, r)
		}
	}
	return acc
}

func (ts *TerraformState) FilterNoOpResources() {
	acc := make(ResourceChanges, 0)
	for _, r := range ts.ResourceChanges {
		if len(r.Change.Actions) == 1 && r.Change.Actions[0] == "no-op" && r.Change.Importing.ID == "" {
			continue
		}
		acc = append(acc, r)
	}
	ts.ResourceChanges = acc
}

func (ts *TerraformState) AllResourceChanges() map[string]ResourceChanges {
	addedResources := addedResources(ts.ResourceChanges)
	deletedResources := deletedResources(ts.ResourceChanges)
	updatedResources := updatedResources(ts.ResourceChanges)
	recreatedResources := recreatedResources(ts.ResourceChanges)
	importedResources := importedResources(ts.ResourceChanges)

	return map[string]ResourceChanges{
		"import":   importedResources,
		"add":      addedResources,
		"delete":   deletedResources,
		"update":   updatedResources,
		"recreate": recreatedResources,
	}
}

func (ts *TerraformState) AllOutputChanges() map[string][]string {
	// create, update, and delete are the only available actions for outputChanges
	// https://developer.hashicorp.com/terraform/internals/json-format
	addedResources := filterOutputs(ts.OutputChanges, "create")
	deletedResources := filterOutputs(ts.OutputChanges, "delete")
	updatedResources := filterOutputs(ts.OutputChanges, "update")

	return map[string][]string{
		"add":    addedResources,
		"delete": deletedResources,
		"update": updatedResources,
	}
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

func filterOutputs(outputChanges map[string]OutputValues, action string) []string {
	acc := make([]string, 0)
	for k, v := range outputChanges {
		if len(v.Actions) == 1 && v.Actions[0] == action {
			acc = append(acc, k)
		}
	}
	return acc
}
