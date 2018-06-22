# opentweet
This project is a few things, see `opentweet.png` for a diagram

1. Most importantly, a new internet application layer protocol for fetching tweets. This allows a standard for distributed fetching of tweets from any opentweet server. See `opentweet-spec.txt`
2. A reference implementation of that protocol, as a golang package in the `protocol/` dir
3. A basic server to serve the tweets, allows registering and posting over REST.
4. An HTML5 client, in Chrome app and Cordova app form, for fetching opentweets and posting to a server. 
5. A test client.

## SECURITY WARNING!!!

Opentweet is not secure, use a throw-away password with any server using this code. Passwords are transmitted in plain text and not hashed on the server.

## Getting Started

If you have the ability to run your own server, please do so! The point of the project is de-centralization.

Otherwise, please feel free to use my server `opentweet.cloudapp.net`

## Running your own server

git and docker is needed

````
git clone https://github.com/gophergala/opentweet.git
cd opentweet
docker build -t opentweet .
docker run --name mysql_vol mysql echo hello
docker run --name mysql -e MYSQL_ROOT_PASSWORD=mysecretpassword -d --volumes-from=mysql_vol mysql
docker run --name opentweet -d --link mysql:mysql -p 12315:12315 -p 80:8080 opentweet
````

## Loading the Chrome App

bower is needed http://bower.io/
````
git clone https://github.com/gophergala/opentweet.git
cd opentweet/client/chromeapp
bower install
````
* In Chrome, go to `chrome://extensions/`
* Check "Developer mode"
* Click "Load unpacked extension..."
* Select "opentweet/client/chromeapp" from where you cloned opentweet

Or, instead of using bower, check the releases and download a zip from after 'bower install'

Some notes on usage. Type the whole server name in your settings: `http://opentweet.cloudapp.net:80` Also, some icons are not showing up for me. Click around the bottom and top right to find all the tabs/pages.

## Registering and posting tweets with Curl

````
curl -i -d '{"user":"me","password":"secret"}' opentweet.cloudapp.net/users
curl -i -d '{"tweet":"Hello Opentweet!!"}' me:secret@opentweet.cloudapp.net/tweet
````

## Building and using the golang test client

setup a go environment

````
go get github.com/gophergala/opentweet
go install github.com/gophergala/opentweet/test-client
$GOHOME/bin/test-client opentweet.cloudapp.net jeff
````


## List of known servers

`opentweet.cloudapp.net`

## People to follow

`jeff@opentweet.cloudapp.net`

Please submit PRs for servers / people
