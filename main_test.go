package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const (
	testVersion    string = "test"
	testExecutable string = "tf-summarize-test"
)

func TestMain(m *testing.M) {
	// compile a 'tf-summarize' for use in running tests
	exe := exec.Command("go", "build", "-ldflags", fmt.Sprintf("-X main.version=%s", testVersion), "-o", testExecutable)
	err := exe.Run()
	if err != nil {
		os.Exit(1)
	}

	m.Run()

	// delete the compiled tf-summarize
	err = os.Remove(testExecutable)
	if err != nil {
		log.Fatal(err)
	}
}

func TestVersionArg(t *testing.T) {
	args := []string{
		"-v",
	}

	for _, arg := range args {
		t.Run(fmt.Sprintf("when tf-summarize is passed '%s'", arg), func(t *testing.T) {
			output, err := exec.Command(fmt.Sprintf("./%s", testExecutable), arg).CombinedOutput()
			if err != nil {
				t.Errorf("expected '%s' not to cause error; got '%v'", arg, err)
			}

			if !strings.Contains(string(output), testVersion) {
				t.Errorf("expected '%s' to output version '%s'; got '%s'", arg, testVersion, output)
			}
		})
	}
}

func TestTFSummarize(t *testing.T) {
	tests := []struct {
		command        string
		expectedError  error
		expectedOutput string
	}{{
		command:        fmt.Sprintf("./%s -md example/tfplan.json", testExecutable),
		expectedOutput: "basic.txt",
	}, {
		command:        fmt.Sprintf("cat example/tfplan.json | ./%s -md", testExecutable),
		expectedOutput: "basic.txt",
	}}

	for _, test := range tests {
		t.Run(fmt.Sprintf("when tf-summarize is passed '%q'", test.command), func(t *testing.T) {
			output, err := exec.Command("/bin/sh", "-c", test.command).CombinedOutput()
			if err != nil && test.expectedError == nil {
				t.Errorf("expected '%s' not to error; got '%v'", test.command, err)
			}

			b, err := os.ReadFile(fmt.Sprintf("testdata/%s", test.expectedOutput))
			if err != nil {
				t.Errorf("error reading file '%s': '%v'", test.expectedOutput, err)
			}

			expected := string(b)

			if test.expectedError != nil && err == nil {
				t.Errorf("expected error '%s'; got '%v'", test.expectedError.Error(), err)
			}

			if test.expectedError != nil && err != nil && test.expectedError.Error() != err.Error() {
				t.Errorf("expected error '%s'; got '%v'", test.expectedError.Error(), err.Error())
			}

			if string(output) != expected {
				t.Logf("expected output: \n%s", expected)
				t.Logf("got output: \n%s", output)
				t.Errorf("received unexpected output from '%s'", test.command)
			}
		})
	}
}
