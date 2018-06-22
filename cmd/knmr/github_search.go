package main

import (
	"context"
	"fmt"

	"github.com/Ken2mer/knmr/logger"
	"github.com/google/go-github/github"
)

func getCodeSearchResult(ctx context.Context, client *github.Client) (*github.CodeSearchResult, error) {
	result, _, err := client.Search.Code(ctx, username, nil)
	if err != nil {
		return nil, err
	}
	return result, err
}

func dumpCodeSearchResult(code *github.CodeSearchResult) {
	for i, result := range code.CodeResults {
		repo := result.GetRepository().GetFullName()
		path := result.GetPath()
		url := result.GetHTMLURL()

		logger.Debugf("index: %d\n", i)
		fmt.Printf("github.com/%s/%s\n", repo, path)
		fmt.Printf("%s\n", url)
	}
}
