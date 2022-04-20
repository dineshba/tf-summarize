package terraform_state

import (
	"encoding/json"
	"fmt"
)

const ColorReset = "\033[0m"
const ColorRed = "\033[31m"
const ColorGreen = "\033[32m"
const ColorYellow = "\033[33m"

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

func (rc ResourceChange) ColorPrefixAndSuffixText() (string, string) {
	var colorPrefix, suffix string
	actions := rc.Change.Actions
	if len(actions) == 1 {
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
	} else {
		colorPrefix = ColorRed
		suffix = "(+/-)"
	}
	return colorPrefix, suffix
}

type ResourceChanges []ResourceChange

type TerraformState struct {
	ResourceChanges ResourceChanges `json:"resource_changes"`
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

func (ts *TerraformState) FilterNoOpResources() {
	acc := make(ResourceChanges, 0)
	for _, r := range ts.ResourceChanges {
		if len(r.Change.Actions) == 1 && r.Change.Actions[0] == "no-op" {
			continue
		}
		acc = append(acc, r)
	}
	ts.ResourceChanges = acc
}

func (ts *TerraformState) AllChanges() map[string]ResourceChanges {
	addedResources := addedResources(ts.ResourceChanges)
	deletedResources := deletedResources(ts.ResourceChanges)
	updatedResources := updatedResources(ts.ResourceChanges)
	recreatedResources := recreatedResources(ts.ResourceChanges)

	return map[string]ResourceChanges{
		"add":      addedResources,
		"delete":   deletedResources,
		"update":   updatedResources,
		"recreate": recreatedResources,
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
