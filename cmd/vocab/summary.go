package main

import (
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/spf13/cobra"
)

// Summary generates aggregate statistics for a set of concepts.
type Summary struct {
	domains  map[string]struct{}
	vocabs   map[string]struct{}
	concepts int
}

func (s *Summary) Domains() []string {
	a := make([]string, len(s.domains))
	i := 0

	for k, _ := range s.domains {
		a[i] = k
		i++
	}

	sort.Strings(a)

	return a
}

func (s *Summary) Vocabs() []string {
	a := make([]string, len(s.vocabs))
	i := 0

	for k, _ := range s.vocabs {
		a[i] = k
		i++
	}

	sort.Strings(a)

	return a
}

func (s *Summary) DomainCount() int {
	return len(s.domains)
}

func (s *Summary) VocabCount() int {
	return len(s.vocabs)
}

func (s *Summary) ConceptCount() int {
	return s.concepts
}

func (s *Summary) Index(c *Concept) {
	if _, ok := s.domains[c.DomainID]; !ok {
		s.domains[c.DomainID] = struct{}{}
	}

	if _, ok := s.vocabs[c.VocabularyID]; !ok {
		s.vocabs[c.VocabularyID] = struct{}{}
	}

	s.concepts++
}

func GenerateSummary(r *ConceptReader) (*Summary, error) {
	s := Summary{
		domains: make(map[string]struct{}),
		vocabs:  make(map[string]struct{}),
	}

	for {
		c, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return nil, err
		}

		s.Index(c)
	}

	return &s, nil
}

var summaryCmd = &cobra.Command{
	Use: "summary <path>",

	Short: "Prints a summary of the vocabulary.",

	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			cmd.Usage()
			os.Exit(1)
		}

		f, err := os.Open(args[0])

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		cr, err := NewConceptReader(f)

		if err != nil {
			fmt.Fprintln(os.Stderr, "Error reading file: %s", err)
		}

		cr.Comma = detectDelimiter(args[0])

		fmt.Fprintf(os.Stderr, "Detected %q delimiter\n", cr.Comma)

		s, err := GenerateSummary(cr)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating summary\n%s\n", err)
			os.Exit(1)
		}

		fmt.Println("Counts\n---")
		fmt.Printf("%d\tDomains\n", s.DomainCount())
		fmt.Printf("%d\tVocabularies\n", s.VocabCount())
		fmt.Printf("%d\tConcepts\n", s.ConceptCount())
	},
}
