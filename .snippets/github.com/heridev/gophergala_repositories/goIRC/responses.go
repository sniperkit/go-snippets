package main

import "fmt"

var canned_responses map[int]string

const HOST_STRING = "goirc.capitalonelabs.com"

const (
	ERR_NOSUCHCHANNEL int = iota
	ERR_NOSUCHNICK
	ERR_UNKNOWNERROR
	ERR_CANNOTSENDTOCHAN
	RPL_WELCOME
	RPL_YOURHOST
	RPL_CREATED
	RPL_MYINFO
	RPL_ISUPPORT
	RPL_YOURID
	RPL_MOTDSTART
	RPL_MOTD
	RPL_ENDOFMOTD
)

func loadMessages() {
	canned_responses = make(map[int]string)

	canned_responses[RPL_WELCOME] = " 001 %s :Welcome to the Capital One Labs IRC Network"
	canned_responses[RPL_YOURHOST] = " 002 %s :Your host is goirc.capitalonelabs.com, running goIRC-0.0.1"
	canned_responses[RPL_CREATED] = " 003 %s :This server was created at some point in the past"
	canned_responses[RPL_MYINFO] = " 004 %s :some server modes go here or something"
	canned_responses[RPL_ISUPPORT] = " 005 %s :info about limits and so env variables will go here"
	canned_responses[RPL_YOURID] = " 006 %s :unique id goes here maybe? (ircnet)"
	canned_responses[RPL_MOTDSTART] = " 372 %s: we don't have an motd yet!!"
	canned_responses[RPL_MOTD] = " 375 %s :" + HOST_STRING + " message of the day"
	canned_responses[RPL_ENDOFMOTD] = " 376 %s :end of motd"
	canned_responses[ERR_UNKNOWNERROR] = " 400 %s : unknown error"
	canned_responses[ERR_NOSUCHNICK] = " 400 %s :no such nick"
	canned_responses[ERR_NOSUCHCHANNEL] = " 403 %s :no such channel"
	canned_responses[ERR_CANNOTSENDTOCHAN] = " 404 %s cannot send to channel"

	for i, v := range canned_responses {
		canned_responses[i] = ":" + HOST_STRING + v
	}
}

func sendWelcome(user *User) {
	user.Write("PING :" + HOST_STRING)
	user.Write(":" + HOST_STRING + " NOTICE Auth :welcome!")
	types := []int{RPL_WELCOME, RPL_CREATED, RPL_YOURHOST, RPL_MYINFO, RPL_ISUPPORT, RPL_YOURID, RPL_MOTDSTART, RPL_MOTD, RPL_ENDOFMOTD}

	for _, val := range types {
		user.Write(fmt.Sprintf(canned_responses[val], user.Nick))
	}
}
