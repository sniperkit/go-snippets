package main

import (
	"fmt"
	"os"

	"github.com/git-hook/git-ledger"
	"github.com/urfave/cli"
)

// List all records in the git-ledger.
func ls(c *cli.Context) error {
	records, err := ledger.GetRecords()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return err
	}

	for _, rec := range records {
		fmt.Println(fmt.Sprintf("%s: %s", rec.Slug, rec.Path))
	}

	return nil
}
