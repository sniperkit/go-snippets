package kiwix

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
)

// url to scrape
const kiwixContentListUrl = "http://wiki.kiwix.org/wiki/Content_in_all_languages"

// Zim represents a Zim package on the Kiwix service
type Zim struct {
	// Name of the project that's been archived
	Project string
	// ISO Language code
	Language string
	// Size String (eg: "11M" or "5.4G")
	Size string
	// Created Date in form YYYY-MM
	Created string
	// Weather it's all content or not
	FullContent bool
	// url to get MD5 entry
	Md5Url string
	// md5 checksum for download
	Md5 string
	// download url
	Url string
}

func (z *Zim) Title() string {
	return fmt.Sprintf("%s - %s", z.Project, z.Language)
}

func (z *Zim) FetchMd5() error {
	res, err := http.Get(z.Md5Url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	md5 := make([]byte, 32)
	res.Body.Read(md5)
	z.Md5 = string(md5)
	return nil
}

func FetchZims() ([]*Zim, error) {
	zims := []*Zim{}

	doc, err := goquery.NewDocument(kiwixContentListUrl)
	if err != nil {
		return nil, err
	}

	// relatively brittle scraping code, will need to be kept in sync kiwixContentListUrl
	doc.Find("#zimtable tbody tr").Each(func(i int, el *goquery.Selection) {
		if i == 0 {
			// skip header row
			return
		}

		zim := &Zim{}

		el.ChildrenFiltered("td").Each(func(j int, cell *goquery.Selection) {
			switch j {
			case 0:
				zim.Project = strings.TrimSpace(cell.Text())
			case 1:
				zim.Language = strings.TrimSpace(cell.Text())
			case 2:
				zim.Size = strings.TrimSpace(cell.Text())
			case 3:
				zim.Created = strings.TrimSpace(cell.Text())
			case 4:
				zim.FullContent = strings.TrimSpace(cell.Text()) == "all"
			case 5:
				zim.Url, _ = el.Find("a").Attr("href")
				md5Url, _ := el.Find("sup small a").Attr("href")
				if md5Url != "" {
					zim.Md5Url = md5Url
				}
			}
		})

		zims = append(zims, zim)
	})

	return zims, nil
}
