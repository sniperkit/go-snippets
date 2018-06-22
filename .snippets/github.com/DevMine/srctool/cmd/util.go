// Copyright 2014-2015 The DevMine authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"archive/zip"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/ioprogress"

	"github.com/DevMine/srctool/config"
	"github.com/DevMine/srctool/log"
)

func getRemoteParsers(serverURL string) []string {
	var remoteParsers []string

	md5sums, err := fetchChecksumsFile(serverURL)
	if err != nil {
		log.Fatal(err)
	}

	for _, line := range strings.Split(md5sums, "\n") {
		tmp := strings.Split(line, " ")
		if len(tmp) != 2 {
			continue
		}

		path, parser := filepath.Split(tmp[1])
		if isSupported(path) {
			remoteParsers = append(remoteParsers, formatParserName(parser))
		}
	}

	return remoteParsers
}

func getInstalledParsers() []string {
	var installedParsers []string

	fis, err := ioutil.ReadDir(filepath.Join(config.DataDir(), config.ParsersFolder))
	if err != nil {
		log.Fatal(err)
	}

	for _, fi := range fis {
		if !fi.IsDir() {
			continue
		}

		if marched, err := filepath.Match("parser-*", fi.Name()); err != nil {
			log.Debug(err)
			continue
		} else if !marched {
			continue
		}

		installedParsers = append(installedParsers, formatParserName(fi.Name()))
	}

	return installedParsers
}

// isSupported checks whether the current OS and architecture are supported.
func isSupported(path string) bool {
	suppPath := filepath.Join(runtime.GOOS, runtime.GOARCH) + string(filepath.Separator)
	return path == suppPath
}

func formatParserName(fileName string) string {
	return strings.Replace(removeExt(fileName), "parser-", "", -1)
}

func genParserName(lang string) string {
	return fmt.Sprintf("parser-%s", lang)
}

// removeExt removes the extension of a given file name.
func removeExt(fileName string) string {
	ext := filepath.Ext(fileName)
	return fileName[0 : len(fileName)-len(ext)]
}

func downloadParser(serverURL, parserName string, verbose bool) error {
	resp, err := http.Get(config.ParserURI(serverURL, parserName))
	if err != nil {
		log.Debug(err)
		return errors.New("failed to download " + parserName)
	}
	defer resp.Body.Close()

	contentLen := resp.Header.Get("Content-Length")
	size, err := strconv.ParseInt(contentLen, 10, 64)
	if err != nil {
		log.Debug(err)
		return errors.New("malformed or missing Content-Length header")
	}

	out, err := os.Create(config.TempPath(parserName))
	if err != nil {
		return err
	}
	defer out.Close()

	progressR := &ioprogress.Reader{
		Reader:       resp.Body,
		Size:         size,
		DrawInterval: time.Millisecond,
		DrawFunc: func(progress, total int64) error {
			if progress == total {
				// Small hack to make sure that the progress text is up to date
				// at the end of the download.
				fmt.Printf("\rDownloading: %s%10s", ioprogress.DrawTextFormatBytes(size, size), "")
				return nil
			}

			fmt.Printf("\rDownloading: %s%10s", ioprogress.DrawTextFormatBytes(progress, total), "")

			return nil
		},
	}

	if _, err = io.Copy(out, progressR); err != nil {
		log.Debug(err)
		return errors.New("failed to download " + parserName)
	}

	fmt.Println()
	if verbose {
		log.Success(parserName, " successfully downloaded")
	}

	expectedSum, err := fetchChecksum(serverURL, parserName)
	if err != nil {
		return err
	}

	if ok, err := verifyChecksum(config.TempPath(parserName), expectedSum); err != nil {
		return err
	} else if !ok {
		return errors.New("MD5 sum mismatch")
	}

	if verbose {
		log.Success("MD5 sum verified")
	}
	return nil
}

func uncompressParser(parserName, target string) error {
	if _, err := os.Stat(filepath.Join(target, parserName)); os.IsExist(err) {
		return errors.New(parserName + " already installed, if you want to update it, use the 'update' command")
	}

	r, err := zip.OpenReader(config.TempPath(parserName))
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}

		path := filepath.Join(target, f.Name)

		// create target file or directory
		if f.FileInfo().IsDir() {
			if err := os.Mkdir(path, f.Mode()); err != nil {
				return err
			}
			continue
		}

		out, err := os.Create(path)
		if err != nil {
			return err
		}

		if err = out.Chmod(f.Mode()); err != nil {
			return err
		}

		// unzip file content
		if _, err = io.Copy(out, rc); err != nil {
			return err
		}
		rc.Close()
	}

	return nil
}

func deleteParserTarball(parserName string) error {
	return os.Remove(config.TempPath(parserName))
}

func fetchChecksumsFile(serverURL string) (string, error) {
	resp, err := http.Get(config.RemoteChecksumsPath(serverURL))
	if err != nil {
		log.Debug(err)
		return "", errors.New("failed to fetch the MD5SUMS file")
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debug(err)
		return "", errors.New("failed to download MD5SUMS file")
	}

	return string(bs), nil
}

func fetchChecksum(serverURL, parserName string) (string, error) {
	checksums, err := fetchChecksumsFile(serverURL)
	if err != nil {
		return "", err
	}

	for _, line := range strings.Split(checksums, "\n") {
		tmp := strings.Split(line, " ")
		if len(tmp) != 2 {
			return "", errors.New("malformed MD5SUMS file")
		}

		sum, remotePath := tmp[0], tmp[1]

		if remotePath == config.RemoteParserPath(parserName) {
			return sum, nil
		}
	}

	return "", fmt.Errorf("no MD5 sum found for file %s", parserName)
}

// verifyChecksum verifies the checksum of a given file.
func verifyChecksum(path, expectedSum string) (bool, error) {
	md5sum, err := checksum(path)
	if err != nil {
		return false, err
	}

	log.Debug("expected MD5 sum:", expectedSum)
	log.Debug("MD5 sum found:", md5sum)

	return expectedSum == md5sum, nil
}

// checksum computes the MD5 checksum.
func checksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Debug(err)
		return "", errors.New("cannot open file")
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Debug(err)
		return "", fmt.Errorf("failed to compute MD5 sum of %s", path)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
