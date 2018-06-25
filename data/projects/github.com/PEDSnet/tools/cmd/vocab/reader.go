package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	ErrDateFormat = errors.New("could not parse date")

	tsvExt = regexp.MustCompile(`(?i)\.tsv\b`)

	zeroTime = time.Time{}

	dateFormats = []string{
		"20060102",
		"2006-01-02",
	}
)

func parseDate(s string) (time.Time, error) {
	for _, l := range dateFormats {
		t, err := time.Parse(l, s)

		if err != nil {
			return t, nil
		}
	}

	return zeroTime, ErrDateFormat
}

func detectDelimiter(s string) rune {
	if tsvExt.MatchString(s) {
		return '\t'
	}

	// Assume CSV.
	return ','
}

// UniversalReader wraps an io.Reader and replaces carriage returns with newlines.
type UniversalReader struct {
	r io.Reader
}

func (c *UniversalReader) Read(buf []byte) (int, error) {
	n, err := c.r.Read(buf)

	// Replace carriage returns with newlines
	for i, b := range buf {
		if b == '\r' {
			buf[i] = '\n'
		}
	}

	return n, err
}

// NewUniversalReader returns a UniversalReader that wraps the passed io.Reader.
func NewUniversalReader(r io.Reader) *UniversalReader {
	return &UniversalReader{r}
}

type Header struct {
	ConceptID       int
	ConceptName     int
	DomainID        int
	VocabularyID    int
	ConceptClassID  int
	StandardConcept int
	ConceptCode     int
	ConceptLevel    int
	ValidStartDate  int
	ValidEndDate    int
	InvalidReason   int

	raw []string
}

func parseHeader(head []string) *Header {
	h := Header{
		raw: head,
	}

	for i, col := range head {
		switch strings.ToLower(col) {
		case "concept_id":
			h.ConceptID = i
		case "concept_name":
			h.ConceptName = i
		case "domain_id":
			h.DomainID = i
		case "vocabulary_id":
			h.VocabularyID = i
		case "concept_class_id":
			h.ConceptClassID = i
		case "standard_concept":
			h.StandardConcept = i
		case "concept_code":
			h.ConceptCode = i
		case "valid_start_date":
			h.ValidStartDate = i
		case "valid_end_date":
			h.ValidEndDate = i
		case "invalid_reason":
			h.InvalidReason = i
		case "concept_level":
			h.ConceptLevel = i
		}
	}

	return &h
}

type ConceptReader struct {
	csv.Reader

	head *Header
	line int
}

func (r *ConceptReader) parse(row []string) (*Concept, error) {
	if len(row) != len(r.head.raw) {
		return nil, fmt.Errorf("Wrong number of values on line %d. Expected %d got %d\n%s", len(r.head.raw), r.line+1, len(row), strings.Join(row, ","))
	}

	id, err := strconv.Atoi(row[r.head.ConceptID])

	if err != nil {
		return nil, fmt.Errorf("Error parsing concept_id: %s", err)
	}

	var (
		sd time.Time
		ed time.Time
	)

	if row[r.head.ValidStartDate] == "" {
		sd = zeroTime
	} else {
		sd, err = parseDate(row[r.head.ValidStartDate])

		if err != nil {
			return nil, fmt.Errorf("Error parsing valid_start_date: %s", row[r.head.ValidStartDate])
		}
	}

	if row[r.head.ValidEndDate] == "20991231" || row[r.head.ValidEndDate] == "" {
		ed = zeroTime
	} else {
		ed, err = parseDate(row[r.head.ValidEndDate])

		if err != nil {
			return nil, fmt.Errorf("Error parsing valid_end_date: %s", row[r.head.ValidEndDate])
		}
	}

	c := Concept{
		ConceptID:       id,
		ConceptName:     strings.TrimSpace(row[r.head.ConceptName]),
		DomainID:        strings.TrimSpace(row[r.head.DomainID]),
		VocabularyID:    strings.TrimSpace(row[r.head.VocabularyID]),
		ConceptClassID:  strings.TrimSpace(row[r.head.ConceptClassID]),
		StandardConcept: strings.TrimSpace(row[r.head.StandardConcept]),
		ConceptCode:     strings.TrimSpace(row[r.head.ConceptCode]),
		ValidStartDate:  sd,
		ValidEndDate:    ed,
		InvalidReason:   strings.TrimSpace(row[r.head.InvalidReason]),
	}

	return &c, nil
}

func (r *ConceptReader) Read() (*Concept, error) {
	row, err := r.Reader.Read()

	// End of file.
	if err == io.EOF {
		return nil, err
	}

	if err != nil {
		return nil, fmt.Errorf("Could not read line: %s", err)
	}

	r.line++

	// Sentinel? Skip.
	if row[0] == "0" {
		return r.Read()
	}

	return r.parse(row)
}

func NewConceptReader(r io.Reader) (*ConceptReader, error) {
	cr := ConceptReader{
		Reader: *csv.NewReader(NewUniversalReader(r)),
	}

	cr.Comment = '#'
	cr.LazyQuotes = true
	cr.TrimLeadingSpace = true
	cr.FieldsPerRecord = -1

	// Read the header.
	head, err := cr.Reader.Read()

	if err != nil {
		return nil, err
	}

	cr.line++
	cr.head = parseHeader(head)

	return &cr, nil
}
