# fluxio: a tool for JSON/PostgreSQL-jsonb I/Os

[![Build Status](https://travis-ci.org/DevMine/fluxio.png?branch=master)](https://travis-ci.org/DevMine/fluxio)
[![GoDoc](http://godoc.org/github.com/DevMine/fluxio?status.svg)](http://godoc.org/github.com/DevMine/fluxio)
[![GoWalker](http://img.shields.io/badge/doc-gowalker-blue.svg?style=flat)](https://gowalker.org/github.com/DevMine/fluxio)

`fluxio` is a command line tool to insert JSON data from its standard input into
a [PostgreSQL](http://www.postgresql.org/) 9.4+ jsonb store and which can also
output JSON from the database to its standard output.

`fluxio` needs to be provided a configuration file from which it can read
database connexion information. You can simply rename `fluxio.conf.sample` to
`fluxio.conf` and adjust options where necessary.

`fluxio` expects one table with two columns; one of type `character varying`
which acts as the primary key (default name: key) and another one of type
`jsonb` which acts as the JSON store (default name: content). You can specify
the column names with the `-col-content` and `-col-key` flags. If you do not use
the default PostgreSQL schema, you can specify the schema name with the `-s`
option.

`fluxio` expects to be provided a key in order to insert data into the database
or to extract content from it. The key is provided using the `-k` flag.

`fluxio` can either import JSON into the specified table with the `-w` flag or
read from the table and export JSON to standard output using the `-r` flag.

To install `fluxio`, simply issue the following command:

    go get -u github.com/DevMine/fluxio
