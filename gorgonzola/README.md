# gorgonzola [![GoDoc](https://godoc.org/github.com/gophergala/gorgonzola?status.svg)](https://godoc.org/github.com/gophergala/gorgonzola)

Gorgonzola is simple barebone job board using the Json-job format to aggregate job offers.

Demo on: http://gorgonzola-gophergala.appspot.com/

Supported features:

- Add jobs by submitting `jobs.json` file. More informatin [here](http://lukasz-madon.github.io/json-job/).
- Json file validation using json-schema.
- Automatic update for job offers
- [Microformats](http://schema.org/JobPosting) for job offers
- [Google App Engine](https://cloud.google.com/appengine/docs) hosted