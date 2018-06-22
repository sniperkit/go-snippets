# Realtime tweet map

This project streams tweets with location data from Twitter and broadcasts
them to the listening browsers through websockets. The browsers then plot
markers to tweet location on the Raphael based world map.

Demo available at [tweetmap.korri.fi](http://tweetmap.korri.fi/)

## Running

To get this running a CREDENTIALS file is needed. This file needs to contain
Twitter application credentials. If you have a Twitter development account
you can get that info from the Twitter development site. The file needs to
contain four rows with the following data:

1. Consumer Key
2. Consumer Secret
3. Access Token
4. Access Token Secret

Once you have the CREDENTIALS file in place, you can start the server with:

```
$ bin/bitodd
```

## Libraries used

### Go

* [github.com/gorilla/mux](https://github.com/gorilla/mux)
* [github.com/gorilla/websocket](https://github.com/gorilla/websocket)
* [github.com/darkhelmet/twitterstream](https://github.com/darkhelmet/twitterstream)

### JS

* [JQuery](https://jquery.com/)
* [Raphaël](http://raphaeljs.com/)

### CSS

* [Skeleton](http://getskeleton.com/)
* [Normalize](http://git.io/normalize)