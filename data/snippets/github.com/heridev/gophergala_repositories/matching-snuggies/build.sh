#!/bin/bash

go get ./cmd/...
go build -o bin/snuggied ./cmd/snuggied
go build -o bin/snuggier ./cmd/snuggier
