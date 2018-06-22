package main

// peers holds the set of peer servers (in "host:port" format). You don't have
// to use this variable to store the peers, but if you store peers in a
// different way, you'll have to modify the tests (because they modify peers
// during tests).
var peers = make(map[string]struct{})
