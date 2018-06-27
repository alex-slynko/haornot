package main_test

import (
	"path/filepath"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

func TestHaornot(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Haornot Suite")
}

var cwd, pathToCLI string

var _ = BeforeSuite(func() {
	var err error
	pathToCLI, err = gexec.Build("github.com/alex-slynko/haornot")
	Expect(err).ToNot(HaveOccurred())
	_, filename, _, _ := runtime.Caller(0)
	cwd = filepath.Dir(filename)
})

var _ = AfterSuite(func() {
	gexec.CleanupBuildArtifacts()
})
