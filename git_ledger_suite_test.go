package ledger_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestGitLedger(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GitLedger Suite")
}
