===== Install =====

To install the example API please use the follow commands:

	cd static/
	bower install ng-admin --save

	go get github.com/Onefootball/entity-rest-api
	go get github.com/ant0ine/go-json-rest/rest
	go get github.com/mattn/go-sqlite3


===== Run =====

Once you run the API with the following command:

	go run main.go // run on port 8080 at 127.0.0.1

You will be able now to check the admin endpoint at:

	http://127.0.0.1:8080/admin/

Or you can request the following end points:

	GET http://localhost:8080/api/:entity
	POST http://localhost:8080/api/:entity
	GET http://localhost:8080/api/:entity/:id
	PUT http://localhost:8080/api/:entity/:id
	DELETE http://localhost:8080/api/:entity/:id

Where the `entity` parameter is a reflection to the table name. Sample requests:

	GET /user #get all users
	GET /user/1 #get user with id 1

	GET /post/1 #get post with id 1