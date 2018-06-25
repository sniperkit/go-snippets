package main

import (
	"fmt"
	"time"
)

type Concept struct {
	ConceptID       int
	ConceptName     string
	DomainID        string
	VocabularyID    string
	ConceptClassID  string
	ConceptLevel    string
	StandardConcept string
	ConceptCode     string
	ValidStartDate  time.Time
	ValidEndDate    time.Time
	InvalidReason   string
}

func (c *Concept) Row() []string {
	return []string{
		fmt.Sprintf("%d", c.ConceptID),
		c.ConceptName,
		c.DomainID,
		c.VocabularyID,
		c.ConceptClassID,
		c.ConceptLevel,
		c.StandardConcept,
		c.ConceptCode,
		c.ValidStartDate.String(),
		c.ValidEndDate.String(),
		c.InvalidReason,
	}
}
