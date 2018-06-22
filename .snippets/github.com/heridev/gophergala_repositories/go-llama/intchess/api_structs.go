package intchess

type APITypeOnly struct {
	Type string `json:"type"`
}

/*
{
"type":"authentication_request",
"username":"test",
"user_token":"test"
}
*/
type APIAuthenticationRequest struct {
	Type      string `json:"type"`
	Username  string `json:"username"`
	UserToken string `json:"user_token"`
}

type APIAuthenticationResponse struct {
	Type     string `json:"type"`
	Response string `json:"response"`
	User     *User  `json:"user,omitempty"` //This will thus be blank (not NULL) if they are not logged in
}

type APIGameRequest struct {
	Type     string `json:"type"`
	Opponent *User  `json:"opponent"`
}

/*
{
"type":"game_response",
"response":"ok"
}
*/
type APIGameResponse struct {
	Type     string `json:"type"`
	Response string `json:"response"`
}

type APIGameOutput struct {
	Type string     `json:"type"`
	Game *ChessGame `json:"game"`
}

/*
{
"type":"signup_request",
"username": "some-other-username",
"user_token": "some-access-token",
"is_ai": false,
"verses_ai": true
}
*/
type APISignupRequest struct {
	Type      string `json:"type"`
	Username  string `json:"username"`
	UserToken string `json:"user_token"`
	IsAi      bool   `json:"is_ai"`
	VersesAi  bool   `json:"verses_ai"`
}

/*
{
"type":"game_move_request",
"move":"e5-e3"
}
*/
type APIGameMoveRequest struct {
	Type string `json:"type"`
	Move string `json:"move"`
}

type APIGameChatRequest struct {
	Type      string `json:"type"`
	MessageId int    `json:"message_id"`
}

type APIGameChat struct {
	Type      string `json:"type"`
	From      *User  `json:"from"`
	MessageId int    `json:"message_id"`
}

/*
{
"type":"game_get_valid_moves_request",
"location":"a2"
}
*/
type APIGameGetValidMovesRequest struct {
	Type     string `json:"type"`
	Location string `json:"location"`
}

type APIGameGetValidMovesResponse struct {
	Type  string   `json:"type"`
	Moves [][]byte `json:"moves"`
}
