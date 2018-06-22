# featscomp: compute DevMine features

[![Build Status](https://travis-ci.org/DevMine/featscomp.png?branch=master)](https://travis-ci.org/DevMine/featscomp)
[![GoDoc](http://godoc.org/github.com/DevMine/featscomp?status.svg)](http://godoc.org/github.com/DevMine/featscomp)
[![GoWalker](http://img.shields.io/badge/doc-gowalker-blue.svg?style=flat)](https://gowalker.org/github.com/DevMine/featscomp)

`featscomp` computes 'features', as defined by the DevMine project and inserts
them into the database.

## Installation

To install `featscomp`, run this command in a terminal, assuming
[Go](http://golang.org/) is installed:

    go get github.com/DevMine/featscomp

Or you can download a binary for your platform from the DevMine project's
[downloads page](http://devmine.ch/downloads).

You also need to setup a [PostgreSQL](http://www.postgresql.org/) database. Look
at the
[README file](https://github.com/DevMine/featscomp/blob/master/db/README.md)
in the `db` sub-folder for details.

## Usage and configuration

Copy `featscomp.conf.sample` to `featscomp.conf` and edit it according to you
needs. It has two sections, one to configure access to the database, the other
to specify which feature you would like to compute:

 * **database**: allows you to configure access to your PostgreSQL
   database.
   - **hostname**: hostname of the machine.
   - **port**: PostgreSQL port.
   - **username**: PostgreSQL user that has access to the database.
   - **password**: password of the database user.
   - **dbname**: database name.
   - **ssl\_mode**: takes any of these 4 values: "disable",
     "require", "verify-ca", "verify-null". Refer to PostgreSQL
     [documentation](http://www.postgresql.org/docs/9.4/static/libpq-ssl.html)
     for details.
 * **features**: allows you to specify which feature to compute.

Once configuring is done, simply run `featscomp` and give it a configuration
file as argument:

    featscomp featscomp.conf
