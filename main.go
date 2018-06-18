package main

import (
	"fmt"
	"github.com/google/go-github/github"
	"io"
	"log"
	"math"
	"os"
	"time"
)

const remaingThreshold = 1

func main() {

	t := &github.UnauthenticatedRateLimitedTransport{
		ClientID:     "SOME_CLIENT_ID",
		ClientSecret: "SOME_CLIENT_SECRET",
	}
	client := github.NewClient(t.Client())

	fmt.Println("Repos that contain magento and PHP code.")

	// create a file to be used for geocoder
	filename := "/tmp/locations.txt"

	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		//this is a fatal error, quit
		return
	}
	defer f.Close()

	// slice the queries into batches to get around the API limit of 1000

	queries := []string{"\"2008-06-01 .. 2012-09-01\"",
		"\"2008-06-01 .. 2012-09-01\"",
		"\"2012-09-02 .. 2013-04-20\"",
		"\"2013-04-21 .. 2013-10-20\"",
		"\"2013-10-21 .. 2014-03-10\"",
		"\"2014-03-10 .. 2014-07-10\"",
		"\"2014-07-10 .. 2014-09-30\""}

	for _, q := range queries {

		query := fmt.Sprintf("magento language:PHP created:" + q)

		page := 1
		maxPage := math.MaxInt32

		opts := &github.SearchOptions{
			Sort:  "updated",
			Order: "desc",
			ListOptions: github.ListOptions{
				PerPage: 100,
			},
		}

		for ; page <= maxPage; page++ {
			opts.Page = page
			result, response, err := client.Search.Repositories(query, opts)
			wait(response)

			if err != nil {
				log.Fatal("FindRepos:", err)
				break
			}

			maxPage = response.LastPage

			msg := fmt.Sprintf("page: %v/%v, size: %v, total: %v",
				page, maxPage, len(result.Repositories), *result.Total)
			log.Println(msg)

			for _, repo := range result.Repositories {

				repoName := *repo.FullName
				username := *repo.Owner.Login
				createdAt := repo.CreatedAt.String()

				fmt.Println("repo: ", repoName)
				fmt.Println("owner: ", username)
				fmt.Println("created at: ", createdAt)

				user, response, err := client.Users.Get(username)
				wait(response)

				if err != nil {
					fmt.Println("error getting userinfo for:", username, err)
					continue
				}

				userLocation := ""
				if user.Location == nil {
					userLocation = "not found"
				} else {
					userLocation = *user.Location
				}

				n, err := io.WriteString(f, "\""+username+"\",\""+userLocation+"\",\""+repoName+"\",\""+createdAt+"\"\n")
				if err != nil {
					fmt.Println(n, err)
				}

				time.Sleep(time.Millisecond * 500)
			}
		}
	}

}

func wait(response *github.Response) {
	if response != nil && response.Remaining <= remaingThreshold {
		gap := time.Duration(response.Reset.Local().Unix() - time.Now().Unix())
		sleep := gap * time.Second
		if sleep < 0 {
			sleep = -sleep
		}

		time.Sleep(sleep)
	}
}
