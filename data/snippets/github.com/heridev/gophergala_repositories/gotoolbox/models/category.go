package models

import "time"

type Category struct {
  Id int64 `json:"id"`

  Name        string `json:"name",schema:"name"`
  Description string `json:"description",schema:"description"`

  CreatedAt time.Time `json:"createdAt",schema:"createdAt"`
  UpdatedAt time.Time `json:"updatedAt",schema:"updatedAt"`

  Projects []Project `json:"projects"`
}
