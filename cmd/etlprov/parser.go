package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	dms "github.com/chop-dbhi/data-models-service/client"
)

const (
	TokenDelim  = ","
	EntityDelim = "."

	EntitiesFile = "entities.csv"
	StepsFile    = "steps.csv"
	PeopleFile   = "people.csv"
	ToolsFile    = "tools.csv"
	SourcesFile  = "sources.csv"
)

// Regex for splitting entity list strings.
var entitySplitter = regexp.MustCompile(`[,\s+]`)

// Regex for a step range, taking the form "N-M".
var stepRange = regexp.MustCompile(`(\d+)\s*-\s*(\d+)`)

// Parses the availability string ensuring it is valid.
func parseAvailability(v string) (string, error) {
	x := strings.ToLower(v)

	switch x {
	case "available", "unavailable", "unknown":
		return x, nil
	}

	return "", fmt.Errorf("Invalid choice for availability: %s", v)
}

// Parses the transmitting boolean ensuring it is valid.
func parseTransmitting(v string) (bool, error) {
	x := strings.ToLower(v)

	switch x {
	case "", "yes", "1", "true":
		return true, nil
	case "no", "0", "false":
		return false, nil
	}

	return false, fmt.Errorf("Invalid choice for transmitting: %s", v)
}

type Parser struct {
	Model *dms.Model

	entities map[string]*Entity
	steps    map[int]*Step
	tools    map[string]*Tool
	people   map[string]*Person
	sources  map[string]*Source
}

func (p *Parser) parseTool(record []string) (*Tool, error) {
	if record[0] == "" {
		return nil, fmt.Errorf("Tool name required")
	}

	t := &Tool{
		Name:    record[0],
		Usage:   record[2],
		Version: record[3],
	}

	steps, err := p.parseStepString(record[1])

	if err != nil {
		return t, err
	}

	t.Steps = steps

	return t, nil
}

// Parses a step value which could be a single step, list of steps or a range.
func (p *Parser) parseStepString(v string) ([]*Step, error) {
	var steps []*Step

	toks := strings.Split(v, TokenDelim)

	for _, t := range toks {
		t = strings.TrimSpace(t)

		if match := stepRange.FindStringSubmatch(t); match != nil {
			// Ignore errors since regexp matched on them
			l, _ := strconv.Atoi(match[1])
			u, _ := strconv.Atoi(match[2])

			if u < l {
				return nil, fmt.Errorf("Invalid step range %s", t)
			}

			// Ensure the upper and lower bound explicitly match in the
			// available steps.
			var lm, um bool

			// Get all steps within the range, inclusive.
			for id, s := range p.steps {
				if id >= l && id <= u {
					if id == l {
						lm = true
					}

					if id == u {
						um = true
					}

					steps = append(steps, s)
				}
			}

			if !lm {
				return nil, fmt.Errorf("Low step in range not defined %d", l)
			}

			if !um {
				return nil, fmt.Errorf("High step in range not defined %d", u)
			}
		} else if t != "" {
			i, err := strconv.Atoi(t)

			if err != nil {
				return nil, fmt.Errorf("Invalid step number %s", t)
			}

			s, ok := p.steps[i]

			if !ok {
				return nil, fmt.Errorf("Step %d not defined", i)
			}

			steps = append(steps, s)
		}
	}

	return steps, nil
}

func (p *Parser) parseStep(record []string) (*Step, error) {
	var (
		id, pid int
		ok      bool
		prev    *Step
		err     error
	)

	if record[0] == "" {
		return nil, fmt.Errorf("Step ID required")
	}

	// Try to convert into a number.
	id, err = strconv.Atoi(record[0])

	if err != nil {
		return nil, err
	}

	s := &Step{
		ID:          id,
		Time:        record[4],
		Entities:    nil,
		Description: record[1],
	}

	// Parse the entities the step applies to.
	entities, err := p.parseEntityString(record[2])

	if err != nil {
		return s, err
	}

	s.Entities = entities

	// Special case of returning the step wit the error to
	// prevent cascading errors.
	if record[3] != "" {
		// Parse the previous step if specified.
		if pid, err = strconv.Atoi(record[3]); err != nil {
			err = fmt.Errorf("Error parsing previous step '%s'", record[3])
		} else if prev, ok = p.steps[pid]; !ok {
			err = fmt.Errorf("Step %d does not exist", pid)
		}

		s.PreviousStep = prev
	}

	return s, err
}

