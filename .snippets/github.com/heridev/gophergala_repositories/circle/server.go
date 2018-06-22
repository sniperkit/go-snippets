package main

import (
    "os"
    "log"
    "net/http"
    "html/template"
    "github.com/gorilla/mux"
    "github.com/mitchellh/go-homedir"
    
)

type FList struct {
    Files []FullPath
}

type FullPath struct {
    Name string
    Path string
}

func Run(){
    
    
    r := mux.NewRouter()
    
    r.HandleFunc("/", DefaultHandler).Methods("GET")
    r.HandleFunc("/", PathHandler).Methods("POST")
    r.HandleFunc("/sort", SortHandler).Methods("POST")
    r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("./public/"))))
    
    http.Handle("/", r)
    http.ListenAndServe(":8080", nil)
}

func DefaultHandler(res http.ResponseWriter, req *http.Request){
    path, err := homedir.Dir()
    log.Println(path)
    checkErr(err)
    
    fobj := CreateFList(path)
    
    t, err := template.ParseFiles("index.html")
    checkErr(err)
    
    err = t.Execute(res, fobj)
    checkErr(err)
}

func PathHandler(res http.ResponseWriter, req *http.Request){
    path := req.FormValue("path")
    log.Println(path)
    fobj := CreateFList(path)
    
    t, err := template.ParseFiles("index.html")
    checkErr(err)
    
    err = t.Execute(res, fobj)
    checkErr(err)
}

type SortObj struct {
    Path string
}

func SortHandler(res http.ResponseWriter, req *http.Request){
    path := req.FormValue("path")
    Sort(path)
    
    t, err := template.ParseFiles("sortpage.html")
    checkErr(err)
    
    err = t.Execute(res, &SortObj{Path: path})
    checkErr(err)
}

func checkErr(err error){
    if err != nil {
        log.Fatal(err)
    }
}

func CreateFList(path string) *FList {
    fobj := &FList{Files: make([]FullPath, 0, 100)}
    
    
    dir, err := os.Open(path)
    checkErr(err)
    
    
    fi, err := dir.Readdir(100)
    //checkErr(err)
    
    if(len(fi)==0){    
    fobj.Files = append(fobj.Files, FullPath{Name:"This Directory is empty", Path: path})
    return fobj
    }
    
    for _, file := range fi {
        if file.IsDir() {
            fobj.Files = append(fobj.Files, 
                                FullPath{Name:file.Name(), Path: path+"/"+file.Name()})
            log.Println(file.Name())
        }
    }
    
    //fobj := &FList{Files: fi}
    
    return fobj
}
    
    