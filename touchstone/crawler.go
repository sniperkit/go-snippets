package gophertv

import (
  "fmt"
  "log"
  "net/http"
  "strings"

  "appengine"
  "appengine/urlfetch"

  "code.google.com/p/google-api-go-client/googleapi/transport"
  "code.google.com/p/google-api-go-client/youtube/v3"
)

// TODO (sunil): Read it from Conf or commandline param
const developerKey = "xxxxx"
const maxResults = 50

func getYoutubeService(c appengine.Context) (*youtube.Service, error) {
  client := &http.Client{
    Transport: &transport.APIKey{
      Key:       developerKey,
      Transport: &urlfetch.Transport{Context: c},
    },
  }
  return youtube.New(client)
}

// searchPlaylists fetches playlists for a given query.
func searchPlaylists(c appengine.Context, query string) ([]string, error) {
  service, err := getYoutubeService(c)
  if err != nil {
    log.Printf("error creating youtube service :: %v", err)
    return nil, err
  }

  nextPageToken := ""
  var playlists []string

  for {
    // Make the API call to YouTube.
    call := service.Search.List("id").
      Q(query).
      Type("playlist").
      PageToken(nextPageToken).
      MaxResults(maxResults)

    response, err := call.Do()
    if err != nil {
      log.Fatalf("Error making search API call: %v", err)
    }

    // Iterate through each item and add it to the correct list.
    for _, item := range response.Items {
      switch item.Id.Kind {
      case "youtube#playlist":
        playlists = append(playlists, item.Id.PlaylistId)
      }
    }

    nextPageToken = response.NextPageToken
    if nextPageToken == "" {
      break
    }
  }
  return playlists, nil
}

func printPlaylists(l []*youtube.SearchResult) {
  for i, item := range l {
    fmt.Printf(`%2d
ID: %v
Title: %v
Description: %v
`, i, item.Id.PlaylistId, item.Snippet.Title, item.Snippet.Description)
  }
  fmt.Printf("\n\n")
}

func fetchPlaylist(c appengine.Context, plID string) ([]*youtube.PlaylistItem, error) {

  service, err := getYoutubeService(c)
  if err != nil {
    log.Printf("error creating youtube service :: %v", err)
    return nil, err
  }
  // Group video, channel, and playlist results in separate lists.
  var plItems []*youtube.PlaylistItem

  nextPageToken := ""

  for {
    // Make the API call to YouTube.
    call := service.PlaylistItems.List("id, snippet, contentDetails, status").
      PlaylistId(plID).
      MaxResults(maxResults).
      PageToken(nextPageToken)

    response, err := call.Do()
    if err != nil {
      log.Fatalf("Error making search API call: %v", err)
    }
    // Iterate through each item and add it to the correct list.
    for _, item := range response.Items {
      switch item.Kind {
      case "youtube#playlistItem":
        plItems = append(plItems, item)
      }
    }

    nextPageToken = response.NextPageToken
    if nextPageToken == "" {
      break
    }
  }
  return plItems, nil
}

func printPlaylistItems(items []*youtube.PlaylistItem) {
  for i, item := range items {
    i++
    log.Printf(`%2d:
    Snippet: %+v
    Status: %+v
    ContentDetails: %v
    ID: %v \n`, i, item.Snippet.ResourceId.VideoId, item.Status, item.ContentDetails, item.Id)
  }
}

func getVideoIDs(c appengine.Context, plID string) ([]string, error) {
  service, err := getYoutubeService(c)
  if err != nil {
    log.Printf("error creating youtube service :: %v", err)
    return nil, err
  }

  // Group video, channel, and playlist results in separate lists.
  var videoIDs []string
  nextPageToken := ""

  for {
    // Make the API call to YouTube.
    call := service.PlaylistItems.List("snippet").
      PlaylistId(plID).
      MaxResults(maxResults).
      PageToken(nextPageToken)

    response, err := call.Do()
    if err != nil {
      log.Fatalf("Error making search API call: %v", err)
    }
    // Iterate through each item and add it to the correct list.
    for _, item := range response.Items {
      switch item.Kind {
      case "youtube#playlistItem":
        videoIDs = append(videoIDs, item.Snippet.ResourceId.VideoId)
      }
    }

    nextPageToken = response.NextPageToken
    if nextPageToken == "" {
      break
    }
  }
  return videoIDs, nil
}

func fetchVideos(c appengine.Context, ids ...string) ([]*youtube.Video, error) {
  service, err := getYoutubeService(c)
  if err != nil {
    log.Printf("error creating youtube service :: %v", err)
    return nil, err
  }

  var videos []*youtube.Video
  idStr := strings.Join(ids, ",")

  log.Printf("video id str: %s", idStr)
  nextPageToken := ""

  for {
    // Make the API call to YouTube.
    var call *youtube.VideosListCall
    part := "id, snippet, contentDetails, statistics"
    if nextPageToken == "" {
      call = service.Videos.List(part).Id(idStr).MaxResults(maxResults)
    } else {
      call = service.Videos.List(part).
        Id(idStr).
        MaxResults(maxResults).
        PageToken(nextPageToken)
    }

    response, err := call.Do()
    if err != nil {
      return nil, err
    }
    // Iterate through each item and add it to the correct list.
    for _, item := range response.Items {
      switch item.Kind {
      case "youtube#video":
        videos = append(videos, item)
      }
    }

    nextPageToken = response.NextPageToken
    if nextPageToken == "" {
      break
    }
  }
  return videos, nil
}

func printVideos(videos ...*youtube.Video) {
  for i, v := range videos {
    log.Printf(`%2d:
    ID: %v
    Snippet: %+v
    Statistics: %+v
    ContentDetails: %v
   `, i, v.Id, v.Snippet, v.Statistics, v.ContentDetails)
  }
}
