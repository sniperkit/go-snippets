/*
ViewLater Chrome Extension
Copyright (C) 2017 Andrii Zablodskyi (andrey.zablodskiy@gmail.com)

This program is free software; you can redistribute it and/or
modify it under the terms of the GNU General Public License
as published by the Free Software Foundation; either version 2
of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
*/
package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"github.com/fabioberger/chrome"
	"github.com/gopherjs/gopherjs/js"
	"github.com/nfnt/resize"
	"honnef.co/go/js/dom"
	"image/jpeg"
	"strconv"
)

type Link struct {
	Id          int    `json:"id"`
	Url         string `json:"url"`
	Title       string `json:"title"`
	Description string `json:"desc"`
	Screenshot  string `json:"screenshot"`
}

type LinksList []Link

const (
	THUMBNAIL_IMAGE_WIDTH  = 640
	THUMBNAIL_IMAGE_HEIGHT = 480
	IMAGE_DATA_START_INDEX = 23
)

var (
	storage    = js.Global.Get("localStorage")
	urls       = LinksList{}
	screenshot = ""
)

func main() {

	c := chrome.NewChrome()
	d := dom.GetWindow().Document()

	rows := d.GetElementByID("links")

	// get json from storage
	linksJson := storage.Get("readLaterLinks").String()

	if linksJson != "" {
		err := json.Unmarshal([]byte(linksJson), &urls)
		if err != nil {
			println("Error has been occurred ", err.Error())
		}
	}

	for _, link := range urls {
		row := d.CreateElement("div").(*dom.HTMLDivElement)
		row.SetClass("row row-link")

		addScreenshot(d, row, link.Screenshot, false)
		addTitle(d, row, link.Id, link.Url, link.Title)
		addRowButtons(d, rows, row, link.Id, link.Url)

		rows.AppendChild(row)
	}

	c.Tabs.CaptureVisibleTab(c.Windows.WINDOW_ID_CURRENT, nil, func(dataUrl string) {
		screenshot = dataUrl
	})

	add := d.GetElementByID("add-button").(*dom.HTMLAnchorElement)
	add.Call("addEventListener", "click", func(event *js.Object) {
		row := d.CreateElement("div").(*dom.HTMLDivElement)
		row.SetClass("row row-link")

		c.Windows.GetCurrent(chrome.Object{}, func(window chrome.Window) {
			id := window.Id

			c.Tabs.GetSelected(id, func(tab chrome.Tab) {

				currentUrl := tab.Url

				nextId := 0
				if len(urls) > 0 {
					nextId = urls[len(urls)-1].Id + 1
				}

				addScreenshot(d, row, screenshot, true)
				addTitle(d, row, nextId, currentUrl, tab.Title)
				addRowButtons(d, rows, row, nextId, currentUrl)

				rows.AppendChild(row)

				urls = append(urls, Link{Id: nextId, Url: currentUrl, Title: tab.Title, Screenshot: screenshot})
				marshalUrlsToStorage()
			})

		})
	})
}

func addScreenshot(d dom.Document, row *dom.HTMLDivElement, screenshot string, newScreenshot bool) {
	div := d.CreateElement("div").(*dom.HTMLDivElement)
	div.SetClass("col-6 d-flex flex-column")
	row.AppendChild(div)

	img := d.CreateElement("img").(*dom.HTMLImageElement)
	img.SetClass("thumbnail")

	if newScreenshot {
		img.Src = resizeScreenshot(screenshot)
	} else {
		img.Src = screenshot
	}

	div.AppendChild(img)
}

func addTitle(d dom.Document, row *dom.HTMLDivElement, id int, url, desc string) {
	div := d.CreateElement("div").(*dom.HTMLDivElement)
	div.SetClass("col-4 d-flex flex-column")
	row.AppendChild(div)

	p := d.CreateElement("p").(*dom.HTMLParagraphElement)
	div.AppendChild(p)

	idStr := strconv.Itoa(id)

	titleLink := d.CreateElement("a").(*dom.HTMLAnchorElement)
	titleLink.SetID(idStr)
	titleLink.Href = url
	titleLink.SetInnerHTML(desc)
	p.AppendChild(titleLink)

	// "hidden" text area that will be used for clipboard copy
	textArea := d.CreateElement("textarea").(*dom.HTMLTextAreaElement)
	textArea.SetID("ta-" + idStr)
	textArea.SetTextContent(url)
	div.AppendChild(textArea)
}

func addRowButtons(d dom.Document, rows dom.Element, row *dom.HTMLDivElement, urlId int, url string) {
	// copy row url to clipboard
	div := d.CreateElement("div").(*dom.HTMLDivElement)
	div.SetClass("col-2 d-flex flex-column align-items-right text-right")
	row.AppendChild(div)

	p := d.CreateElement("p").(*dom.HTMLParagraphElement)
	div.AppendChild(p)

	copyLink := d.CreateElement("a").(*dom.HTMLAnchorElement)
	copyLink.SetClass("btn btn-default btn-sm btn-light")
	copyLink.SetInnerHTML("<i class=\"fa fa-clipboard\"></i> Copy")
	copyLink.Call("addEventListener", "click", func(event *js.Object) {
		document := js.Global.Get("document")
		textArea := document.Call("getElementById", "ta-"+strconv.Itoa(urlId))
		textArea.Call("select")

		document.Call("execCommand", "copy")
	})
	p.AppendChild(copyLink)

	// remove row
	removeLink := d.CreateElement("a").(*dom.HTMLAnchorElement)
	removeLink.SetClass("btn btn-default btn-sm btn-danger")
	removeLink.SetInnerHTML("<i class=\"fa fa-trash-o\"></i> Delete")
	removeLink.Call("addEventListener", "click", func(event *js.Object) {
		removeUrl(urlId)
		marshalUrlsToStorage()
		rows.RemoveChild(row)
	})
	p.AppendChild(removeLink)
}

func removeUrl(id int) {
	modified := LinksList{}
	for _, link := range urls {
		if link.Id == id {
			continue
		}

		modified = append(modified, link)
	}

	urls = modified
}

func resizeScreenshot(screenshot string) string {
	decodedBytes, err := base64.StdEncoding.DecodeString(screenshot[IMAGE_DATA_START_INDEX:])
	if err != nil {
		println("Cannot decode base64 bytes." + err.Error())
		return screenshot
	}

	image, err := jpeg.Decode(bytes.NewReader(decodedBytes))
	if err != nil {
		println("Cannot decode screenshot image." + err.Error())
		return screenshot
	}

	resizedImage := resize.Resize(THUMBNAIL_IMAGE_WIDTH, THUMBNAIL_IMAGE_HEIGHT, image, resize.Lanczos3)

	var buf bytes.Buffer

	jpeg.Encode(&buf, resizedImage, nil)

	return "data:image/jpeg;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
}

func marshalUrlsToStorage() {
	bytes, err := json.Marshal(urls)
	if err != nil {
		println("Error has occurred during marshalling ", err.Error())
	}

	storage.Set("readLaterLinks", string(bytes))
}