// parseSource parses a source record.
func (p *Parser) parseSource(record []string) (*Source, error) {
	if record[0] == "" {
		return nil, fmt.Errorf("Source name required")
	}

	s := &Source{
		Name:    record[0],
		Usage:   record[2],
		Version: record[3],
	}

	steps, err := p.parseStepString(record[1])

	if err != nil {
		return s, err
	}

	s.Steps = steps

	return s, nil
}

// parsePerson parses a person record.
func (p *Parser) parsePerson(record []string) (*Person, error) {
	if record[0] == "" {
		return nil, fmt.Errorf("Person name required")
	}

	b := &Person{
		Name:  record[0],
		Email: record[1],
		Role:  record[2],
	}

	steps, err := p.parseStepString(record[3])

	if err != nil {
		return b, err
	}

	b.Steps = steps

	return b, nil
}

func (p *Parser) parseEntity(record []string) (*Entity, error) {
	name, err := p.validateEntityName(record[0])

	if err != nil {
		return nil, err
	}

	e := &Entity{
		Name:       name,
		Comment:    record[3],
		Truncation: record[4],
		Limit:      record[5],
	}

	avail, err := parseAvailability(record[1])

	if err != nil {
		return e, err
	}

	e.Availability = avail

	trans, err := parseTransmitting(record[2])

	if err != nil {
		return e, err
	}

	e.Transmitting = trans

	return e, nil
}

func (p *Parser) validateEntityName(name string) (string, error) {
	toks := strings.SplitN(strings.ToLower(name), EntityDelim, 2)

	table := p.Model.Tables.Get(toks[0])

	if table == nil {
		if len(toks) > 1 {
			return "", fmt.Errorf("Unknown table '%s' for field '%s'", toks[0], toks[1])
		}
		return "", fmt.Errorf("Unknown table '%s'", toks[0])
	}

	if len(toks) == 1 {
		return name, nil
	}

	field := table.Fields.Get(toks[1])

	if field == nil {
		return "", fmt.Errorf("Unknown field '%s'", name)
	}

	return name, nil
}

// parseEntityString parses a string and validates it against an entity parser
// that is populated with entities.
func (p *Parser) parseEntityString(s string) ([]*Entity, error) {
	var (
		err      error
		entities []*Entity
	)

	for _, name := range entitySplitter.Split(strings.ToLower(s), -1) {
		if name == "all" || name == "" {
			var i int
			entities = make([]*Entity, len(p.entities))

			for _, e := range p.entities {
				entities[i] = e
				i++
			}

			return entities, nil
		}

		name, err = p.validateEntityName(name)

		if err != nil {
			return nil, err
		}

		// Table; add all fields for the table verifying they were defined in
		// the entities file.
		if !strings.Contains(name, EntityDelim) {
			t := p.Model.Tables.Get(name)

			// Include the table itself as an entity if defined in the entity list.
			if e, ok := p.entities[name]; ok {
				entities = append(entities, e)
			} else {
				return nil, fmt.Errorf("Entity '%s' not defined", name)
			}

			for _, f := range t.Fields.List() {
				n := fmt.Sprintf("%s%s%s", name, EntityDelim, f.Name)

				if e, ok := p.entities[n]; ok {
					entities = append(entities, e)
				} else {
					return nil, fmt.Errorf("Entity '%s' not defined", n)
				}
			}

			continue
		}

		if e, ok := p.entities[name]; ok {
			entities = append(entities, e)
		} else {
			return nil, fmt.Errorf("Entity '%s' not defined", name)
		}
	}

	return entities, nil
}

