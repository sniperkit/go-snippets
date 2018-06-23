package main

import (
	"database/sql"
	"encoding/json"
	"github.com/cloudflare/ahocorasick"
	_ "github.com/lib/pq"
)

// An entity (which has an ID and several representative strings)
type Entity struct {
	id string
}

type Entities struct {
	terms              []string
	termOffsetToEntity []*Entity
}

type Extractor struct {
	cfg      *Config
	entities Entities
	matcher  *ahocorasick.Matcher
}

func NewExtractor(cfg *Config) *Extractor {
	return &Extractor{cfg: cfg}
}

func (self *Entities) addEntity(id string, termsJson string) error {
	entity := Entity{id}

	var terms []string

	err := json.Unmarshal([]byte(termsJson), &terms)
	if err != nil {
		return err
	}

	for _, term := range terms {
		self.terms = append(self.terms, term)
		self.termOffsetToEntity = append(self.termOffsetToEntity, &entity)
	}

	return nil
}

func (extr *Extractor) LoadEntities() error {
	return extr.loadEntitiesFromDatabase(extr.cfg.dbConnectionString)
}

func (extr *Extractor) loadEntitiesFromDatabase(dbConnectionString string) error {
	db, err := sql.Open("postgres", dbConnectionString)
	if err != nil {
		return err
	}

	rows, err := db.Query("SELECT id, terms FROM entities")
	if err != nil {
		return err
	}
	defer rows.Close()

	var entities Entities

	for rows.Next() {
		var id, terms string
		if err := rows.Scan(&id, &terms); err != nil {
			return err
		}

		if err := entities.addEntity(id, terms); err != nil {
			return err
		}
	}

	if err := rows.Err(); err != nil {
		return err
	}

	extr.matcher = ahocorasick.NewStringMatcher(entities.terms)
	logInfo("Loaded", len(entities.terms), "terms")

	extr.entities = entities

	return nil
}

func (extr *Extractor) Extract(text string) []string {
	matchIndexes := extr.matcher.Match([]byte(text))
	matchingTermIds := make([]string, 0, 1000)
	for _, termIndex := range matchIndexes {
		matchingTermIds = append(matchingTermIds, extr.entities.termOffsetToEntity[termIndex].id)
	}

	return matchingTermIds
}
