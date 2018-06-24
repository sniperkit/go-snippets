package core_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestSuite(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "core")
}

type (
	customNumeric int
	customString  string
)

func intPtr(n int) *int { return &n }
