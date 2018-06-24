package ledger_test

import (
	. "github.com/git-hook/git-ledger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"os"
	"fmt"
	"path"
)

var (
	r Record
	s Record
)

// Move the current git-ledger to a new location so testing can create
// a temporary one
func pushGitLedger() {
	os.Rename(Path(), fmt.Sprintf("%s.backup", Path()))
	f, _ := os.Create(Path())
	defer f.Close()
}

// Remove the temporary git-ledger and replace it with the backup
func popGitLedger() {
	os.Remove(Path())
	os.Rename(fmt.Sprintf("%s.backup", Path()), Path())
}

// Create a dummy ledger-file to act as the device-under-test.
func createDummyLedger() {
	r = Record{Path: "/tmp", Slug: "pony/fill"}
	r.AddToLedger()
	s = Record{Path: "/home", Slug: "badger/bear"}
	s.AddToLedger()
}

var _ = Describe("GitLedger", func() {

	BeforeEach(func() {
		pushGitLedger()
		createDummyLedger()
	})
	AfterEach(popGitLedger)

	// TODO: test that git-ledger silently creates the ledger if asked to read it

	Describe("CLI", func() {
		Context("add", func() {
			// TODO: test about duplicates
			Specify("should accept fully-initialized records", func() {
				Record{Path: "/proc", Slug: "tron/man"}.AddToLedger()
				records, _ := GetRecords()
				Expect(len(records)).To(Equal(3))
				Record{Path: "/dev", Slug: "cron/man"}.AddToLedger()
				records, _ = GetRecords()
				Expect(len(records)).To(Equal(4))
			})
		})

		Context("remove", func() {
			Specify("should accept fully-initialized records", func() {
				s.RemoveFromLedger()
				records, _ := GetRecords()
				Expect(len(records)).To(Equal(1))
				r.RemoveFromLedger()
				records, _ = GetRecords()
				Expect(len(records)).To(Equal(0))
			})
		})

		// Context("find", func() {
		// 	var (
		// 		r Record
		// 		s Record
		// 	)

		// 	Specify("")
		// })
	})

	Context("The ledger", func() {
		Specify("should reside in the user's home directory", func() {
			Expect(path.Join(os.Getenv("HOME"), ".git-ledger")).To(Equal(Path()))
		})
		Specify("should be created if it is accessed and it does not exist", func() {
			os.Remove(Path())
			records, err := GetRecords()
			Expect(err).To(BeNil())
			Expect(len(records)).To(Equal(0))
		})
	})

	Describe("Records", func() {

		Context("When the ledger is empty", func() {
			It("should return 0 records", func() {
				os.Remove(Path())
				records, _ := GetRecords()
				Expect(len(records)).To(Equal(0))
			})
		})

		// fixme: make this format more awesome
		Specify("should print in an expected format", func() {
			path := "/path/to/nowhere"
			slug := "sclopio/peepio"
			record := Record{Path: path, Slug: slug}
			expected := fmt.Sprintf("[[Record]]\npath = \"%s\"\nslug = \"%s\"\n\n", record.Path, record.Slug)
			Expect(record.String()).To(Equal(expected))
		})

		Specify("should be index-able by slug", func() {
			rec, _ := GetBySlug("fill")
			Expect(rec).To(Equal(r))
			rec, _ = GetBySlug("pony/fill")
			Expect(rec).To(Equal(r))
			rec, _ = GetBySlug("badger")
			Expect(rec).To(Equal(s))
		})

		Specify("should be index-able by path", func() {
			rec, _ := GetByPath("/tmp")
			Expect(rec).To(Equal(r))
			rec, _ = GetByPath("/home")
			Expect(rec).To(Equal(s))
			_, err := GetByPath("DNE")
			Ω(err).ShouldNot(BeNil())
		})
	})

	Describe("Detecting git meta-data", func() {
		Context("Git repositories", func() {
			var (
				dir string
				gitdir string
			)

			BeforeEach(func() {
				dir = "/tmp"
				gitdir = path.Join(dir, ".git")
			})

			AfterEach(func() {
				os.Remove(gitdir)
			})

			Specify("should identify as such", func() {
				os.MkdirAll(gitdir, os.ModePerm)
				Ω(IsGitProject(dir)).Should(BeTrue())
			})

			Specify("should be differentiable from non-git repositories", func() {
				Ω(IsGitProject(dir)).Should(BeFalse())
			})
		})
	})
})
