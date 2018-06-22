package models

import "time"

type User struct {
  Id int64 `json:"id"`

  GitHubName        string `json:"gitHubName"`
  GitHubEmail       string `json:"gitHubEmail"`
  GitHubID          string `json:"gitHubId"`
  GitHubAvatarURL   string `json:"gitHubAvatarUrl"`
  GitHubAccessToken string `json:",omitempty"`
  Admin             bool   `json:",omitempty"`

  CreatedAt time.Time `json:"createdAt"`
  UpdatedAt time.Time `json:"updatedAt"`

  Projects []Project `json:"projects"`
}
