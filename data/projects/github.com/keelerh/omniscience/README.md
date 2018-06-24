[![Build Status](https://travis-ci.org/keelerh/omniscience.svg?branch=master)](https://travis-ci.org/keelerh/omniscience)
[![Coverage Status](https://coveralls.io/repos/github/keelerh/omniscience/badge.svg?branch=master)](https://coveralls.io/github/keelerh/omniscience?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/keelerh/omniscience)](https://goreportcard.com/report/github.com/keelerh/omniscience)

<div align="center">
  <img src="images/logo.png" width="200" height="150">
</div>

Omniscience is a search utility for searching over private documents aggregated from a collection of sources.

```
+----------+
|Google    |
|Drive     | --+
|Service   |   |
+----------+   |
               |
+----------+   |    +----------+     +----------+     +----------+                       
|GitHub    |   |    |Ingestion |     |DB        |     |Frontend  |                   
|Service   | --+--> |Service   | --> |          | --> |          |
|          |   |    |          |     |          |     |          |
+----------+   |    +----------+     +----------+     +----------+
               |
+----------+   |
|Confluence|   |
|Service   | --+
|          |
+----------+
```


## Elasticsearch

### Install

```
$ brew update
$ brew install elasticsearch
```

### Using local ES server

To avoid getting a cors related error when connecting the UI to your local elasticsearch instance, add the following to your `config/elasticsearch.yml` file (if you installed Elasticsearch using brew, this file can be found at `/usr/local/Cellar/elasticsearch/6.2.2/libexec/config/elasticsearch.yml`).

```
http.cors.enabled : true  
http.cors.allow-origin : "*"
http.cors.allow-methods : OPTIONS, HEAD, GET, POST, PUT, DELETE
http.cors.allow-headers : X-Requested-With,X-Auth-Token,Content-Type, Content-Length
```

### Start local ES server

```
$ elasticsearch
```

## Golang backend

### Start Go server

To start the omniscience server, run

```
$ go run cmd/omniscience_server/main.go --fGoogleServiceAccountPath=<path-to-service-account-json>
```

from the root of the project.

### Manually invoke indexing

Indexing is invoked on a per service basis. To index files from a given service navigate to the service's directory in the `cmd/omniscience_client` directory and run 

```
$ go run main.go
```

with the relevant flags for the service. More detailed instructions for running each service can be found in the README of each service's directory.

Elasticsearch and the omniscience server must be running for indexing to work.

## React app

### Install dependencies

From the project root run:

```
$ cd client && npm i
```

### Start local node server

From the project root run:

```
$ cd client && npm start
```

## Search everything!

Navigate to [http://localhost:3000](http://localhost:3000)

## License
This project is licensed under the MIT License - see the LICENSE.md file for details