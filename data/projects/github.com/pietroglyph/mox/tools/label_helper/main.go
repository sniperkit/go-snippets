package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	scryfall "github.com/BlueMonday/go-scryfall"
	flag "github.com/ogier/pflag"
)

type configuration struct {
	Bind           string
	IndexLocation  string
	OutputLocation string
}

type labelledCard struct {
	UUID            string      `json:"uuid"`
	CardName        boundingBox `json:"cardName"`
	TypeLine        boundingBox `json:"typeLine"`
	CollectorNumber boundingBox `json:"collectorNumber"`
	SetSymbol       boundingBox `json:"setSymbol"`
}

type boundingBox struct {
	StartX int `json:"startX"`
	StartY int `json:"startY"`
	EndX   int `json:"endX"`
	EndY   int `json:"endY"`
}

type helperTextCSS struct {
	FontFamily    string
	LetterSpacing string
	Top           string
	Left          string
	FontSize      string
	PaddingBottom int
}

var (
	config          configuration
	indexTemplate   *template.Template
	ctx             context.Context
	client          *scryfall.Client
	recentImages    map[string][]byte
	recentImagesMux sync.RWMutex
)

func init() {
	flag.StringVarP(&config.Bind, "bind", "b", "localhost:8000", "Host and port to bind the webserver to.")
	flag.StringVarP(&config.IndexLocation, "index-location", "i", "./index.html", "Location of the file to serve on the webserver.")
	flag.StringVarP(&config.OutputLocation, "output-location", "o", "./output/", "Location to output annotations and files.")
}

