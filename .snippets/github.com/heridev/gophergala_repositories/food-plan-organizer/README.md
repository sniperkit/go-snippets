# Food Plan Organizer

The Food Plan Organizer helps people who cook often by creating per-week grocery lists, as well as
creating detailed calory intake overviews.

See the releases for easy-to-use prepacked archives. Note that the current state requires you to have Atom-Shell installed on your machine.
Detailed instructions are available in the release instructions.

## Architecture

The application consists of three parts:

1. ETL binary which converts [SR27 nutrition data](http://www.ars.usda.gov/Services/docs.htm?docid=24912) into a [SQLite](https://github.com/mattn/go-sqlite3) database.
   The resulting database file is shipped with the application to skip the database creation on the client.
2. A webserver which handles interaction with the SR27 data as well as managing all user input, including recipe management.
3. A single page application written in angularjs, which provides the UI

We're using [atom-shell](https://github.com/atom/atom-shell) to wrap parts 2 & 3 in a executable binary.

## Hacking

Assuming you've got go installed:

1. run the ETL process once to generate a sr27 database
2. run the backend server: `make run_server`
3. install frontend dependencies `npm install && bower install`
4. run the frontend `grunt serve`
5. open the browser and navigate to `http://127.0.0.1:9000`

If you want to develop in the atom-shell you need to uncomment `main.js:47` and comment `main.js:48`

## TODOs

- Makefile to generate atom-shell
- Tests
- Planning UI
