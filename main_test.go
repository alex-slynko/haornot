package main_test

import (
	"os/exec"
	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Haornot", func() {
	var (
		spec    string
		session *gexec.Session
	)

	BeforeEach(func() {
		spec = path.Join(cwd, "fixtures", "nginx.yml")
	})

	JustBeforeEach(func() {
		var err error
		command := exec.Command(pathToCLI, spec)
		session, err = gexec.Start(command, GinkgoWriter, GinkgoWriter)
		Expect(err).NotTo(HaveOccurred())
	})

	It("exits with 0 status code when good spec is passed", func() {
		Eventually(session).Should(gexec.Exit(0))
	})

	Context("when file is missing", func() {
		BeforeEach(func() {
			spec = path.Join(cwd, "fixtures", "file_that_should_not_exist")
		})

		It("exits with error", func() {
			Eventually(session).Should(gexec.Exit())
			Expect(session.ExitCode()).NotTo(Equal(0))
		})
	})

	Context("when bad spec is passed", func() {
		BeforeEach(func() {
			spec = path.Join(cwd, "fixtures", "bad_nginx.yml")
		})

		It("exits with error", func() {
			Eventually(session).Should(gexec.Exit())
			Expect(session.ExitCode()).NotTo(Equal(0))
		})
	})

	Context("when no spec is passed", func() {
		It("exists with error", func() {
			command := exec.Command(pathToCLI)
			session, err := gexec.Start(command, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
			Eventually(session).Should(gexec.Exit())
			Expect(session.ExitCode()).NotTo(Equal(0))
		})
	})

})
