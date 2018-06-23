# Golang Project Template

## Introduction

This is the starting point for a Golang project with Gometalinter, Govendor, TDD, Precommit git hooks, security best practices, etc.

## Contents

- [Prerequisities](#prerequisities)

## Prerequisities

```bash
# Govendor
go get -u github.com/kardianos/govendor
govendor sync

# Certs (Self signed for now)
mkdir certs
openssl req -x509 -out certs/server.crt -newkey rsa:4096 -keyout certs/server.key -days 365 -nodes -config openssl.cnf -extensions v3_ext
```

## Usage

```bash
go build
./golang-project-template serve
```
