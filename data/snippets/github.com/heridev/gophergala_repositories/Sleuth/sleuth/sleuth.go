package sleuth

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
	"github.com/yanatan16/golang-instagram/instagram"
)

type Post struct {
	Link     string `json:"link"`
	Image    string `json:"image"`
	Tags     []string `json:"tags"`
	Author   string `json:"author"`
	Caption  string `json:"caption"`
	Likes    uint `json:"likes"`
	Comments []string `json:"comments"`
}

type Result struct {
	Posts       []*Post `json:"images"`
	SearchesLeft uint `json:"searchsleft"`
}

var api instagram.Api

var batches int

func Init(_batches int) {
	if _batches == 0 {
		_batches = 20
	}

	batches = _batches
}

func GetTimeStamps(interval int64) (timestamps []string) {
	timestamps = make([]string, batches)
	now := time.Now().Unix()
	for i := range timestamps {
		timestamps[i] = strconv.FormatInt(now+interval, 10)
	}
	return timestamps
}

// Search for occurence of tag in expectedTags (minus first element)
func hasTags(media instagram.Media, expectedTags []string) (found bool) {
	expectedTags = expectedTags[1:]
	tagMap := make(map [string]bool, len(expectedTags))
	found = false
	for _, tag := range expectedTags {
		tagMap[tag] = true
	}
	for _, mT := range media.Tags {
		if tagMap[mT] {
			found = true
			break;
		}
	}
	for _, eT := range expectedTags {
		if strings.Contains(parseCaption(media.Caption), eT) {
			found = true
			return
		}
	}
	for _, comment := range media.Comments.Data {
		for _, eT := range expectedTags {
			if strings.Contains(comment.Text, eT) {
				found = true
				break
			}
		}
	}
	return
}

func parseComments(comments *instagram.Comments) (parsedComments []string) {
	parsedComments = make([]string, len(comments.Data))
	for i, c := range comments.Data {
		parsedComments[i] = c.Text
	}
	return parsedComments
}

func parseCaption (caption *instagram.Caption) (parsedCaption string) {
	parsedCaption = ""
	if caption != nil {
		parsedCaption = caption.Text
	}
	return
}

func GetPosts(out chan<- *Post, done chan<- int, api *instagram.Api, tags []string, time string) {
	params := url.Values{}
	params.Set("max_id", time)
	res, err := api.GetTagRecentMedia(tags[0], params)
	matches := 0
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error %v\n", err)
	}
	for _, m := range res.MediasResponse.Medias {
		if m.Type != "image" || len(tags) > 1 && !hasTags(m, tags) {
			continue
		}
		captionText := parseCaption(m.Caption)
		image := &Post{m.Link, m.Images.LowResolution.Url, m.Tags, m.User.Username, captionText, uint(m.Likes.Count), parseComments(m.Comments)}
		matches++
		out <- image
	}

	done <- matches
}

func AuthenticatedSearch(tags []string, accessToken string) (result *Result) {
	api := &instagram.Api{
		AccessToken: accessToken,
	}
	result = search(api, tags)
	return
}

func UnauthenticatedSearch(tags []string, clientId string) (result *Result) {
	api := &instagram.Api{
		ClientId: clientId,
	}
	result = search(api, tags)
	return
}

func search(api *instagram.Api, tags []string) (result *Result) {
	// Get timestamps minutes apart
	timestamps := GetTimeStamps(100)
	maxImages := batches*20
	images := make([]*Post, maxImages)
	out := make(chan *Post, maxImages)
	done := make(chan int, len(timestamps))
	i := 0
	totalImages := 0
	for _, time := range timestamps {
		go GetPosts(out, done, api, tags, time)
	}
	for matches := range done {
		i++
		totalImages += matches
		if i >= batches {
			close(done)
		}
	}
	fmt.Println("done...found %d images", totalImages)
	for i = 0 ; i < totalImages ; i++ {
		images[i] = <-out	
	}
	result = &Result{images[:totalImages], 1000}
	return
}