func (p *Parser) ReadEntities(r io.Reader) []error {
	var (
		e    *Entity
		err  error
		errs []error
	)

	err = ReadRows(r, func(row []string) {
		e, err = p.parseEntity(row)

		if err != nil {
			errs = append(errs, err)
		}

		if e != nil {
			if _, ok := p.entities[e.Name]; ok {
				err = fmt.Errorf("Duplicate entity '%s' found", e.Name)
				errs = append(errs, err)
				return
			}

			p.entities[e.Name] = e
		}
	})

	// Error returned while reading.
	if err != nil {
		return []error{err}
	}

	return errs
}

func (p *Parser) ReadSteps(r io.Reader) []error {
	var (
		s    *Step
		err  error
		errs []error
	)

	err = ReadRows(r, func(row []string) {
		s, err = p.parseStep(row)

		// Special case to log steps to prevent cascading errors.
		if err != nil {
			errs = append(errs, err)
		}

		if s != nil {
			if _, ok := p.steps[s.ID]; ok {
				err = fmt.Errorf("Duplicate step '%d', found", s.ID)
				errs = append(errs, err)
				return
			}

			p.steps[s.ID] = s
		}
	})

	// Error returned while reading.
	if err != nil {
		return []error{err}
	}

	return errs
}

func (p *Parser) ReadTools(r io.Reader) []error {
	var (
		t    *Tool
		err  error
		errs []error
	)

	err = ReadRows(r, func(row []string) {
		t, err = p.parseTool(row)

		if err != nil {
			errs = append(errs, err)
		}

		if t != nil {
			if _, ok := p.tools[t.Name]; ok {
				err = fmt.Errorf("Duplicate tool '%s', found", t.Name)
				errs = append(errs, err)
				return
			}

			p.tools[t.Name] = t
		}
	})

	// Error returned while reading.
	if err != nil {
		return []error{err}
	}

	return errs
}

func (p *Parser) ReadSources(r io.Reader) []error {
	var (
		s    *Source
		err  error
		errs []error
	)

	err = ReadRows(r, func(row []string) {
		s, err = p.parseSource(row)

		if err != nil {
			errs = append(errs, err)
			return
		}

		if s != nil {
			if _, ok := p.sources[s.Name]; ok {
				err = fmt.Errorf("Duplicate source '%s', found", s.Name)
				errs = append(errs, err)
				return
			}

			p.sources[s.Name] = s
		}
	})

	// Error returned while reading.
	if err != nil {
		return []error{err}
	}

	return errs
}

func (p *Parser) ReadPeople(r io.Reader) []error {
	var (
		v    *Person
		err  error
		errs []error
	)

	err = ReadRows(r, func(row []string) {
		v, err = p.parsePerson(row)

		if err != nil {
			errs = append(errs, err)
			return
		}

		if v != nil {
			if _, ok := p.people[v.Name]; ok {
				err = fmt.Errorf("Duplicate step '%v', found", v.Name)
				errs = append(errs, err)
				return
			}

			p.people[v.Name] = v
		}
	})

	// Error returned while reading.
	if err != nil {
		return []error{err}
	}

	return errs
}

// EntitiesWithoutSteps returns a list of entities without steps associated with them.
func (p *Parser) EntitiesWithoutSteps() []*Entity {
	// Copy the map.
	m := make(map[string]*Entity, len(p.entities))

	for n, e := range p.entities {
		// Only steps that are being transmitted need to have steps
		// associated with them.
		if e.Transmitting {
			m[n] = e
		}
	}

	for _, s := range p.steps {
		for _, e := range s.Entities {
			if _, ok := m[e.Name]; ok {
				delete(m, e.Name)
			}
		}
	}

	i := 0
	o := make([]*Entity, len(m))

	for _, e := range m {
		o[i] = e
		i++
	}

	return o
}

// MissingEntities determines the entities defined in the model, but that
// have not be evaluated by the parser.
func (p *Parser) MissingEntities(req bool, ignores []string) []string {
	var (
		ok    bool
		name  string
		names []string
	)

	igidx := make(map[string]struct{}, len(ignores))

	for _, name := range ignores {
		igidx[name] = struct{}{}
	}

	for _, t := range p.Model.Tables.List() {
		if _, ok = igidx[t.Name]; ok {
			continue
		}

		for _, f := range t.Fields.List() {
			name = fmt.Sprintf("%s%s%s", t.Name, EntityDelim, f.Name)

			if _, ok = igidx[name]; ok {
				continue
			}

			if _, ok = p.entities[name]; !ok && (req && f.Required) {
				names = append(names, name)
			}
		}
	}

	return names
}

