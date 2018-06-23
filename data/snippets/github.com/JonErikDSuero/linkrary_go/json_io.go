package main

import (
  "net/http"
  "io/ioutil"
  "log"
  "encoding/json"
  "strconv"
)

func readJson(r *http.Request, v interface{}) bool {
  defer r.Body.Close()
  var (
    body []byte
    err  error
  )
  body, err = ioutil.ReadAll(r.Body)
  if err != nil {
    log.Printf("ReadJson couldn't read request body %v", err)
    return false
  }
  if err = json.Unmarshal(body, v); err != nil {
    log.Printf("ReadJson couldn't parse request body %v", err)
    return false
  }
  return true
}


func writeJson(w http.ResponseWriter, v interface{}) {
  // avoid json vulnerabilities, always wrap v in an object literal
  doc := map[string]interface{}{"d": v}

  if data, err := json.Marshal(doc); err != nil {
    log.Printf("Error marshalling json: %v", err)
  } else {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Content-Length", strconv.Itoa(len(data)))
    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
  }
}
