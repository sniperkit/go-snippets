package models

import (
  "time"
)

type Subscriber struct {
  Id         int64        `db:"id" json:"id"`
  CreatedAt  time.Time    `db:"created" json:"created_at"`
  UpdatedAt  time.Time    `db:"updated" json:"updated_at"`
  Name       string       `db:"name" json:"name"`
  Email      string       `db:"email" json:"email"`
}