// Entities returns the entities that have been parsed.
func (p *Parser) Entities() []*Entity {
	i := 0
	a := make([]*Entity, len(p.entities))

	for _, e := range p.entities {
		a[i] = e
		i++
	}

	return a
}

// Steps returns the steps that have been parsed.
func (p *Parser) Steps() []*Step {
	i := 0
	a := make([]*Step, len(p.steps))

	for _, s := range p.steps {
		a[i] = s
		i++
	}

	return a
}

// Sources returns the sources that have been parsed.
func (p *Parser) Sources() []*Source {
	i := 0
	a := make([]*Source, len(p.sources))

	for _, s := range p.sources {
		a[i] = s
		i++
	}

	return a
}

// Tools returns the tools that have been parsed.
func (p *Parser) Tools() []*Tool {
	i := 0
	a := make([]*Tool, len(p.tools))

	for _, s := range p.tools {
		a[i] = s
		i++
	}

	return a
}

// People returns the people that have been parsed.
func (p *Parser) People() []*Person {
	i := 0
	a := make([]*Person, len(p.people))

	for _, s := range p.people {
		a[i] = s
		i++
	}

	return a
}

func (p *Parser) ReadDir(dir string, handle ErrorHandler) error {
	var (
		name string
		file *os.File
		err  error
		errs []error
	)

	// Entities.
	name = filepath.Join(dir, EntitiesFile)
	file, err = os.Open(name)

	if err != nil {
		return err
	}

	errs = p.ReadEntities(file)

	file.Close()
	handle(EntitiesFile, errs)

	// Steps.
	name = filepath.Join(dir, StepsFile)
	file, err = os.Open(name)

	if err != nil {
		return err
	}

	errs = p.ReadSteps(file)

	file.Close()
	handle(StepsFile, errs)

	// Tools.
	name = filepath.Join(dir, ToolsFile)
	file, err = os.Open(name)

	if err != nil {
		return err
	}

	errs = p.ReadTools(file)

	file.Close()
	handle(ToolsFile, errs)

	// Sources.
	name = filepath.Join(dir, SourcesFile)
	file, err = os.Open(name)

	if err != nil {
		return err
	}

	errs = p.ReadSources(file)

	file.Close()
	handle(SourcesFile, errs)

	// People.
	name = filepath.Join(dir, PeopleFile)
	file, err = os.Open(name)

	if err != nil {
		return err
	}

	errs = p.ReadPeople(file)
	file.Close()

	handle(PeopleFile, errs)

	return nil
}

// NewParser initializes a new parser.
func NewParser(model *dms.Model) *Parser {
	return &Parser{
		Model:    model,
		entities: make(map[string]*Entity),
		steps:    make(map[int]*Step),
		tools:    make(map[string]*Tool),
		people:   make(map[string]*Person),
		sources:  make(map[string]*Source),
	}
}

type ErrorHandler func(name string, errs []error)

type ErrorPrinter struct {
	Limit  int
	Writer io.Writer
}

// Print implements the ErrorHandler
func (p *ErrorPrinter) Handle(name string, errs []error) {
	if len(errs) > 0 {
		fmt.Fprintln(p.Writer, "---")

		if len(errs) == 1 {
			fmt.Fprintf(p.Writer, "1 error has been detected for '%s'\n", name)
		} else {
			fmt.Fprintf(p.Writer, "%d errors have been detected for '%s'\n", len(errs), name)
		}

		printErrors(p.Writer, errs, p.Limit)
	}
}

func printErrors(w io.Writer, errs []error, trunc int) {
	for i, err := range errs {
		if trunc > 0 && i == trunc {
			fmt.Fprintf(w, "[truncated %d errors]\n", len(errs)-trunc)
			break
		}

		fmt.Fprintln(w, "*", err)
	}
}
