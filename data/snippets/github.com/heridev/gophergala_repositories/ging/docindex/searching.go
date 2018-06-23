package docindex

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"path"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search"
)

// SearchResult ...
// TODO(alvivi): doc this
type SearchResult struct {
	Name       string
	Type       DocKind
	Link       string
	Match      string
	Highlights SearchHighlights
}

// SearchHighlights ...
// TODO(alvivi): doc this
type SearchHighlights struct {
	Name    template.HTML
	Content template.HTML
}

// Search ...
// TODO(alvivi): doc this
func Search(index bleve.Index, queryString string) ([]*SearchResult, *bleve.SearchResult, error) {
	matchQuery := bleve.NewMatchPhraseQuery(queryString)
	entries, sr, err := performSearch(index, matchQuery)
	if err != nil {
		return nil, nil, err
	}
	if sr.Total > 0 {
		return entries, sr, nil
	}
	fuzzyQuery := bleve.NewMatchQuery(queryString)
	fuzzyQuery.SetFuzziness(2)
	entries, fsr, err := performSearch(index, fuzzyQuery)
	if err != nil {
		return nil, nil, err
	}
	fsr.Merge(sr)
	return entries, fsr, nil
}

func performSearch(index bleve.Index, query bleve.Query) ([]*SearchResult, *bleve.SearchResult, error) {
	search := bleve.NewSearchRequest(query)
	search.Fields = []string{
		"name",
		"kind",
		"import",
	}
	search.Highlight = bleve.NewHighlightWithStyle("html")
	search.Explain = false
	sr, err := index.Search(search) // sr, err := ...
	if err != nil {
		return []*SearchResult{}, nil, err
	}

	entries := []*SearchResult{}
	for _, hit := range sr.Hits {
		entry, err := newSearchResult(hit.Fields, hit.Fragments)
		if err == nil {
			entries = append(entries, entry)
		} else {
			log.Printf("Error building a search result entry: %s.\n", err.Error())
		}
	}
	return entries, sr, nil
}

func newSearchResult(fields map[string]interface{}, fragments search.FieldFragmentMap) (*SearchResult, error) {
	// Name
	nameValue, ok := fields["name"]
	if !ok {
		return nil, errors.New("Required field 'name' not found")
	}
	name := nameValue.(string)
	// Type
	doctypeValue, ok := fields["kind"]
	if !ok {
		return nil, errors.New("Required field 'kind' not found")
	}
	doctype := DocKind(doctypeValue.(string))
	// Import Path
	importPathValue, ok := fields["import"]
	if !ok {
		return nil, errors.New("Required field 'import' not found")
	}
	importPath := importPathValue.(string)
	// Link
	var link string
	switch doctype {
	case PackageKind:
		link = "http://" + path.Join("godoc.org/", importPath)
	case FuncKind, MethodKind, ConstKind, VarKind, TypeKind:
		basepath := "http://" + path.Join("godoc.org/", importPath)
		link = fmt.Sprintf("%s#%s", basepath, name)
	}
	// Highlights - Name
	var highlightName string
	if hname, ok := fragments["name"]; ok {
		highlightName = hname[0]
	}
	// Highlights - Content
	var highlightContent string
	if hcontent, ok := fragments["doc"]; ok {
		highlightContent = hcontent[0]
	}

	return &SearchResult{
		Name: name,
		Type: DocKind(doctype),
		Link: link,
		Highlights: SearchHighlights{
			Name:    template.HTML(highlightName),
			Content: template.HTML(highlightContent),
		},
	}, nil
}
