package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/docker/docker/pkg/parsers"
	"github.com/docker/docker/registry"
	"sourcegraph.com/sourcegraph/srclib/dep"
	"sourcegraph.com/sourcegraph/srclib/unit"
)

func init() {
	_, err := parser.AddCommand("depresolve",
		"resolve a Dockerfile's dependencies",
		"Resolve a Dockerfile's FROM, RUN, etc., dependencies to their repository clone URLs.",
		&depResolveCmd,
	)
	if err != nil {
		log.Fatal(err)
	}
}

type DepResolveCmd struct {
}

var depResolveCmd DepResolveCmd

func (c *DepResolveCmd) Execute(args []string) error {
	var unit *unit.SourceUnit
	if err := json.NewDecoder(os.Stdin).Decode(&unit); err != nil {
		return err
	}
	if err := os.Stdin.Close(); err != nil {
		return err
	}

	res := make([]*dep.Resolution, len(unit.Dependencies))
	for i, rawDep := range unit.Dependencies {
		baseImage, ok := rawDep.(string)
		if !ok {
			return fmt.Errorf("Dockerfile raw dep is not a string base image: %v (%T)", rawDep, rawDep)
		}

		res[i] = &dep.Resolution{Raw: rawDep}

		hostname, name, repo, tag, err := resolveImageRef(baseImage)
		// _, _, _, tag, err := resolveImageRef(baseImage)
		if err == nil {
			// res[i].Target = &dep.ResolvedTarget{
			// 	ToRepoCloneURL:  DockerImageDynRefCloneURL,
			// 	ToUnitType:      DockerfileUnitType,
			// 	ToUnit:          baseImage,
			// 	ToVersionString: tag,
			// }
			res[i].Error = fmt.Sprintf("Don't know how to resolve Docker image hostname %q name %q repo %q tag %q to a VCS repository.", hostname, name, repo, tag)
		} else {
			res[i].Error = err.Error()
		}
	}

	b, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		return err
	}
	if _, err := os.Stdout.Write(b); err != nil {
		return err
	}
	fmt.Println()
	return nil
}

func resolveImageRef(imageRef string) (hostname, name, repo, tag string, err error) {
	repo, tag = parsers.ParseRepositoryTag(imageRef)
	hostname, name, err = registry.ResolveRepositoryName(repo)
	return hostname, name, repo, tag, err
}
