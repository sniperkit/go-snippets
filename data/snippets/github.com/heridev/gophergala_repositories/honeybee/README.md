# honeybee

(proof of concept)

This is a service to monitor all your TCP connections between all your servers.

It is compound of an agent service and a server.

First of all you have to setup and run your server instance, after that you have to put and agent (need root permissions for get all TCP ports) in each server you want to monitor.

Then if you go to `localhost:3007` you can see all TCP connections as a link graph