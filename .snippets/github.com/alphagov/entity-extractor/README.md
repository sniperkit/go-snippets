# Entity Extractor

Service for extracting named entities from text

## Overview/purpose

An entity is any well known real world *thing* such as a person, place,
organisation, policy etc.

A 'named entity' is any such thing which has a well known name, for example
'Government Digital Service'. Many entities will be known by multiple names,
for example Government Digital Service is also known as 'GDS'.

Therefore each entity my be identified by a set of terms, and is ascribed an
identifier.

```
{"terms":["Government digital service","GDS"],"id":"1"}
```

We anticipate building a corpus of known named entities collated from various
sources.

Once we have a well defined corpus of known named entities, we can use that
during document indexing and search.

The entity-extractor service will be used during both of these operations.

When indexing the document, the entity-extractor will be used to find out
which entities are contained in the document. The list of entity ids will be
stored alongside the document.

When searching, the entity-extractor will be used to determine which entities
were mentioned in the search query. These will then be added to the search
command to boost query results which mention the same entity.

## Installation

The `entity-extractor` service can built using make:

```
$ cd entity-extractor
$ make
```

this uses the [`gom`](https://github.com/mattn/gom) dependency manager to
fetch required packages, before building the tool itself.

You can see the list of dependencies in the [Gomfile](Gomfile).

## Running tests

The tests can be run with `make test` which will also fetch dependencies and
compile the package if needed:

```
$ make test
...
gom test -v
=== RUN TestParseEntityFromJson
--- PASS: TestParseEntityFromJson (0.00s)
=== RUN TestLoadEntities
2015/01/15 12:04:41 Loading entities from /Users/davidheath/alphagov/entity-extractor/data/entities.jsonl
2015/01/15 12:04:41 Loaded 4 terms
--- PASS: TestLoadEntities (0.00s)
=== RUN TestExtract
2015/01/15 12:04:41 Loading entities from /Users/davidheath/alphagov/entity-extractor/data/entities.jsonl
2015/01/15 12:04:41 Loaded 4 terms
--- PASS: TestExtract (0.00s)
PASS
ok    _/Users/davidheath/alphagov/entity-extractor  0.012s
```

## Running the service

The service is configured using the following environment variables:

  * `EXTRACTOR_EXTRACT_ADDR` - address which the server listens on. Default `:3096`
  * `EXTRACTOR_ENTITIES_PATH` - path to the [json lines](http://jsonlines.org/) file containing known entities to be loaded when the server starts up. Default `data/entities.jsonl`
  * `EXTRACTOR_LOG_PATH` - file on which to write logging output. Default `STDERR`

To start the service, set any environment variables you may need and then just run
`entity-extractor`:

```
$ ./entity-extractor
2015/01/15 11:39:17 logging JSON to STDERR
2015/01/15 11:39:17 using GOMAXPROCS value of 2
2015/01/15 11:39:17 Loading entities from data/entities.jsonl
2015/01/15 11:39:17 Loaded 4 terms
2015/01/15 11:39:17 listening for requests on :3096
```

## Using the service

The entity-extractor service offers a single API endpoint `POST /extract`.
This accepts a document and returns the entity IDs of entities found in that
document.

```
$ curl -XPOST -d'document mentioning GDS' localhost:3096/extract
["1"]

$ curl -XPOST -d'document mentioning GDS and MoJ' localhost:3096/extract
["1","2"]

$ curl -XPOST -d'document mentioning neither' localhost:3096/extract
[]
```

