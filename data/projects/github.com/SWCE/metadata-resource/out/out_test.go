package main_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"fmt"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Out", func() {
	var tmpdir string
	var source string

	var session *gexec.Session
	var outCmd *exec.Cmd

	BeforeEach(func() {
		var err error

		tmpdir, err = ioutil.TempDir("", "out-source")
		Expect(err).NotTo(HaveOccurred())

		source = path.Join(tmpdir, "out-dir")
		os.MkdirAll(source, 0755)
		outCmd = exec.Command(outPath, source)
		fmt.Printf("%s", tmpdir)
	})

	AfterEach(func() {
		os.RemoveAll(tmpdir)
	})

	Context("when executed", func() {

		JustBeforeEach(func() {
			var err error

			session, err = gexec.Start(outCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())
		})

		It("reports error", func() {
			<-session.Exited
			Expect(session.Err).To(gbytes.Say("out should not be used"))
			Expect(session.ExitCode()).To(Equal(1))
		})

	})

})
