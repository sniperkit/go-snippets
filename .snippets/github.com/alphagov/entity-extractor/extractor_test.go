package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"path"
	"runtime"
	"testing"
)

type ParseExample struct {
	raw           string
	expectedTerms []string
	expectedId    string
	expectedError error
}

var parseExamples = []ParseExample{
	{
		"{}",
		make([]string, 0),
		"",
		nil,
	},
	{
		`{"terms":["Government digital service","GDS"],"id":"1"}`,
		[]string{"Government digital service", "GDS"},
		"1",
		nil,
	},
}

func fixturePath(fixtureFile string) string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "data", fixtureFile)
}

func exampleConfig() *Config {
	return &Config{
		extractAddress:     ":3096",
		dbConnectionString: "host=/var/run/postgresql dbname=entity-extractor_test sslmode=disable",
		logPath:            "STDERR",
	}
}

func TestLoadEntities(t *testing.T) {
	config := exampleConfig()
	extractor := NewExtractor(config)
	err := extractor.LoadEntities()
	if err != nil {
		fmt.Println(err)
		t.FailNow()
	}

	expectedTerms := [...]string{
		"Government digital service",
		"GDS",
		"Ministry of Justice",
		"MoJ",
	}
	require.Equal(t, len(expectedTerms), len(extractor.entities.terms))
	for i, expectedTerm := range expectedTerms {
		assert.Equal(t, expectedTerm, extractor.entities.terms[i])
	}
	expectedIds := [...]string{
		"1",
		"1",
		"2",
		"2",
	}
	for i, expectedId := range expectedIds {
		assert.Equal(t, expectedId, extractor.entities.termOffsetToEntity[i].id)
	}
}

type ExtractionExample struct {
	comment         string
	document        string
	expectedTermIds []string
}

var extractionExamples = []ExtractionExample{
	{
		"a document matching a single word term",
		"This document mentions GDS but it doesn't mention the Ministry of J...",
		[]string{"1"},
	},
	{
		"a document matching a multi-word term",
		"Government digital service",
		[]string{"1"},
	},
	{
		"terms are matched case sensitively",
		"gds",
		[]string{},
	},
}

func TestExtract(t *testing.T) {
	config := exampleConfig()
	extractor := NewExtractor(config)
	err := extractor.LoadEntities()
	if err != nil {
		fmt.Println("ERROR: %v", err)
		t.FailNow()
	}

	for _, example := range extractionExamples {
		matchedTermIds := extractor.Extract(example.document)
		assert.Equal(t, example.expectedTermIds, matchedTermIds, example.comment)
	}
}
