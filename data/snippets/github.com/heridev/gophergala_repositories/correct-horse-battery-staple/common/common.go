package common

type ClientState struct {
	Id       int64
	Name     string
	Lat, Lng float64
	Accuracy float64
}

type ServerUpdate struct {
	Clients ClientStates
}

type ClientStates []ClientState

func (cs ClientStates) Len() int           { return len(cs) }
func (cs ClientStates) Swap(i, j int)      { cs[i], cs[j] = cs[j], cs[i] }
func (cs ClientStates) Less(i, j int) bool { return cs[i].Id < cs[j].Id }
