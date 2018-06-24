package main

import (
	"bytes"
	"context"
	"errors"
	"image"
	"image/jpeg"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "image/png"

	scryfall "github.com/BlueMonday/go-scryfall"
	"github.com/disintegration/imaging"
	tf "github.com/tensorflow/tensorflow/tensorflow/go"
)

type setSymbol struct {
	Set    string
	Rarity string
}

type deferredFile struct {
	FileInfo os.FileInfo

	containsBytes bool
	bytes         []byte
}

const (
	setSymbolFileExtension = ".jpg"
	setSymbolFileSeparator = "-"

	pastedImagePadding = 5 // In pixels
)

func setupSetSymbols(ctx context.Context, client *scryfall.Client, setSymbolsDir string, backgroundPath string, session *tf.Session, graph *tf.Graph, getNew bool) (map[setSymbol]*deferredFile, error) {
	ssFiles, err := ioutil.ReadDir(setSymbolsDir)
	if err != nil {
		return nil, err
	}

	bgReader, err := os.Open(backgroundPath)
	defer bgReader.Close()
	if err != nil {
		return nil, err
	}

	bg, _, err := image.Decode(bgReader)
	if err != nil {
		return nil, err
	}

	sets, err := client.ListSets(ctx)
	if err != nil {
		return nil, err
	}

	setsMap := make(map[setSymbol]*deferredFile)

S:
	for i := range sets {
		var foundCommon bool
		var foundUncommon bool
		var foundRare bool
		var foundMythicRare bool

		for _, v := range ssFiles {
			ss, err := getSetSymbolFromFileName(v.Name())
			if err != nil || ss.Set != sets[i].Code {
				continue // This probably isn't the file we're looking for
			}

			switch ss.Rarity {
			case "common":
				foundCommon = true
			case "uncommon":
				foundUncommon = true
			case "rare":
				foundRare = true
			case "mythic":
				foundMythicRare = true
			}

			setsMap[ss] = &deferredFile{FileInfo: v}
		}

		if !getNew || (foundCommon && foundUncommon && foundRare && foundMythicRare) {
			continue S
		}

		cards, err := client.SearchCards(ctx, "e:"+sets[i].Code, scryfall.SearchCardsOptions{})
		if err != nil {
			log.Println("Couldn't find any cards for", sets[i].Code)
			continue S
		}

		code := sets[i].Code

		var release *scryfall.Date
		if sets[i].ReleasedAt == nil {
			release = &scryfall.Date{Time: time.Now()}
		} else {
			release = sets[i].ReleasedAt
		}

		if code == "lea" || code == "leb" || code == "2ed" ||
			code == "ced" || code == "cei" || code == "3ed" || code == "sum" ||
			(cards.Cards[0].Frame == scryfall.Frame1993 && cards.Cards[0].BorderColor == "white") {
			foundCommon = true
			foundUncommon = true
			foundRare = true
			foundMythicRare = true
		} else if cards.Cards[0].Frame == scryfall.Frame1993 ||
			release.Before(time.Date(1998, 6, 14, 0, 0, 0, 0, time.UTC)) {
			// Catches cards with the 1997 frame released before Exodus
			foundUncommon = true
			foundRare = true
			foundMythicRare = true
		} else if cards.Cards[0].Frame == scryfall.Frame1997 ||
			release.Before(time.Date(2008, 10, 2, 0, 0, 0, 0, time.UTC)) {
			// Catches cards with either the Exodus and later 1997 frame, or the 2003 frame before Shards of Alara
			// Also catches cards with FrameFuture...
			foundMythicRare = true
		} else if cards.Cards[0].Frame != scryfall.Frame2003 &&
			cards.Cards[0].Frame != scryfall.Frame2015 {
			// This is a fallback for everything else
			foundUncommon = true
			foundRare = true
			foundMythicRare = true
		}

		for _, v := range cards.Cards {
			if foundCommon && foundUncommon && foundRare && foundMythicRare {
				log.Println("Got new set symbols for", sets[i].Code, "("+sets[i].Name+").")
				continue S
			}

			switch v.Rarity {
			case "common":
				if foundCommon {
					continue
				}
				foundCommon = true
			case "uncommon":
				if foundUncommon {
					continue
				}
				foundUncommon = true
			case "rare":
				if foundRare {
					continue
				}
				foundRare = true
			case "mythic":
				if foundMythicRare {
					continue
				}
				foundMythicRare = true
			}

			if v.ImageURIs == nil {
				log.Println("Couldn't find images for", v.ID)
				continue
			}

			resp, err := http.Get(v.ImageURIs.PNG)
			if err != nil {
				return nil, err
			}

			img, _, err := image.Decode(resp.Body)
			if err != nil {
				return nil, err
			}

			// XXX: Assumes that the pasted card isn't too big!
			img = imaging.Paste(bg, img, image.Point{X: pastedImagePadding, Y: pastedImagePadding})

			buf := new(bytes.Buffer)
			err = jpeg.Encode(buf, img, nil)
			if err != nil {
				return nil, err
			}

			crops, classIndicies, err := inferCroppedSections(buf.Bytes(), session, graph)
			if err != nil {
				return nil, err
			}

			var ssCrop image.Image
			for ind := range crops {
				if classIndicies[ind] == setSymbolIndex {
					ssCrop = crops[ind]
					break
				}
			}

			if ssCrop == nil {
				continue // We can't find what we want, try again with another card
			}

			buf = new(bytes.Buffer)
			err = jpeg.Encode(buf, ssCrop, nil)
			if err != nil {
				return nil, err
			}

			ss := setSymbol{Set: v.Set, Rarity: v.Rarity}

			err = ioutil.WriteFile(filepath.Join(setSymbolsDir, ss.fileName()), buf.Bytes(), 0644)
			if err != nil {
				return nil, err
			}

			setsMap[ss] = getDeferredBytes(buf.Bytes())
		}
	}

	return setsMap, nil
}

func getSetSymbolFromFileName(filename string) (setSymbol, error) {
	filename = strings.TrimSuffix(filename, setSymbolFileExtension)
	nameParts := strings.Split(filename, setSymbolFileSeparator)

	if len(nameParts) != 2 {
		return setSymbol{}, errors.New("malformed set symbol filename")
	}

	return setSymbol{Set: nameParts[0], Rarity: nameParts[1]}, nil
}

func (s setSymbol) fileName() string {
	return s.Set + setSymbolFileSeparator + s.Rarity + setSymbolFileExtension
}

func getDeferredBytes(buf []byte) *deferredFile {
	return &deferredFile{containsBytes: true, bytes: buf}
}

func (i *deferredFile) getBytes(baseDir string) ([]byte, error) {
	if i.containsBytes {
		return i.bytes, nil
	}

	buf, err := ioutil.ReadFile(filepath.Join(baseDir, i.FileInfo.Name()))
	if err != nil {
		return nil, err
	}

	i.bytes = buf
	i.containsBytes = true

	return buf, nil
}
