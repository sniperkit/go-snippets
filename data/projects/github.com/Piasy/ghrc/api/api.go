package api

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "math"
  "net/http"
  "sort"
  "time"

  "github.com/Piasy/ghrc/models"
)

const (
  API_ERROR = iota
)

func GetUsers(page int, token string) ([]*models.User, int) {
  resp, err := http.Get(fmt.Sprintf("https://api.github.com/search/users?q=location:china&page=%d&access_token=%s", page, token))
  if err != nil {
    return nil, API_ERROR
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, API_ERROR
  }

  var pagedUser models.PagedUser
  err = json.Unmarshal(body, &pagedUser)
  if err != nil {
    return nil, API_ERROR
  }

  for _, user := range pagedUser.Items {
    fillUserInfo(user, token)
    fmt.Println("fill", user.Login, "finished,", user.DashboardStar, "DashboardStar")
  }

  return pagedUser.Items, 0
}

func fillUserInfo(user *models.User, token string) int {
  resp, err := http.Get(user.Url + "?access_token=" + token)
  if err != nil {
    return API_ERROR
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return API_ERROR
  }

  var fullUser models.User
  err = json.Unmarshal(body, &fullUser)
  if err != nil {
    return API_ERROR
  }

  user.Name = fullUser.Name
  user.Location = fullUser.Location
  user.Followers = fullUser.Followers
  user.UpdatedAt = fullUser.UpdatedAt
  user.Repos = fullUser.Repos
  getDashboardStar(user, token)

  return 0
}

func getDashboardStar(user *models.User, token string) int {
  reposNum := user.Repos
  page := int(math.Ceil(float64(reposNum) / 30))
  repos := make([]*models.Repo, 0, reposNum)
  for i := 1; i <= page; i++ {
    onePage, err := listReposOf(user, i, token)
    if err != 0 {
      return API_ERROR
    }
    repos = append(repos, onePage...)
    time.Sleep(time.Second)
  }

  sort.Sort(sort.Reverse(models.ReposByStar(repos)))
  for i, repo := range repos {
    if i >= 5 {
      break
    }
    user.DashboardStar += repo.StargazersCount
  }

  return 0
}

func listReposOf(user *models.User, page int, token string) ([]*models.Repo, int) {
  resp, err := http.Get(fmt.Sprintf("https://api.github.com/users/%s/repos?page=%d&access_token=%s", user.Login, page, token))
  if err != nil {
    return nil, API_ERROR
  }

  defer resp.Body.Close()
  body, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, API_ERROR
  }

  var repos []*models.Repo
  err = json.Unmarshal(body, &repos)
  if err != nil {
    return nil, API_ERROR
  }
  return repos, 0
}
