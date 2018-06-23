package main

import "errors"

type Location interface {
	GetButtonTitle() string
	GetHeader() string
	GetDescription() string
	Link(Location)
	GetLinkedLocations() []Location
	Enter(*State)
}

type BaseLocation struct {
	Header          string
	ButtonTitle     string
	Description     string
	LinkedLocations []Location
	LinkedActions   []Action
}

func (location *BaseLocation) Link(nextLocation Location) {
	if nextLocation == nil {
		panic(errors.New("NEXTLOCATION = NIL"))
	}

	location.LinkedLocations = append(location.LinkedLocations, nextLocation)
}

func (location *BaseLocation) GetLinkedLocations() []Location {
	return location.LinkedLocations
}

func (location *BaseLocation) LinkAction(nextAction Action) {
	location.LinkedActions = append(location.LinkedActions, nextAction)
}

func (location *BaseLocation) GetLinkedActions() []Action {
	return location.LinkedActions
}

func (location *BaseLocation) GetButtonTitle() string {
	return location.ButtonTitle
}

func (location *BaseLocation) GetHeader() string {
	return location.Header
}

func (location *BaseLocation) GetDescription() string {
	return location.Description
}

func (location *BaseLocation) Enter(*State) {
	// noop
}
