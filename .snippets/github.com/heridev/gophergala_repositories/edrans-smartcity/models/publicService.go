package models

import (
	"math/rand"
)

const (
	SERVICE_HOSPITAL    = "Hospital"
	SERVICE_POLICE      = "PoliceDept"
	SERVICE_FIREFIGHTER = "FireDept"

	VEHICLE_AMBULANCE  = "ambulance"
	VEHICLE_POLICE_CAR = "police car"
	VEHICLE_PUMPER     = "pumper"

	CALL_SERVICE_MEDIC   = "Medic"
	CALL_SERVICE_POLICE  = "Police"
	CALL_SERVICE_FIREMAN = "Fireman"
)

type PublicService struct {
	Location int //NodeID
	Service  string
	Vehicles []Vehicle
	Errors   chan error
}

//will be used only for patrolmen
func (s *PublicService) readErrors(c *City) {
	for {
		<-s.Errors
		newPatrolman := Vehicle{
			Service:      SERVICE_POLICE,
			MinWeight:    5,
			Alert:        make(chan Path, 2),
			Errors:       s.Errors,
			InCity:       c,
			Position:     c.GetNode(s.Location),
			BasePosition: c.GetNode(s.Location),
		}
		s.Vehicles = append(s.Vehicles, newPatrolman)
		go newPatrolman.patrol(rand.Int() % c.GetNumNodes())
	}
}

func NewPublicServicePosition(city *City, numNodes int) int {
	var pos int
	for {
		pos = rand.Intn(numNodes) + 1
		node := city.GetNode(pos)
		if len(node.Outputs) != 0 && len(node.Sem.Inputs) != 0 {
			break
		}
	}
	return pos
}
