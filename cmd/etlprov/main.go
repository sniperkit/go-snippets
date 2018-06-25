package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	dms "github.com/chop-dbhi/data-models-service/client"
)

var pedsnetIgnores = []string{
	"concept",
	"concept_ancestor",
	"concept_class",
	"concept_relationship",
	"concept_synonym",
	"domain",
	"drug_strength",
	"relationship",
	"source_to_concept_map",
	"vocabulary",
}

func main() {
	var (
		dir      string
		model    string
		version  string
		service  string
		ignore   string
		truncate bool
		ver      bool
	)

	flag.StringVar(&model, "model", "pedsnet", "Name of the data model to validate against.")
	flag.StringVar(&version, "version", "", "Version of the data model to validate against.")
	flag.StringVar(&service, "service", "http://data-models.origins.link", "URL to the data models service.")
	flag.StringVar(&ignore, "ignore", "", "Comma-separated list of entities to ignore.")
	flag.BoolVar(&truncate, "truncate", true, "Truncate the list of errors.")
	flag.BoolVar(&ver, "v", false, "Prints the version.")

	flag.Parse()

	if ver {
		fmt.Println(progVersion)
		return
	}

	args := flag.Args()

	if len(args) == 0 {
		dir = "."
	} else {
		dir = args[0]
	}

	// Ensure the directory exists. This is performed here to save a remote call.
	stat, err := os.Stat(dir)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if !stat.IsDir() {
		fmt.Fprintf(os.Stderr, "'%s' not a directory\n", stat.Name())
		os.Exit(1)
	}

	client, err := dms.New(service)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err = client.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "Error communicating with service\n> %s\n", err)
		os.Exit(1)
	}

	dm, err := client.ModelRevision(model, version)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Problem fetching model data\n> %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Validating against model '%s/%s'\n", model, version)
	fmt.Printf("Scanning files in '%s'\n", dir)

	p := NewParser(dm)

	var limit int

	if truncate {
		limit = 10
	}

	errPrinter := ErrorPrinter{
		Limit:  limit,
		Writer: os.Stderr,
	}

	if err = p.ReadDir(dir, errPrinter.Handle); err != nil {
		fmt.Fprintf(os.Stderr, "Problem parsing files in directory '%s'\n> %s\n", dir, err)
		os.Exit(1)
	}

	// Determine missing entities.
	var ignores []string

	if ignore != "" {
		ignores = strings.Split(ignore, ",")
	} else {
		// Special case defaults.
		if dm.Name == "pedsnet" {
			ignores = pedsnetIgnores
		}
	}

	names := p.MissingEntities(true, ignores)

	if len(names) > 0 {
		fmt.Println("---")

		if len(names) == 1 {
			fmt.Fprintln(os.Stderr, "1 entity is missing from the model")
		} else {
			fmt.Fprintf(os.Stderr, "%d entities are missing from the model\n", len(names))
		}

		sort.Strings(names)
		errs := make([]error, len(names))

		for i, n := range names {
			errs[i] = fmt.Errorf(n)
		}

		printErrors(os.Stderr, errs, limit)
	}

	entities := p.EntitiesWithoutSteps()

	if len(entities) > 0 {
		fmt.Println("---")

		if len(entities) == 1 {
			fmt.Fprintln(os.Stderr, "1 entity does not have steps")
		} else {
			fmt.Fprintf(os.Stderr, "%d entities do not have steps\n", len(entities))
		}

		errs := make([]error, len(entities))

		for i, e := range entities {
			errs[i] = fmt.Errorf(e.Name)
		}

		printErrors(os.Stderr, errs, limit)
	}

	fmt.Println("---")

	fmt.Printf("%d entities\n", len(p.Entities()))
	fmt.Printf("%d steps\n", len(p.Steps()))
	fmt.Printf("%d tools\n", len(p.Tools()))
	fmt.Printf("%d sources\n", len(p.Sources()))
	fmt.Printf("%d persons\n", len(p.People()))
}
