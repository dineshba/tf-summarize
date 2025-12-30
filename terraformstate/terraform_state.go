package terraformstate

import (
	"encoding/json"
	"fmt"
	"sort"

	tfjson "github.com/hashicorp/terraform-json"
)

const ColorReset = "\033[0m"
const ColorRed = "\033[31m"
const ColorGreen = "\033[32m"
const ColorMagenta = "\033[35m"
const ColorYellow = "\033[33m"
const ColorCyan = "\033[36m"

type ResourceChanges = []*tfjson.ResourceChange //Type alias for brevity

func GetColorPrefixAndSuffixText(rc *tfjson.ResourceChange) (string, string) {
	var colorPrefix, suffix string
	actions := (*rc).Change.Actions
	if len(actions) == 1 && !actions.NoOp() {
		if actions.Create() {
			colorPrefix = ColorGreen
			suffix = "(+)"
		} else if actions.Delete() {
			colorPrefix = ColorRed
			suffix = "(-)"
		} else {
			colorPrefix = ColorYellow
			suffix = "(~)"
		}
	} else if rc.Change.Importing != nil && rc.Change.Importing.ID != "" {
		colorPrefix = ColorCyan
		suffix = "(i)"
	} else if actions.DestroyBeforeCreate() {
		colorPrefix = ColorMagenta
		suffix = "(-/+)"
	} else if actions.CreateBeforeDestroy() {
		colorPrefix = ColorMagenta
		suffix = "(+/-)"
	}
	return colorPrefix, suffix
}

func Parse(input []byte) (tfjson.Plan, error) {
	plan := tfjson.Plan{}
	err := json.Unmarshal(input, &plan)
	if err != nil {
		return tfjson.Plan{}, fmt.Errorf("error when parsing input: %s", err.Error())
	}
	return plan, nil
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
		if r.Change.Importing == nil {
			continue
		}

		id := r.Change.Importing.ID
		identity := r.Change.Importing.Identity
		if id != "" || identity != nil {
			acc = append(acc, r)
		}
	}
	return acc
}

func FilterNoOpResources(ts *tfjson.Plan) {
	acc := make(ResourceChanges, 0)
	for _, r := range ts.ResourceChanges {
		// ID-based importing
		oldImporting := r.Change.Importing != nil && r.Change.Importing.ID != ""

		// New identity-based importing introduced in terraform 1.12
		newImporting := r.Change.Importing != nil && r.Change.Importing.Identity != nil

		if r.Change.Actions.NoOp() && !oldImporting && !newImporting {
			continue
		}
		acc = append(acc, r)
	}
	ts.ResourceChanges = acc
}

func GetAllResourceChanges(plan tfjson.Plan) map[string]ResourceChanges {
	addedResources := addedResources(plan.ResourceChanges)
	deletedResources := deletedResources(plan.ResourceChanges)
	updatedResources := updatedResources(plan.ResourceChanges)
	recreatedResources := recreatedResources(plan.ResourceChanges)
	importedResources := importedResources(plan.ResourceChanges)

	sortResources := func(resources ResourceChanges) {
		sort.Slice(resources, func(i, j int) bool {
			return resources[i].Address < resources[j].Address
		})
	}

	sortResources(addedResources)
	sortResources(deletedResources)
	sortResources(updatedResources)
	sortResources(recreatedResources)
	sortResources(importedResources)

	return map[string]ResourceChanges{
		"import":   importedResources,
		"add":      addedResources,
		"delete":   deletedResources,
		"update":   updatedResources,
		"recreate": recreatedResources,
	}
}

func GetAllOutputChanges(plan tfjson.Plan) map[string][]string {
	// create, update, and delete are the only available actions for outputChanges
	// https://developer.hashicorp.com/terraform/internals/json-format
	addedResources := filterOutputs(plan.OutputChanges, "create")
	deletedResources := filterOutputs(plan.OutputChanges, "delete")
	updatedResources := filterOutputs(plan.OutputChanges, "update")

	sort.Strings(addedResources)
	sort.Strings(deletedResources)
	sort.Strings(updatedResources)

	return map[string][]string{
		"add":    addedResources,
		"delete": deletedResources,
		"update": updatedResources,
	}
}

func filterResources(resources ResourceChanges, action tfjson.Action) ResourceChanges {
	acc := make(ResourceChanges, 0)
	for _, r := range resources {
		if len(r.Change.Actions) == 1 && r.Change.Actions[0] == action {
			acc = append(acc, r)
		}
	}
	return acc
}

func filterOutputs(outputChanges map[string]*tfjson.Change, action tfjson.Action) []string {
	acc := make([]string, 0)
	for k, v := range outputChanges {
		if len(v.Actions) == 1 && v.Actions[0] == action {
			acc = append(acc, k)
		}
	}
	return acc
}
