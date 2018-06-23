# I Love Indexes

## About

This project is about massive information indexation. We wanted to be able to
get data from "Index Of" pages, you know, the default file listing pages of
most web servers. Our goal was to aggregate several of these websites and to make
real time searches in them.

It is hosted here: https://i-love-indexes.scalingo.io. Before making searches and
having fun, please read the following.

## Application structure

The application has the following structure.

### API

The first component of the project (written in Go, of course) is a simple HTTP
API, designed to submit and list "Index Of"pages and overall, to request search
results.

### Front-End

One (little) page application, making async requests to the API. We'd have liked to
write this in Go, but https://github.com/gopherjs/gopherjs is not mature enough to
build something consistent during a 2-days hackathon.

### Crawler

Its job is to crawl the websites received by the API. Both programs are communicating
thank to the https://github.com/jrallison/go-workers package, using a redis backend.
The API send crawling tasks to the Crawler which handles them.

The crawler send indexation data to Elasticsearch in order to be able to query it later.

## Features

What can you do with this stuff?!

* Add an "Index Of"-like website. We are currently supporting (and parsing HTTP
  with https://github.com/PuerkitoBio/goquery) of 3 different web servers:
  * Apache
  * Nginx
  * Lighttpd 

  → Example: http://index.l3o.eu:8080/ (it has already been crawled, too bad,
  you'll get a validation error if you try to add it.)

* Website is crawled recursively and titles of every file is indexed (after
  sanitizing for Elasticsearch, because the way ES tokenize strings can be a
  bit tricky)

* When you click in the "list" icon on the website, you get the list of all the
  indexed websites, and the count of indexed files.

* Writing in the search field trigger a fuzzy search among all the indexed
  content and return in a paginated way the results.

* You can also use regexp over the filenames, try `.*\.mkv` for instance

* If you are looking for a particular type of files, we are making category
  sorting: `audio`, `video`, `ebook`, it will automatically restrain your
  search scope to this category

* For MP3s, we download the first KB of the file then break the connection
  with the webserver to reduce the bandwidth consumption.

* For all the network related operation, there is a threashold of the maximal
  number of connections per host, as a result no one get DoS-ed.

## Example

With the original soundtrack of (_Inglorious
Basterds_)[http://index.l3o.eu:8080/Quentin%20Tarantino%27s%20Inglourious%20Basterds%20V0/] 

You can search for `basterds` and find all the different songs of the album,
but as we are doing fuzzy matching, you can make typos and type `bastrd`, it
will still work!

You can have an insight of its power there: https://www.youtube.com/watch?v=WI0RJbco_l4
mp3 headers have been fetched to get their metadata.

## Disclaimer

Of course, we do not encourage piracy, this project has been built for the
technical challenge first, and to have fun.
