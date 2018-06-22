// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"io/ioutil"

	"github.com/codegangsta/cli"

	"github.com/DevMine/srctool/config"
	"github.com/DevMine/srctool/log"
)

// Install command installs one or all language parser(s).
func Install(c *cli.Context) {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	installedParsers := make(map[string]struct{})
	for _, parser := range getInstalledParsers() {
		installedParsers[parser] = struct{}{}
	}

	if !c.Args().Present() {
		installAll(cfg, installedParsers)
	} else {
		parser := c.Args().First()
		parserName := genParserName(parser)
		if _, ok := installedParsers[parser]; !ok {
			if err := installParser(cfg, parserName, true); err != nil {
				log.Fatal(err)
			}
		} else {
			log.Info(parserName, " already installed")
		}
	}
}

func installAll(cfg *config.Config, installedParsers map[string]struct{}) {
	parsers := getRemoteParsers(cfg.DownloadServerURL)
	for _, parser := range parsers {
		if _, ok := installedParsers[parser]; !ok {
			if err := installParser(cfg, genParserName(parser), true); err != nil {
				log.Fail(err)
			}
		} else {
			log.Info(genParserName(parser), " already installed")
		}
	}
}

func installParser(cfg *config.Config, parserName string, verbose bool) error {
	if err := downloadParser(cfg.DownloadServerURL, parserName, verbose); err != nil {
		return err
	}

	defer func() {
		if err := deleteParserTarball(parserName); err != nil {
			log.Fail(err)
		}
	}()

	if err := uncompressParser(parserName, config.ParsersDir()); err != nil {
		log.Debug(err)
		return errors.New("failed to uncompress the " + parserName + " archive")
	}

	md5sum, err := checksum(config.TempPath(parserName))
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(config.LocalChecksumPath(parserName), []byte(md5sum), 0644); err != nil {
		log.Debug(err)
		return errors.New("failed to write MD5SUM file in the parser directory")
	}

	if verbose {
		log.Success(parserName, " successfully installed")
	}
	return nil
}
