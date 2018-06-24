// Package ledger provides primitives for operating on the user's
// git-ledger.
package ledger

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/jmalloc/grit/src/grit/pathutil"
)

// Record entry in the git-ledger.
type Record struct {
	Path string
	Slug string
}

// Struct of Records, for toml parsing.
type records struct {
	Record []Record
}

// Get the path of the git-ledger.
func Path() string {
	home, _ := pathutil.HomeDir()
	return path.Join(home, ".git-ledger")
}

// Return true if dir contains a .git directory.
func IsGitProject(dir string) bool {
	stat, err := os.Stat(path.Join(dir, ".git"))
	isProject := err == nil && stat.IsDir()
	return isProject
}

// Get a list of records currently in the git-ledger.
func GetRecords() (record []Record, err error) {
	if _, err = os.Stat(Path()); os.IsNotExist(err) {
		os.OpenFile(Path(), os.O_RDONLY|os.O_CREATE, 0644)
	}

	var b []byte
	b, err = ioutil.ReadFile(Path())
	if err != nil {
		return
	}
	str := string(b)

	var ledger records
	if _, err = toml.Decode(str, &ledger); err != nil {
		return
	}

	record = ledger.Record
	return
}

// Look-up a record from the git-ledger by slug.
func GetBySlug(slug string) (Record, error) {
	var match Record

	records, err := GetRecords()
	if err != nil {
		return match, err
	}
	for _, r := range records {
		if strings.Contains(r.Slug, slug) {
			return r, nil
		}
	}
	return match, errors.New(fmt.Sprintf("Unknown project: %s", slug))
}

// Look-up a record from the git-ledger by path.
func GetByPath(path string) (Record, error) {
	var match Record

	records, err := GetRecords()
	if err != nil {
		return match, err
	}
	for _, r := range records {
		if r.Path == path {
			return r, nil
		}
	}
	return match, errors.New(fmt.Sprintf("Unknown project: %s", path))
}

// Return the current record as a toml string.
func (r Record) String() string {
	return fmt.Sprintf("[[Record]]\npath = \"%s\"\nslug = \"%s\"\n\n", r.Path, r.Slug)
}

// Remove record from the git-ledger.  Will remove all records
// matching this records path from the ledger.
func (r Record) RemoveFromLedger() error {
	records, err := GetRecords()
	if err != nil {
		return err
	}

	var content string
	for _, s := range records {
		if s.Path != r.Path {
			content = content + s.String() + "\n"
		}
	}

	ioutil.WriteFile(Path(), []byte(content), 0644)
	return nil
}

// Add record to the git-ledger.
func (r Record) AddToLedger() {
	f, err := os.OpenFile(Path(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fmt.Fprintf(f, r.String())
}
