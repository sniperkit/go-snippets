entity-rest-api
===============

[![Build Status](https://travis-ci.org/Onefootball/entity-rest-api.svg?branch=feature%2Fadd-travis-file)](https://travis-ci.org/Onefootball/entity-rest-api)

Description
-----------

entity-rest-api aims to offer a abstract layer over the `go-json-rest` repository, allowing your project to easily query using REST standards to do it.

Installation
------------

This package can be installed with the go get command:

    go get github.com/Onefootball/entity-rest-api

    # dependencies
    go get github.com/ant0ine/go-json-rest/rest

Usage
-----

Into you project, import the following repositories:

	import (
		"github.com/ant0ine/go-json-rest/rest"
		eram "github.com/Onefootball/entity-rest-api/manager"
		era "github.com/Onefootball/entity-rest-api/api"
		"net/http"
	)

Setup the default API REST with `go-json-rest`

	api := rest.NewApi()
	api.Use(rest.DefaultProdStack...)

You will need a database driver, and do the following:

	entityManager := eram.NewEntityDbManager(db)
	entityRestApi := era.NewEntityRestAPI(entityManager)

Then you must setup a router if you want to request something:

	router, err := rest.MakeRouter(
		rest.Get("/api/:entity", entityRestApi.GetAllEntities),
		rest.Post("/api/:entity", entityRestApi.PostEntity),
		rest.Get("/api/:entity/:id", entityRestApi.GetEntity),
		rest.Put("/api/:entity/:id", entityRestApi.PutEntity),
		rest.Delete("/api/:entity/:id", entityRestApi.DeleteEntity),
	)

Finally bind the router to the API and the API to the http handler:

	api.SetApp(router)
	http.Handle("/api/", api.MakeHandler())

Endpoints
---------

Once you run your project based on this library and the example above, you will be able to request the following end points painless:

	GET http://localhost:8080/api/:entity
	POST http://localhost:8080/api/:entity
	GET http://localhost:8080/api/:entity/:id
	PUT http://localhost:8080/api/:entity/:id
	DELETE http://localhost:8080/api/:entity/:id

Where the `entity` parameter is a reflection to the table name. Sample requests:

	GET /user #get all users
	GET /user/1 #get user with id 1

It also allow the insertion of entities based on json.

Queries
-------

Beside the default entity structure, you can do increment your request with queryStrings that allow to order, filter, partition, and more with the entity set.

	_perPage // if you want to use pagination
	_page // current page
	_sortField // the field to sort the query
	_sortDir // the direction of the sort

All the remaining parameters passed by queryString will be treated as filters, for example:

	name=test

This will search by `test` in the column `name` of the entity table.

Tests
-----

To run the tests that right now cover mostly the `api` features do:

	cd api
	go test

License
-------

Copyright (c) 2015 Onefootball GmbH

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
