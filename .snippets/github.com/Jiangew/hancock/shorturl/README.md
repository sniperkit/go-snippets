# Shorturl Service
This sample is a API application based on beego. It has two API func:

- /v1/shorten
- /v1/expand

## Build
```
cd $GOPATH/src/github.com/jiangew/hancock/shorturl
go build
```

## Usage
```
# shortening url example
http://localhost:8080/v1/shorten/?longurl=http://google.com

{
  "UrlShort": "5laZG",
  "UrlLong": "http://google.com"
}

# expanding url example
http://localhost:8080/v1/expand/?shorturl=5laZI

{
  "UrlShort": "5laZG",
  "UrlLong": "http://google.com"
}
```
