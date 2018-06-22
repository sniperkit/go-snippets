package models

import (
	"fmt"
	"time"
)

type Vehicle struct {
	Service      string
	MinWeight    int
	Busy         bool
	Alert        chan Path
	Errors       chan error
	InCity       *City
	Position     *Node
	BasePosition *Node
}

func (v *Vehicle) patrol(start int) {
	patrol := time.After(1 * time.Millisecond)
	for {
		select {
		case <-patrol:
			node := v.InCity.GetNode(start)
			if node == nil {
				v.Errors <- fmt.Errorf("can not go on patrol")
				return
			}
			if len(node.Outputs) == 0 {
				v.Errors <- fmt.Errorf("can not go on patrol")
				return
			}
			v.Position = node
			patrol = time.After(time.Duration(node.Outputs[0].Weight) * time.Second)
		HasNext:
			for {
				for i := 0; i < len(node.Outputs); i++ {
					next := v.InCity.GetNode(node.Outputs[i].DestinyID)
					if next.Sem.ActiveInput.Name == node.Outputs[i].Name {
						start = node.Outputs[i].DestinyID
						break HasNext
					}
				}
				time.Sleep(5 * time.Second)
			}
		case path := <-v.Alert:
			v.run(path)
			v.Position = v.InCity.GetNode(path.Links[len(path.Links)-1].DestinyID)
			start = v.Position.ID
		}
	}
}

func (v *Vehicle) wait() {
	for {
		path := <-v.Alert
		v.run(path)
		path = <-v.Alert
		switch v.Service {
		case SERVICE_HOSPITAL:
			v.run(path)
		case SERVICE_FIREFIGHTER:
			v.back(path)
		}
	}
}

func (v *Vehicle) run(path Path) time.Duration {
	v.Busy = true
	now := time.Now()
	var i int
	for i = 0; i < len(path.Links); i++ {
		v.InCity.GetNode(path.Links[i].DestinyID).Sem.Status <- SemRequest{Status: true, Allow: path.Links[i].Name}
		time.Sleep(3 * time.Second)
		v.InCity.GetNode(path.Links[i].DestinyID).Sem.Status <- SemRequest{Status: false, Allow: path.Links[i].Name}
		v.Position = v.InCity.GetNode(path.Links[i].DestinyID)
	}
	v.Busy = false
	return time.Since(now)
}

func (v *Vehicle) back(path Path) time.Duration {
	v.Busy = true
	now := time.Now()
	var i int
	for i = 0; i < len(path.Links); i++ {
		time.Sleep(time.Duration(path.Weights[i]) * time.Second)
		v.Position = v.InCity.GetNode(path.Links[i].DestinyID)
	}
	v.Busy = false
	return time.Since(now)
}
