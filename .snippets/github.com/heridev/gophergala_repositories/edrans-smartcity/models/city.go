package models

import (
	"fmt"
	"math/rand"
)

const (
	POS_SOUTH = 0
	POS_NORTH = 1
	POS_WEST  = 2
	POS_EAST  = 3
)

type City struct {
	nodes    []Node
	Services []PublicService
	Name     string
	Size     []int // height x width

	// Aux vars
	LastError error
}

type Node struct {
	ID       int
	Location []int
	Outputs  []Link
	Sem      *Semaphore
}

type Link struct {
	Name      string
	OriginID  int
	DestinyID int
	Weight    int
}

type Path struct {
	Links            []Link
	Weights          []int
	Estimate         int
	OriginalEstimate int
	Reached          bool
	ForgetMe         bool
}

type Location struct {
	Lat     int
	Long    int
	Vehicle int //-1: none, 0: police, 1: ambulance, 2:pumper
	Input   int
	Weight  int
}

func NewCity(nodeList []Node, name string, height, width int) (city *City, err error) {
	if height <= 0 || width <= 0 {
		err = fmt.Errorf("Trying to create an invalid-sized city")
		return nil, err
	}

	myCity := City{
		nodes: nodeList,
		Name:  name,
		Size:  []int{height, width},
	}
	myCity.generateSem()
	myCity.enableSem()

	return &myCity, nil
}

func (c *City) GetNumNodes() int {
	return len(c.nodes)
}

func (c *City) AddService(service string, location, vehicles, minWeight int) {
	var newservice PublicService
	newservice.Service = service
	newservice.Location = location
	newservice.Errors = make(chan error, 5)
	node := c.GetNode(location)
	if node == nil {
		fmt.Println("invalid node:", location)
		return
	}
	for i := 0; i < vehicles; i++ {
		newservice.Vehicles = append(
			newservice.Vehicles,
			Vehicle{
				Service:      service,
				MinWeight:    minWeight,
				Errors:       newservice.Errors,
				InCity:       c,
				BasePosition: c.GetNode(location),
				Position:     c.GetNode(location),
				Alert:        make(chan Path, 5),
			})
	}
	c.Services = append(c.Services, newservice)
}

func (c *City) LaunchVehicles() {
	for i := 0; i < len(c.Services); i++ {
		if c.Services[i].Service == SERVICE_HOSPITAL || c.Services[i].Service == SERVICE_FIREFIGHTER {
			for j := 0; j < len(c.Services[i].Vehicles); j++ {
				go c.Services[i].Vehicles[j].wait()
			}
		} else {
			go c.Services[i].readErrors(c)
			for j := 0; j < len(c.Services[i].Vehicles); j++ {
				go c.Services[i].Vehicles[j].patrol(rand.Int() % len(c.nodes))
			}
		}
	}
}

func (c *City) CallService(call string) (*Vehicle, error) {
	switch call {
	case CALL_SERVICE_MEDIC:
		return c.callService(SERVICE_HOSPITAL, VEHICLE_AMBULANCE)
	case CALL_SERVICE_FIREMAN:
		return c.callService(SERVICE_FIREFIGHTER, VEHICLE_PUMPER)
	case CALL_SERVICE_POLICE:
		return c.callService(SERVICE_POLICE, VEHICLE_POLICE_CAR)
	}
	return nil, fmt.Errorf("unknown service")
}

func (c *City) callService(service, name string) (*Vehicle, error) {
	var base PublicService
	for i := 0; i < len(c.Services); i++ {
		if c.Services[i].Service == service {
			base = c.Services[i]
		}
	}
	for i := 0; i < len(base.Vehicles); i++ {
		if !base.Vehicles[i].Busy {
			return &base.Vehicles[i], nil
		}
	}
	return nil, fmt.Errorf("There is no %s available", name)
}

func (c *City) enableSem() {
	for i := 0; i < len(c.nodes); i++ {
		c.nodes[i].Sem.Status <- SemRequest{Status: false}
	}
}

func (c *City) getLinked(node int) []Link {
	var links []Link
	for i := 0; i < len(c.nodes); i++ {
		for j := 0; j < len(c.nodes[i].Outputs); j++ {
			if c.nodes[i].Outputs[j].DestinyID == node {
				links = append(links, c.nodes[i].Outputs[j])
			}
		}
	}
	return links
}

func (c *City) generateSem() {
	for i := 0; i < len(c.nodes); i++ {
		links := c.getLinked(c.nodes[i].ID)
		if len(links) == 0 {
			newDefault := defaultSemaphore()
			c.nodes[i].Sem = &newDefault
			continue
		}
		var sem Semaphore
		sem.Interval = defaultInterval
		sem.Inputs = links
		sem.ActiveInput = &sem.Inputs[0]
		sem.Status = make(chan SemRequest, 1)
		c.nodes[i].Sem = &sem
		go sem.Start()
	}
}

func (c *City) GetNode(ID int) *Node {
	if c.LastError != nil {
		return nil
	}
	if len(c.nodes) < ID || ID <= 0 {
		c.LastError = fmt.Errorf("Node %d does not exist", ID)
		return nil
	}
	return &c.nodes[ID-1]
}

func (c *City) GetLocations() []Location {
	var locations = make([]Location, len(c.nodes))
	for i := 0; i < len(locations); i++ {
		locations[i].Lat = c.nodes[i].Location[0]
		locations[i].Long = c.nodes[i].Location[1]
		locations[i].Vehicle = c.getVehicle(c.nodes[i].ID)
		if c.nodes[i].Sem == nil {
			locations[i].Input = -1
			continue
		}
		if len(c.nodes[i].Sem.Inputs) == 0 {
			locations[i].Input = -1
			continue
		}

		input := c.GetNode(c.nodes[i].Sem.ActiveInput.OriginID)
		if input == nil {
			fmt.Printf("Error in city:198\n")
			return nil
		}
		locations[i].Weight = c.nodes[i].Sem.ActiveInput.Weight
		inputNode := input.Location
		switch {
		case inputNode[1] > c.nodes[i].Location[1]:
			locations[i].Input = POS_SOUTH
		case inputNode[1] < c.nodes[i].Location[1]:
			locations[i].Input = POS_NORTH
		case inputNode[0] > c.nodes[i].Location[0]:
			locations[i].Input = POS_WEST
		case inputNode[0] < c.nodes[i].Location[0]:
			locations[i].Input = POS_EAST
		}
	}
	return locations
}

func (c *City) getVehicle(node int) int {
	i := 0
	j := 0
	defer func() {
		_ = recover()
	}()
	for i = 0; i < len(c.Services); i++ {
		for j = 0; j < len(c.Services[i].Vehicles); j++ {
			if c.Services[i].Vehicles[j].Position.ID == node {
				var vehicleType int
				switch c.Services[i].Vehicles[j].Service {
				case SERVICE_HOSPITAL:
					vehicleType = 1
				case SERVICE_FIREFIGHTER:
					vehicleType = 2
				case SERVICE_POLICE:
					vehicleType = 0
				}
				return vehicleType
			}
		}
	}
	return -1
}

func (c *City) CleanError() {
	c.LastError = nil
}
