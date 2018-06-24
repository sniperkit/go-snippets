package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/git-hook/git-ledger"
	"github.com/urfave/cli"
)

// Remove specified path or slug from the ledger.
func rm(c *cli.Context) error {
	input := c.Args().First()

	// Obtain record by slug or, if DNE, path
	record, err := ledger.GetBySlug(input)
	if err != nil {
		absolutePath, _ := filepath.Abs(input)
		record, err = ledger.GetByPath(absolutePath)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return err
	}

	err = record.RemoveFromLedger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return err
	}

	return nil
}
