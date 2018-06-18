package backend

import (
	//"fmt"
	"net/url"
)

func CleanURL(rawurl string) (string, error) {
	//fmt.Println("#################", rawurl)
	// decode the URL if it is already encoded.
	URL, err := url.QueryUnescape(rawurl)
	if err != nil {
		//	fmt.Println("$$$$$$$$$$$$$$$$$$$", err.Error())
		return "", err
	}

	// parse the url
	u, err := url.Parse(URL)
	if err != nil {
		//	fmt.Println("$$$$$$$$$$$$$$$$$$$", err.Error())
		return "", err
	}

	// just get the host name and the rest of the path, ignoring the queries, bookmarks, schema etc
	cleanURL := u.Host + u.Path
	//fmt.Println("#################", u.Host)
	//fmt.Println("#################", u.Path)
	//fmt.Println("#################", cleanURL, len(cleanURL)-1)

	if string(cleanURL[len(cleanURL)-1]) == "/" {

		cleanURL = cleanURL[:len(cleanURL)-1]
	}
	//fmt.Println("#################", cleanURL)
	return cleanURL, nil
}
