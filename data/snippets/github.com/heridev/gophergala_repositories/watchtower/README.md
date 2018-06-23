## Watchtower

Watchtower is a service for easily interacting between many different kinds of devices in (almost) real-time. Think https://pusher.com but for devices ranging from PC’s all the way down to Arduinos with only the barest of network stacks.

It provides **channels** you can communicate on, which carry events from senders to any amount of receivers, so your devices can talk to each other in any kind of setup you choose. It supports endpoints ranging from UDP, TCP all the way up to HTTP - making it easy to choose the right one.

You get freedom to choose the format of your messages. You like JSON? You get JSON! You like binary? You get binary!

### How should you use Watchtower?

	* You download it and set up a server
	* All your flower-watering Arduinos register and joins a channel
	* Your Macbook joins the same channel
	* Your Macbook starts playing Flight of the Valkyries on max volume
	* One by one, the Arduinos start spraying water all over
	* Your wife kicks you out for ruining the curtains for the umpteenth time
	* You are now living on the street, trading Go Programming for food (but you have an awesome flower watering system, if only one of those VC's would invest...)

### So how do I actually use it?

	* Register as a new user
	* Join some channels
	* Broadcast some messages
	* Act on / display received messages

The rest, as they say, is history.

### Ideas

	* Connect all your embedded devices to different channels - you now have a communication protocol that ties everything together
	* Use it as an extremely lightweight IRC server, for chatting
	* Use it to launch nuclear weapons, but require two or more registered users to confirm the launch

### Message formats

#### TCP/IP & UDP

##### Register as a user

Send 'R' to the server, like so

	echo "R" | nc localhost 3033
	
The server will respond with an acknowledgement packet (first byte = 'A') and your four-byte user ID (second to fifth byte)

##### Join a channel

Send 'J' (one byte) plus your user ID (four bytes) and the channel ID to the server, for example (if your user ID is 5308416 and you want to join channel 5308411)

	echo "J00000001" | nc localhost 3033
	
The server will respond with an ack packet (first byte = 'A'), and the channel ID (following four bytes)

##### Send messages

Send 'M', followed by your user ID (four bytes), the channel ID (four bytes), the length of the message (four bytes) and the actual message (N bytes)

Since you must be subscribed to the channels you send to, you will get the same message back on the channel, in the message format described below.

##### Receive messages

A message is denoted by the first byte being set to 'M'. The following four bytes are the sender, then then channel ID (4 bytes), the message type (2 bytes - not yet in use) and the message length (4 bytes). After this, the message itself follows.

A client will receive all messages (for now, private messages in the pipeline) sent to any channel the client is listening to.

### Implementing your own client

Take a look in `interfaces/tcp/handler_test.go` for a code sample on how the system can be used in Go. Adapt to suit whatever environment you are in.

