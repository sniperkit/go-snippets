package cmd

import (
	"bytes"
	"gopkg.in/check.v1"
	"os"
	"os/exec"
	"testing"
)

func Test(t *testing.T) { check.TestingT(t) }

type S struct {
	stdin   *os.File
	recover []string
}

var _ = check.Suite(&S{})
var manager *Manager

func (s *S) SetUpSuite(c *check.C) {
	targetFile := os.Getenv("HOME") + "/.megam"
	_, err := os.Stat(targetFile)
	if err == nil {
		old := targetFile + ".old"
		s.recover = []string{"mv", old, targetFile}
		exec.Command("mv", targetFile, old).Run()
	} else {
		s.recover = []string{"rm", targetFile}
	}
	f, err := os.Create(targetFile)
	c.Assert(err, check.IsNil)
	f.Write([]byte("http://localhost"))
	f.Close()
}

func (s *S) TearDownSuite(c *check.C) {
	exec.Command(s.recover[0], s.recover[1:]...).Run()
}

func (s *S) SetUpTest(c *check.C) {
	var stdout, stderr bytes.Buffer
	manager = NewManager("gulpd", "0.1", "", &stdout, &stderr, os.Stdin)
	var exiter recordingExiter
	manager.e = &exiter
}
