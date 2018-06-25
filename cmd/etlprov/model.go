package main

type Entity struct {
	Name         string
	Availability string
	Transmitting bool
	Comment      string
	Truncation   string
	Limit        string
}

type Step struct {
	ID           int
	Description  string
	Entities     []*Entity
	PreviousStep *Step
	Time         string
}

type Tool struct {
	Name    string
	Steps   []*Step
	Usage   string
	Version string
}

type Source struct {
	Name    string
	Steps   []*Step
	Usage   string
	Version string
}

type Person struct {
	Name  string
	Email string
	Role  string
	Steps []*Step
}
