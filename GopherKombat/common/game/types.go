package game

type ActionCode int

const (
	Nop   ActionCode = 0
	Move  ActionCode = 1
	Shoot ActionCode = 2
)

type GopherData struct {
	X      int  `json:"x"`
	Y      int  `json:"y"`
	Friend bool `json:"friend"`
}

type State struct {
	GopherId int          `json:"gid"`
	Health   int          `json:"health"`
	Ammo     int          `json:"ammo"`
	Me       GopherData   `json:"me"`
	Nearby   []GopherData `json:"nearby"`
}

type Action struct {
	Code ActionCode `json:"code"`
	X    int        `json:"x"`
	Y    int        `json:"y"`
}