func main() {
	flag.Parse()

	b, err := ioutil.ReadFile(config.IndexLocation)
	if err != nil {
		log.Panic(err)
	}

	indexTemplate, err = template.New("index").Parse(string(b))
	if err != nil {
		log.Panic(err)
	}

	ctx = context.Background()
	client, err = scryfall.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	recentImages = make(map[string][]byte)

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/ingest", handleDataIngest)

	log.Println("Listening on", config.Bind)
	log.Panic(http.ListenAndServe(config.Bind, nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	var err error
	var card scryfall.Card
	if r.URL.Query().Get("multiverseid") != "" {
		var id int
		id, err = strconv.Atoi(r.URL.Query().Get("multiverseid"))
		if err != nil {
			loggedHTTPError(w, err.Error(), http.StatusBadRequest)
			return
		}
		card, err = client.GetCardByMultiverseID(ctx, id)
		if err != nil {
			loggedHTTPError(w, err.Error(), http.StatusBadRequest)
			return
		}
	} else {
		for {
			card, err = client.GetRandomCard(ctx)
			if err != nil {
				loggedHTTPError(w, err.Error(), http.StatusInternalServerError)
				return
			} else if r.URL.Query().Get("frame") != "" {
				if r.URL.Query().Get("frame") == string(card.Frame) {
					break
				}
			} else if card.Set != "PRM" && card.Set != "TD0" {
				break
			}
		}
	}

	imgResp, err := http.Get(card.ImageURIs.PNG) // This way we only get the image once
	if err != nil {
		loggedHTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	img, err := ioutil.ReadAll(imgResp.Body)
	if err != nil {
		loggedHTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	set, err := client.GetSet(ctx, card.Set)
	if err != nil {
		loggedHTTPError(w, err.Error(), http.StatusInternalServerError)
	}

	recentImagesMux.Lock()
	recentImages[card.ID] = img
	recentImagesMux.Unlock()

	cname := helperTextCSS{}
	tline := helperTextCSS{}
	cnum := helperTextCSS{}

	cname.PaddingBottom = 0
	switch card.Frame {
	case scryfall.Frame1993:
		tline.FontFamily = "MPlantin"
		tline.Top = "588px"
		tline.Left = "82px"
		tline.FontSize = "32px"
		tline.LetterSpacing = "1.5px"

		cname.FontFamily = "GoudyMedieval"
		cname.LetterSpacing = "-0.6px"
		if card.Set == "LEA" || card.Set == "LEB" || card.Set == "2ED" || card.Set == "CED" {
			cname.Top = "59px"
			cname.Left = "60px"
		} else if card.BorderColor == "white" || card.BorderColor == "gold" {
			cname.Top = "61px"
			cname.Left = "57px"
		} else {
			cname.Top = "45px"
			cname.Left = "64px"
			cname.LetterSpacing = "0px"
		}
		cname.FontSize = "43px"
	case scryfall.Frame1997:
		cname.FontFamily = "GoudyMedieval"
		cname.Top = "56px"
		cname.Left = "88px"
		cname.FontSize = "43px"
		cname.LetterSpacing = "0.2px"

		tline.FontFamily = "Palatino Linotype"
		tline.Top = "585px"
		tline.Left = "87px"
		tline.FontSize = "34px"
		tline.LetterSpacing = "-0.3px"

		if set.ReleasedAt.Before(time.Date(1998, time.June, 14, 0, 0, 0, 0, time.UTC)) {
			break
		} else if set.Name == "PALP" || set.Name == "ATH" || set.Name == "PGRU" || set.Name == "BRB" || set.Name == "PSUS" || set.Name == "P00" {
			// These sets don't have collector numbers, but were released after June 15, 1998
			break
		}

		cnum.FontFamily = "MPlantin"
		cnum.FontSize = "16px"
		cnum.Top = "972px"
		cnum.Left = "500px"
	case scryfall.Frame2003:
		cname.FontFamily = "Matrix"
		cname.Top = "75px"
		cname.Left = "70px"
		cname.FontSize = "48px"
		cname.LetterSpacing = "-0.6px"
		cname.PaddingBottom = 8

		tline.FontFamily = "Matrix"
		tline.Top = "610px"
		tline.Left = "79px"
		tline.FontSize = "38px"
		tline.LetterSpacing = "-0.15px"

		cnum.FontFamily = "MPlantin"
		cnum.Top = "990px"
		cnum.Left = "407px"
		cnum.LetterSpacing = "0.25px"
		cnum.FontSize = "17px"
	case scryfall.Frame2015:
		cname.FontFamily = "Beleren"
		cname.Top = "66px"
		cname.Left = "72px"
		cname.FontSize = "40px"
		cname.LetterSpacing = "-0.3px"

		tline.FontFamily = "Beleren"
		tline.Top = "601px"
		tline.Left = "69px"
		tline.FontSize = "34px"
		tline.LetterSpacing = "0px"

		cnum.FontFamily = "Gotham Medium"
		cnum.FontSize = "16px"
		cnum.Top = "983px"
		cnum.Left = "55px"
		cnum.LetterSpacing = "3px"
	}

	args := map[string]interface{}{
		"UUID":                  card.ID,
		"URI":                   template.URL("data:image/png;base64," + base64.StdEncoding.EncodeToString(img)),
		"CardName":              card.Name,
		"TypeLine":              card.TypeLine,
		"CollectorNumber":       card.CollectorNumber + "/" + strconv.Itoa(set.CardCount),
		"Set":                   card.SetName,
		"CardNameHelper":        cname,
		"TypeLineHelper":        tline,
		"CollectorNumberHelper": cnum,
	}

	if card.Frame == scryfall.Frame2015 {
		colNumInt, err := strconv.Atoi(card.CollectorNumber)
		if err == nil {
			args["CollectorNumber"] = fmt.Sprintf("%03d", colNumInt) + "/" + strconv.Itoa(set.CardCount)
		}
	}

	indexTemplate.ExecuteTemplate(w, "Document", args)
}

func handleDataIngest(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		loggedHTTPError(w, err.Error(), http.StatusBadRequest)
		return
	}

	lc := &labelledCard{}
	err = json.Unmarshal(b, lc)
	if err != nil {
		loggedHTTPError(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ioutil.WriteFile(config.OutputLocation+lc.UUID+".json", b, 0644)
	if err != nil {
		loggedHTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	recentImagesMux.RLock()
	err = ioutil.WriteFile(config.OutputLocation+lc.UUID+".png", recentImages[lc.UUID], 0644)
	recentImagesMux.RUnlock()

	recentImagesMux.Lock()
	delete(recentImages, lc.UUID)
	recentImagesMux.Unlock()
	if err != nil {
		loggedHTTPError(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func loggedHTTPError(w http.ResponseWriter, e string, code int) {
	log.Println(strconv.Itoa(code), e)
	loggedHTTPError(w, e, code)
}
