package main

import (
	"io/ioutil"
	"os"

	"github.com/alex-slynko/haornot/analyzer"
	"github.com/alex-slynko/haornot/formatter"
)

func main() {
	if len(os.Args) < 2 {
		failWith("Spec file is required")
	}
	contents, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		failWith(err.Error())
	}
	output, err := analyzer.Analyze(contents)

	if err != nil {
		failWith(err.Error())
	}
	if len(output) > 0 {
		failWith(prettify(output))
	}
}

func failWith(message string) {
	formatter := formatter.ImageFormatter{}
	formatter.Fail(message)
	os.Exit(1)
}

func prettify(errors []string) string {
	result := ""

	for _, msg := range errors {
		result = result + "😿 " + msg + "\n"
	}
	return result
}
