# goIRC

##IRC server written in Go

[IRC Spec](https://tools.ietf.org/html/rfc1459)

###Connection Steps
 1. ```telnet ec2-54-191-196-95.us-west-2.compute.amazonaws.com 3030```

###Commands
 1. ```PASS <User Name>```
 1. ```JOIN #<Channel Name>```
 1. ```PRIVMSG #<Channel Name>:<Message>```
 1. ```HELP```
 1. ```LIST```
 1. ```PART #<Channel Name>```

### To Run Locally
 1. ```go run *[^_t].go```
 2. ```telnet localhost 3030```

