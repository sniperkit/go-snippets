# go-plugin

[![GoDoc](https://godoc.org/github.com/ajgb/go-plugin?status.png)][godoc]
[![Build Status](https://travis-ci.org/ajgb/go-plugin.svg?branch=master)][travis]
[![Codecov Status](https://codecov.io/gh/ajgb/go-plugin/branch/master/graph/badge.svg)][codecov]

[travis]: https://travis-ci.org/ajgb/go-plugin
[godoc]: https://godoc.org/github.com/ajgb/go-plugin
[codecov]: https://codecov.io/gh/ajgb/go-plugin

## Description

Package plugin provides a Go library which helps writing monitoring plugins.

## Synopsis

    package main

    import (
      "github.com/ajgb/go-plugin"
    )
    
    var opts struct {
      Hostname string `short:"H" long:"hostname" description:"Host" default:"localhost"`
      Port     int    `short:"p" long:"port" description:"Port" default:"123"`
    }

    func main() {
      check := plugin.New("check_service", "v1.0.0")
      if err := check.ParseArgs(&opts); err != nil {
        check.ExitCritical("Error parsing arguments: %s", err)
      }
      defer check.Final()

      check.AddMessage("Service %s:%d", opts.Hostname, opts.Port)

      serviceMetrics, err := .... // retrieve metrics
      if err != nil {
        check.ExitCritical("Connection to %s:%d failed: %s", opts.Hostname, opts.Port, err)
      }
    
      for metric, value := range serviceMetrics {
        check.AddMetric(metric, value)
      }
    }

## API Documentation

Please read full usage documentation at <https://godoc.org/github.com/ajgb/go-plugin>

## License

Copyright (C) 2017  Alex J. G. Burzyński

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
