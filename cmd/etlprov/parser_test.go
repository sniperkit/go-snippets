package main

import (
	"testing"

	dms "github.com/chop-dbhi/data-models-service/client"
)

type errorMap map[string][]error

func (m errorMap) Handle(name string, errs []error) {
	m[name] = errs
}

var model *dms.Model

func init() {
	client, _ := dms.New("http://data-models.origins.link")
	model, _ = client.ModelRevision("pedsnet", "2.0.0")
}

func TestEntityParser(t *testing.T) {
	p := NewParser(model)

	row := []string{
		"care_site.care_site_source_value",
		"available",
		"yes",
		"The care site name and the department id in the source",
		"",
		"",
	}

	_, err := p.parseEntity(row)

	if err != nil {
		t.Errorf("error parsing row: %s\n%v", err, row)
	}
}

func TestStepParser(t *testing.T) {
	p := NewParser(model)

	// Sentinel entity for reference.
	e := &Entity{
		Name: "person.person_id",
	}

	p.entities[e.Name] = e

	row := []string{
		"1",
		"Extracing race information for the cohort",
		"person.person_id",
		"",
		"",
	}

	s1, err := p.parseStep(row)

	if err != nil {
		t.Errorf("error parsing row: %s\n%v", err, row)
	}

	if s1.ID != 1 {
		t.Errorf("expected ID %d, got %d", 1, s1.ID)
	}

	p.steps[s1.ID] = s1

	row[0] = "2"
	row[3] = "1"

	s2, err := p.parseStep(row)

	if err != nil {
		t.Errorf("error parsing row: %s\n%v", err, row)
	}

	if s2.ID != 2 {
		t.Errorf("expected ID %d, got %d", 2, s2.ID)
	}

	if s2.PreviousStep != s1 {
		t.Error("previous step does not match")
	}
}

func TestParser(t *testing.T) {
	p := NewParser(model)

	m := make(errorMap)

	if err := p.ReadDir("./test_data", m.Handle); err != nil {
		t.Fatal(err)
	}

	for k, errs := range m {
		if len(errs) != 0 {
			t.Errorf("got %d errors for %s", len(errs), k)
			t.Errorf("%v", errs)
		}
	}
}
