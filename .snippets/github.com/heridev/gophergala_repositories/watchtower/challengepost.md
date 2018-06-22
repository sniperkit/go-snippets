### Introduction & Background
Like many in my circle of friends I enjoy geeking out with embedded systems.

Rather than taking a minute to water the flowers, everyone "knows" it's wiser to spend hours setting up an Arduino controlling a water pump, controlled by a moisture sensor in the flower pot. To complete the circle, one would of course like to get the moisture sensor readings out from the Arduino, so that instead of "if dry then goto pump_water" we can hook it up to the cloud and replace the original, simple algorithm with a deep neural network that will learn how to water plants optimally.... But now I am getting ahead of myself ;).

Now, as you may know, getting an Arduino connected to the network is not so bad. However, the HTTP stack is not exactly up to par with Go's standard library, so services like Firebase or Pusher are a bit too unwieldy. This is where Watchtower comes in.

I'd been thinking about this kind of project for a while and the Gala was a great opportunity to set aside some time to get started on it.

### Why Go?

Go felt like a natural choice for this project for several reasons. A server app like this requires **performance** (if you want to support thousands of devices on your DigitalOcean droplet), **great concurrency** primitives (to scale up) and a **solid standard library**, so one does not have to reinvent the wheel when it comes to e.g HTTP or TCP servers.

Go has all of this, with loads of sugar on top.

### Key features (project goals)

Hoping to get the app finished enough today that all of these will be true, but below are at least my goals for the project.

* **Multiple interfaces**
  * Coding a controller app for your PC? Use the streaming HTTP endpoints
  * Listening in on a channel and storing data on your server to graph later? Use the HTTP async callback endpoints
  * Sending out data from your Arduinos? Use the TCP/IP or UDP endpoints
* High performance
	* The server should handle thousands of subscribers on thousands of channels without problems
* Reliable
	* High test coverage should give users that warm, fuzzy feeling of reliability

### Things I wanted to do but know I won't have time for

* 100% test coverage
* Support for all of the endpoints
* Add Watchtower support to an existing embedded system so I could make a snazzy video
* Back everything with a database (probably [Bolt](https://github.com/boltdb/bolt)) so state is kept when the server goes down