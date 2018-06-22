Matching Snuggies
=================

Matching Snuggies is a slicing software that exposes a backend slicing program
(currently slic3r) through an HTTP API.  A command line slicing tool is
provided for ease of use and to support eventual integration with host software
like Repetier-Host and OctoPrint.

Matching Snuggies is well suited for integration with host-software that may
run in a resource constrained environment, such as a Raspberry Pi.

Documentation
=============

Install
-------

First install [slic3r](http://slic3r.org/download), the backend slicing
software Matching Snuggies has chosen to support initially.

NOTE: OS X users should symlink the executable at Slicer.app/MacOS/slicer into
their environment's PATH.

./build.sh

Slicing API
-----------

A REST API is exposed to schedule slicing jobs, retrieve resulting gcode, and
get periodic status updates while slicing is in progress.

```
./bin/snuggied -slic3r.configs=testdata
```

See the snuggied documentation on
[godoc.org](http://godoc.org/github.com/gophergala/matching-snuggies/cmd/snuggied).
See the API [doc](API.md) for information about each endpoint.

Command line tool
-----------------

After starting, the daemon can be sent files to slice using the command line
tool.

```
./bin/snuggier -preset=hq -o FirstCube.gcode testdata/FirstCube.amf
```

When `snuggied` is running on another host specify the server when calling `snuggier`.

```
./bin/snuggier -server=10.0.10.123:8888 -preset=hq -o FirstCube.gcode testdata/FirstCube.amf
```

See the snuggier command documentation on godoc.org
[godoc.org](http://godoc.org/github.com/gophergala/matching-snuggies/cmd/snuggier).

Long term goals
---------------

- API authorization
- integration with other backend slicers (Cura)
- a slicing queue that may be consumed by a pool of workers (shared
  configuration; dropbox?)
- cluster health/monitoring dashboard
