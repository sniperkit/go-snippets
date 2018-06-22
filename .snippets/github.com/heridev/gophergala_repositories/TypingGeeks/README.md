TypingGeeks
================

![TypingGeeksGif](https://github.com/gophergala/TypingGeeks/blob/master/TypingGeeks.gif)

"The Will-To-Be best Typing Game Ever"

## Introduction

__Noted:__ This is for now a one-player client player which you can only play locally. It is a Proof-of-Concept (POC).  After contest, it will be extended to multiplayer via network where we have main server to match each player together and have skills where each player can attack the other side.

TypingGeeks will be a peer-to-peer typing game based on Terminal targeting on people who love & always work with Terminal.

The game can be portable anywhere that have Terminal whether on your local machine, enterprise server, or your server on the cloud but not only limited to these.

It will also be designed to be modular and communication via network procotol that anyone will be able to design and create a new TypingGeeks game on any platform such as native app or web service with high graphics. That means a player who plays on Terminal will be able to compete with another player on web browser w/o any difference.

With this mini-game, developers can use it to relax during work 5-10 mins and gain more efficiency after playing it.

## Concept

### 1 Player (Current)
- many levels (word velocity, word pop-up frequency)
- HP life decreasing if fail to type or fail to finish word

### Multiplayer (Future)
- Vs. based game via P2P or socket network protocol
- Each player has HP, MP and skills they can use to attack the other player.
- MP will be increased after successfully finish a word
- Only one game, one stage, play until the other side dies with HP zero.

## Game Design

### Main Server (Future)

- User matching
- Generate text ( many languages will be supported, this can increase difficulty and usefulness)
- Exchange progress of each player and skill that send to attack another player

### Client ( Current - Only 1 player locally, Future - Use as a client for Multiplayer )

- Get texts and random pop-up position
- Render
- Input ( Send all inputs to main server )

Noted: Server will give only text and other necessary communication between 2 clients. This can help reducing lots of network traffics and put all unnecessary calculation to client side. )

#### Implementation

- Thread for display text
- Thread for key inputs
- Thread for receiving data via network or other channel (Future)

## How Golang helps

- __Go routines for all parts:__ Threading wouldn't become easier than this!
- __Modularity:__ easily implementing plugins and scalability
- __Utf8 support:__ to support different languages
- __Ease of deployment:__ just one file, you can just put on any server, cloud, or your local machine and PLAY!
- __Socket programming (Future):__ for multiplayer via network via Go Channel
