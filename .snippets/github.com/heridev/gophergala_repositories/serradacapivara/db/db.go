package db

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Site is an archeological site
type Site struct {
	Id              string
	Name            string
	HasPicture      bool
	HasEngraving    bool
	Other           bool
	IsHistoric      bool
	YearOfDiscovery int
	Latitude        string
	Longitude       string
	City            string
	Circuit         string
	NationalPark    string
	Location        string
	Description     string
	Pictures        []string
}

type Database []Site

var DB Database

// Find the Site by its identifier
func FindByID(id string) (Site, error) {
	for _, s := range DB {
		if s.Id == id {
			return s, nil
		}
	}

	return Site{}, errors.New("Site not found")
}

func AllSites() []Site {
	return DB
}

func Search(query string) []Site {
	var sites []Site

	for _, site := range AllSites() {
		if strings.Contains(strings.ToUpper(site.Name), query) || strings.Contains(strings.ToUpper(site.Description), query) {
			sites = append(sites, site)
		}
	}

	return sites
}

// Initialize the database
func init() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))

	if err != nil {
		panic("Something very bad just happen.")
	}

	file, err := ioutil.ReadFile(path.Join(dir, "db/data.json"))

	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(file, &DB); err != nil {
		panic(err)
	}
}
