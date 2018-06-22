package helpers

import (
  "html/template"
  "time"
)

var FuncMap = template.FuncMap{
  "DateFormat": func(t time.Time) string {
    return t.Format("2006/01/02")
  },
}
