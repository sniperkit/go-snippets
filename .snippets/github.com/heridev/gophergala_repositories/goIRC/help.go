package main

type HelpCommand struct {
	Summary string
	Syntax  string
	// Help		string - more detailed help instructions
	// Ops		boolean - flag to see who can use the command
}

var Help = map[string]HelpCommand{
	"JOIN": HelpCommand{
		Summary: "join a new channel",
		Syntax:  "JOIN <#channel>",
	},
	"PRIVMSG": HelpCommand{
		Summary: "send a message to a channel",
		Syntax:  "PRIVMSG <#channel>: <your message here>",
	},
	"NICK": HelpCommand{
		Summary: "change your nickname",
		Syntax:  "NICK <new_nick_here>",
	},
	"PASS": HelpCommand{
		Summary: "~*~hidden command~*~ short hand for registering and setting a username",
		Syntax:  "PASS <username>",
	},
	"TOPIC": HelpCommand{
		Summary: "see the topic for a certain channel",
		Syntax:  "TOPIC <#channel>: <your topic here>",
	},
	"LIST": HelpCommand{
		Summary: "lists out all available channels to join",
		Syntax:  "LIST",
	},
	// "USERS": HelpCommand{
	// 	Summary: "get a list of all users in a channel",
	// 	Syntax:  "syntax: USERS <#channel>",
	// },
	// "NICK": HelpCommand{
	// 	Summary: "change your nickname",
	// 	Syntax:  "NICK <new_nick_here>",
	// },
	// "QUIT": HelpCommand{
	// 	Summary: "quit the server",
	// 	Syntax:  "<QUIT>",
	// },
	// "PART": HelpCommand{
	// 	Summary: "leave the channel",
	// 	Syntax:  "PART <#channel>",
	// },
	// "KICK": HelpCommand{
	// 	Summary: "kick out a user from the channel - ops only",
	// 	Syntax:  "KICK <#channel>: <user_to_kick>",
	// },
	// "INVITE": HelpCommand{
	// 	Summary: "invite a user to the channel",
	// 	Syntax:  "INVITE <#channel>: <user_to_invite>",
	// },
	// "KILL": HelpCommand{
	// 	Summary: "disconnect a user - ops only",
	// 	Syntax:  "KILL <user_to_kill>",
	// },
}
