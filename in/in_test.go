package main_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/swce/metadata-resource/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"path/filepath"
)

var _ = Describe("In", func() {
	var tmpdir string
	var destination string

	var inCmd *exec.Cmd

	BeforeEach(func() {
		var err error

		tmpdir, err = ioutil.TempDir("", "in-destination")
		Expect(err).NotTo(HaveOccurred())

		destination = path.Join(tmpdir, "in-dir")

		inCmd = exec.Command(inPath, destination)

		inCmd.Env = append(os.Environ(),
			"BUILD_ID=1",
			"BUILD_NAME=2",
			"BUILD_JOB_NAME=3",
			"BUILD_PIPELINE_NAME=4",
			"BUILD_TEAM_NAME=5",
			"ATC_EXTERNAL_URL=6",
		)
	})

	AfterEach(func() {
		os.RemoveAll(tmpdir)
	})

	Context("when executed", func() {
		var request models.InRequest
		var response models.InResponse

		BeforeEach(func() {

			request = models.InRequest{
				Version: models.TimestampVersion{
					Version: "1",
				},
				Source: models.Source{},
			}

			response = models.InResponse{}
		})

		JustBeforeEach(func() {
			stdin, err := inCmd.StdinPipe()
			Expect(err).NotTo(HaveOccurred())

			session, err := gexec.Start(inCmd, GinkgoWriter, GinkgoWriter)
			Expect(err).NotTo(HaveOccurred())

			err = json.NewEncoder(stdin).Encode(request)
			Expect(err).NotTo(HaveOccurred())

			<-session.Exited
			Expect(session.ExitCode()).To(Equal(0))

			err = json.Unmarshal(session.Out.Contents(), &response)
			Expect(err).NotTo(HaveOccurred())
		})

		It("reports the version to be the input version", func() {
			Expect(response.Version.Version).To(Equal("1"))
		})

		It("writes the requested data the destination", func() {

			checkProp(destination, "build-id", "BUILD_ID", "1", response.Metadata[0])
			checkProp(destination, "build-name", "BUILD_NAME", "2", response.Metadata[1])
			checkProp(destination, "build-job-name", "BUILD_JOB_NAME", "3", response.Metadata[2])
			checkProp(destination, "build-pipeline-name", "BUILD_PIPELINE_NAME", "4", response.Metadata[3])
			checkProp(destination, "build-team-name", "BUILD_TEAM_NAME", "5", response.Metadata[4])
			checkProp(destination, "atc-external-url", "ATC_EXTERNAL_URL", "6", response.Metadata[5])
		})

	})
})

func checkProp(destination string, filename string, prop string, valueToCheck string, meta models.MetadataField) {
	output := filepath.Join(destination, filename)
	file, err := ioutil.ReadFile(output)
	Expect(err).NotTo(HaveOccurred())
	val := string(file)
	Expect(val).To(Equal(valueToCheck))
	Expect(meta.Name).To(Equal(prop))
	Expect(meta.Value).To(Equal(valueToCheck))
}
