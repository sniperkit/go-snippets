package main

import "fmt"
import "github.com/chuckpreslar/codex"

func main() {
  users := codex.Table("test")
  sql, err := users.ToSql()
  if err != nil {
    fmt.Println("ToSql error :", err)
  }
  fmt.Println(sql)
}
