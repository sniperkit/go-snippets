TypingGeeks
================

## Introduction

TypingGeeks is a peer-to-peer typing game based on Terminal targeting on people who love & always work with Terminal.
The game can be portable anywhere that have Terminal whether on your local machine, enterprise server, or your server on the cloud but not only limited to these.
It will also be designed to be modular and communication via network procotol that anyone will be able to design and create a new TypingGeeks game on any platform such as native app or web service with high graphics. That means a player who plays on Terminal will be able to compete with another player on web browser w/o any difference.

With this mini-game, developer can use it to relax during work 5-10 mins and gain more efficiency after playing it.

## Concept

- Vs. based game via P2P or socket network protocol
- Each player has HP, MP and skills they can use to attach the other player.
- MP will be increased after successfully finish a word
- Only one game, one stage, play until the other side dies with HP zero.


## Game Design

### Main Server

- Match users
- Generate text
- Exchange progress of each player and skill that send to attach another player

### Client

- Get texts and random pop-up position
- Render
- Input ( Send all inputs to main server )

Noted: Server will give only text and other necessary communication between 2 clients. This can help reducing lots of network traffics and put all unnecessary calculation to client side. )

#### Implementation

- Thread for display text
- Thread for key inputs
- Thread for receiving data via network or other channel


## How Golang can help

- Go routines for all parts ( Threading becomes easier )
- Modularity, easily implementing plugins and scale
- Socket programming

