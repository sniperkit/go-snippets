package main

import (
	"fmt"
	"os"

	"github.com/git-hook/git-ledger"
	"github.com/urfave/cli"
)

// Find path associated with specified path or slug from the ledger.
func find(c *cli.Context) error {
	input := c.Args().First()

	record, err := ledger.GetBySlug(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return err
	}
	fmt.Println(record.Path)

	return nil
}
