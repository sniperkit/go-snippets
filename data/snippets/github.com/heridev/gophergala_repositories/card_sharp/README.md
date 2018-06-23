# CardParty.Today

Play card with your friends, without having to purchase games ahead
of time. Just grab your smartphone, start a game on the TV, and play.


## Current Status

A demo site is available at CardParty.Today, where you can create games
and play them. Winning is really a feature, but you can play until then.
Two games are available, a fill in the phrase and a adjective + noun game.


### Future Directions (that might happen, but don't exist yet)

* User created games
* Addon decks to add to normal decks
* Full size decks instead of demo sized decks
* Winning a game
* Drop-in players mid-game
* Deactivate players that have left
* Casting the games to TV's via Chromecast
* Refactoring the code
* Angular.js for the interfaces
* Properly sized buttons and text for phones
* Animations in the UI
* Mobile Apps

### Deploying

The way I deploy is to 
* build the server binary and place into a bin directory
* zip the project up
* sftp the zip file to the web server
* unzip the file into a directory in /opt
* fix public/js/app.js websocket links to point to correct server
* add an upstart conf for the server
* start the server using upstart



### /etc/init/cards.conf

    description "cardparty.today"
    author "Andrew Sellers <andrew@andrewcsellers.com>"

    start on runlevel [2345]
    stop on runlevel [016]

    respawn
    chdir /opt/cards
    env ADMIN_PASS="fwoeuiafnvwuiewahfiwjfa"
    exec /opt/cards/bin/server

