package models

type User struct {
  Avatar        string `json:"avatar_url"`
  Login         string `json:"login"`
  Name          string `json:"name"`
  DashboardStar int    `json:"dashboard_star"`
  Repos         int    `json:"public_repos"`
  Followers     int    `json:"followers"`
  Location      string `json:"location"`
  UpdatedAt     string `json:"updated_at"`

  Url           string `json:"url"`
  ReposUrl      string `json:"repos_url"`
}

type UserByDashboard []*User

func (users UserByDashboard) Len() int {
  return len(users)
}

func (users UserByDashboard) Swap(i, j int) {
  users[i], users[j] = users[j], users[i]
}

func (users UserByDashboard) Less(i, j int) bool {
  return users[i].DashboardStar < users[j].DashboardStar
}

type PagedUser struct {
  Pagination
  Items []*User `json:"items"`
}

type Repo struct {
  FullName        string `json:"full_name"`
  StargazersCount int    `json:"stargazers_count"`
}

type ReposByStar []*Repo

func (repos ReposByStar) Len() int {
  return len(repos)
}

func (repos ReposByStar) Swap(i, j int) {
  repos[i], repos[j] = repos[j], repos[i]
}

func (repos ReposByStar) Less(i, j int) bool {
  return repos[i].StargazersCount < repos[j].StargazersCount
}

type Pagination struct {
  TotalCount        int  `json:"total_count"`
  IncompleteResults bool `json:"incomplete_results"`
}

type Ranks struct {
  UpdatedAt string  `json:"updated_at"`
  Ranks     []*User `json:"ranks"`
}

type PutRet struct {
  Hash string `json:"hash"`
  Key  string `json:"key"`
}
