package daihinmin

type X string

type WelcomeReply struct {
	X
	Name string
	Sesh sesh
}

type ErrorReply struct {
	X
	Wrt string
	Msg string
}

type YouJoined struct {
	X
	PlayerID int
	Chan     string
}

type UserJoinPartReply struct {
	X
	Chan string
	User string
}

type PlayReply struct {
	X
	Events []Event
	Hand   Cards
}

type GameInfo struct {
	X     `json:",omitempty"`
	ID    string
	Name  string
	Users []string
}

// The status of the current game which is public to all the players.
type GameStatus struct {
	X             `json:",omitempty"`
	CurrentPlayer int
	Stats         []int
	Pile
}

type YourTurn struct {
	X        `json:",omitempty"`
	PlayerID int
	Sesh     sesh
	Hand     Cards
}
