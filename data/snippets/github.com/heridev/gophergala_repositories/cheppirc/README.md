Web-based IRC client in GO. You can use it as a BNC or as a Web IRC service for channels that are important for you.

Story
=====

I use IRC a lot. I have become more and more frustrated about my IRC client disconnecting when my notebook goes to sleep after a couple of minutes without activity to save battery time. The standard answer for my problem would have been a BNC.

However, a BNC requires a VPS and some hacking around even when I install it via apt. I have tried to find a web based IRC software that I can install on my server, but I could not find any that also acts as a BNC and does not drop the IRC connection when I close the browser.

Thus one of my colleauges and I decided to write one, and Go seems to be the best tool for the job. It has good support for web stuff, it's concurrent and it's pretty fast.

CheppIRC is not only a good alternative to a personal BNC, but online communities can also benefit from it. With CheppIRC integrated into a website, it can make it easier for the users to participate in communication on the community's IRC channel.

The status of the project is far from finished. It's rather a prototype then a complete product. We plan to implement the missing IRC features, and replace the session handling with pluggable user-management that allows CheppIRC to be integrated into existing websites.

Github: https://github.com/gophergala/cheppirc

Members: 
Denes Lados https://github.com/mimrock
Daniel Kalmar https://github.com/danikalmar

Demo: http://cheppirc.com:8081/login (Works most of the time:) )

Install
=======

- Clone the repository
- Install go and the dependencies of cheppirc
- go run main.go
- Now Cheppirc listens on port 8081 (harcoded, but easy to change).
- Visit localhost:8081/login to login and use the site.
