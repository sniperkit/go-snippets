// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"os"

	"github.com/codegangsta/cli"

	"github.com/DevMine/srctool/config"
	"github.com/DevMine/srctool/log"
)

// Delete command deletes one or all language parser(s).
func Delete(c *cli.Context) {
	if !c.Args().Present() {
		deleteAll(c.Bool("dry"))
	} else {
		if err := deleteParser(genParserName(c.Args().First()), c.Bool("dry"), true); err != nil {
			log.Fatal(err)
		}
	}
}

func deleteAll(dryMode bool) {
	parsers := getInstalledParsers()
	for _, parser := range parsers {
		if err := deleteParser(genParserName(parser), dryMode, true); err != nil {
			log.Fail(err)
		}
	}
}

func deleteParser(parserName string, dryMode bool, verbose bool) error {
	parserPath := config.ParserPath(parserName)

	if _, err := os.Stat(parserPath); os.IsNotExist(err) {
		log.Debug(err)
		return errors.New(parserName + " is not installed")
	}

	if dryMode {
		log.Info("parser path:", parserPath)
		return nil
	}

	log.Debug("removing ", parserPath)
	if err := os.RemoveAll(parserPath); err != nil {
		log.Debug(err)
		return errors.New("failed to remove " + parserName)
	}

	if verbose {
		log.Success(parserName, " successfully removed")
	}
	return nil
}
