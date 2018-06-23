package main

import (
    "encoding/json"
    "io/ioutil"
    "flag"
    "fmt"
    "log"
    "net/http"
    )

func main() {
    var f = flag.String("f", "", "File to parse for generation")
    var t = flag.String("t", "", "JSON String to parse for generation")
    var s = flag.Bool("s", false, "Run web service (set for true)")
    var port = flag.Int("p", 8080, "Port for web service")
    flag.Parse()
    
    var p interface{}
    var err error
    
    if *f != "" {
        j, err := ioutil.ReadFile(*f)
        if err != nil {
            log.Fatal(err)
        }
        
        err = json.Unmarshal(j, &p)
        if err != nil {
            log.Fatal(err)
        }

        o, err := json.Marshal(p)
        if err != nil {
            log.Fatal(err)
        }

        fmt.Println(string(o))
    }
    
    if *t != "" {
        err = json.Unmarshal([]byte(*t), &p)
        if err != nil { 
            log.Fatal(err)
        }
        
        o, err := json.Marshal(p)
        if err != nil {
            log.Fatal(err)
        }
        
        fmt.Println(string(o))
    }
    
    if *s {
        fs := http.FileServer(http.Dir("."))
        http.Handle("/", fs)
        
        fmt.Println(fmt.Sprintf("Listening on port: %d...\n", *port))
        err = http.ListenAndServe(fmt.Sprintf(":%d", *port), nil)
        if err != nil {
            log.Fatal(err)
        }
    }
    
}
