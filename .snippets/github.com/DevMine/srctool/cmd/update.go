// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"io/ioutil"
	"os"

	"github.com/codegangsta/cli"

	"github.com/DevMine/srctool/config"
	"github.com/DevMine/srctool/log"
)

// Update command updates one or all installed parser(s).
func Update(c *cli.Context) {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	if !c.Args().Present() {
		updateAll(cfg)
	} else {
		updateParser(cfg, genParserName(c.Args().First()))
	}
}

func updateAll(cfg *config.Config) {
	parsers := getInstalledParsers()
	for _, parser := range parsers {
		updateParser(cfg, genParserName(parser))
	}
}

func updateParser(cfg *config.Config, parserName string) {
	if !isAlreadyInstalled(parserName) {
		log.Fail(" parser " + parserName + " not installed, install it first")
		return
	}

	LocalChecksum, err := ioutil.ReadFile(config.LocalChecksumPath(parserName))
	if err != nil {
		log.Debug(err)
		log.Fail("unable to read the MD5 file of the currently installed " + parserName + " parser")
		return
	}

	remoteCheckSum, err := fetchChecksum(cfg.DownloadServerURL, parserName)
	if err != nil {
		log.Fail(err)
		return
	}

	if string(LocalChecksum) == remoteCheckSum {
		log.Info("latest version of " + parserName + " already installed")
		return
	}

	if err = deleteParser(parserName, false, false); err != nil {
		log.Fail(err)
		return
	}

	if err = installParser(cfg, parserName, false); err != nil {
		log.Fail(err)
		return
	}

	log.Success("parser " + parserName + " successfully updated")
}

func isAlreadyInstalled(parserName string) bool {
	parserPath := config.ParserPath(parserName)
	log.Debug("=> ", parserPath)

	if _, err := os.Stat(parserPath); err != nil {
		log.Debug(err)
		return false
	}

	return true
}
