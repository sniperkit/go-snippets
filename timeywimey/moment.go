package main

import "time"

type Moment struct {
	ObjectType   string    `json:"objectType" bson:"objectType"`
	StartDate    time.Time `json:"start" bson:"startDate"`
	EndDate      time.Time `json:"end" bson:"endDate"`
	Repeating    bool      `json:"repeating" bson:"repeating"`
	Title        string    `json:"title" bson:"title"`
	Summary      string    `json:"summary" bson:"summary"`
	CalendarData string    `json:"calendarData" bson:"calendarData"`
	LastModified time.Time `json:"lastModified" bson:"lastModified"`
}
