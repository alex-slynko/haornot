package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/alex-slynko/haornot/analyzer"
	"github.com/alex-slynko/haornot/formatter"
	"github.com/alex-slynko/haornot/types"
)

var hasErrors bool

func main() {
	if len(os.Args) < 2 {
		failWith("Spec file is required")
		os.Exit(1)
	}
	contents, err := ioutil.ReadFile(os.Args[1])

	if err != nil {
		failWith(err.Error())
		os.Exit(1)
	}
	hasErrors = false
	manifests := bytes.Split(contents, []byte("\n---"))
	totalDeployments := 0
	deploymentsWithoutErrors := 0
	for _, manifest := range manifests {
		output, err := analyzer.Analyze(manifest)

		if err == analyzer.ErrNotADeployment {
			continue
		}

		totalDeployments++
		if err != nil {
			showError(err)
			continue
		}
		showDeploymentMessage(output)
		if len(output.Errors) == 0 {
			deploymentsWithoutErrors++
		}
	}

	if totalDeployments == 0 {
		failWith("only deployments can be analyzed")
	}
	if !hasErrors {
		success()
	} else {
		os.Exit(1)
	}
}

func success() {
	formatter := formatter.ImageFormatter{}
	formatter.Success()
	fmt.Println()
	fmt.Println("😸 Your spec file satifies all checks 😸")
}

func failWith(message string) {
	formatter := formatter.ImageFormatter{}
	formatter.CriticalFail(message)
	os.Exit(1)
}

func showError(err error) {
	formatter := formatter.ImageFormatter{}
	formatter.CriticalFail(err.Error())
}

func showDeploymentMessage(em *types.Message) {
	formatter := formatter.ImageFormatter{}
	if len(em.Errors) > 0 {
		formatter.Fail(em)
		hasErrors = true
	} else {
		formatter.Progress(em)
	}
}
