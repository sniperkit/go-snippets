package routines

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/VojtechVitek/go-trello"
	"github.com/garyburd/go-oauth/oauth"
)

var (
	oauthClient = oauth.Client{
		TemporaryCredentialRequestURI: "https://trello.com/1/OAuthGetRequestToken",
		ResourceOwnerAuthorizationURI: "https://trello.com/1/OAuthAuthorizeToken",
		TokenRequestURI:               "https://trello.com/1/OAuthGetAccessToken",
	}

	token = make(chan *oauth.Credentials)
)

func Scrap(w http.ResponseWriter, r *http.Request) {
	if trelloClient == nil {
		fmt.Fprintln(w, "Need to authorized with trello")
		return
	}
	fmt.Println("HERE 1:")

	// User @trelloClient
	user, err := trelloClient.Member("zaquestion")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user.FullName)

	// @trelloClient Boards
	boards, err := user.Boards()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("HERE 2:", len(boards))

	if len(boards) > 0 {
		board := boards[0]
		fmt.Printf("* %v (%v)\n", board.Name, board.ShortUrl)

		// @trelloClient Board Lists
		lists, err := board.Lists()
		if err != nil {
			log.Fatal(err)
		}

		for _, list := range lists {
			fmt.Println("   - ", list.Name)

			// @trelloClient Board List Cards
			cards, _ := list.Cards()
			for _, card := range cards {
				fmt.Println("      + ", card.Name)
			}
		}
	}
}
