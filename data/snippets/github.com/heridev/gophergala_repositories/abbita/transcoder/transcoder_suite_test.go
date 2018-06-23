package transcoder_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTranscoder(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Transcoder Suite")
}
