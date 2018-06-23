package models

import (
  "time"
)

type EmailList struct {
  Id               int64      `db:"id" json:"id"`
  Title            string     `db:"title" json:"title"`
  Content          string     `sql:"size:0" db:"content" json:"content"`
  CreatedAt        time.Time  `db:"created" json:"created_at"`
  UpdatedAt        time.Time  `db:"updated" json:"updated_at"`
  AutoresponderId  int64      `json:"autoresponder_id"`

}
