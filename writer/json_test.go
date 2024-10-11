package writer

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/dineshba/tf-summarize/terraformstate"
	. "github.com/hashicorp/terraform-json"
	"github.com/nsf/jsondiff"
)

// Struct to hold arguments and expected output for each test case
type testCase struct {
	args           terraformstate.ResourceChanges
	expectedOutput map[string]interface{}
}

func TestJSONWriter(t *testing.T) {

	// Define test cases
	testCases := []testCase{
		{
			args: terraformstate.ResourceChanges{
				{
					Address: "module.test.azapi_resource.logical_network",
					Type:    "aws_instance",
					Name:    "example",
					Change: &Change{
						Before:  map[string]interface{}{"name": "old_instance"},
						After:   map[string]interface{}{"name": "new_instance"},
						Actions: Actions{ActionCreate},
					},
				},
			},
			expectedOutput: map[string]interface{}{
				"module": map[string]interface{}{
					"test": map[string]interface{}{
						"azapi_resource": map[string]interface{}{
							"logical_network": map[string]interface{}{
								"(+)": map[string]interface{}{
									"name": "new_instance",
								},
							},
						},
					},
				},
			},
		},
		{
			args: terraformstate.ResourceChanges{
				{
					Address: "module.test.aws_s3_bucket.example",
					Type:    "aws_s3_bucket",
					Name:    "example",
					Change: &Change{
						Before:  map[string]interface{}{"name": "old_bucket"},
						After:   map[string]interface{}{"name": "new_bucket"},
						Actions: Actions{ActionUpdate},
					},
				},
			},
			expectedOutput: map[string]interface{}{
				"module": map[string]interface{}{
					"test": map[string]interface{}{
						"aws_s3_bucket": map[string]interface{}{
							"example": map[string]interface{}{
								"(~)": map[string]interface{}{
									"name": map[string]interface{}{
										"changed": []string{
											"old_bucket",
											"new_bucket",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			args: terraformstate.ResourceChanges{
				{
					Address: "module.test.aws_security_group.example",
					Type:    "aws_security_group",
					Name:    "example",
					Change: &Change{
						Before:  map[string]interface{}{"name": "old_sg"},
						After:   map[string]interface{}{"name": "new_sg"},
						Actions: Actions{ActionDelete},
					},
				},
			},
			expectedOutput: map[string]interface{}{
				"module": map[string]interface{}{
					"test": map[string]interface{}{
						"aws_security_group": map[string]interface{}{
							"example": map[string]interface{}{
								"(-)": map[string]interface{}{
									"name": "old_sg",
								},
							},
						},
					},
				},
			},
		},
	}

	// Iterate through test cases
	for i, tc := range testCases {
		jsonWriter := NewJSONWriter(tc.args)
		var buf bytes.Buffer
		err := jsonWriter.Write(&buf)
		if err != nil {
			t.Fatalf("Unexpected error in test case %d: %v", i+1, err)
		}
		expectedJSON, err := json.Marshal(tc.expectedOutput)
		if err != nil {
			t.Fatalf("Error marshalling expected output in test case %d: %v", i+1, err)
		}
		opts := jsondiff.DefaultJSONOptions()
		diff, str := jsondiff.Compare(expectedJSON, buf.Bytes(), &opts)
		if diff != jsondiff.FullMatch {
			t.Errorf("Output mismatch in test case %d: %s", i+1, str)
		}
	}
}
