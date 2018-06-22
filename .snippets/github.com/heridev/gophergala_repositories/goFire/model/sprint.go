package model

import "time"

type Sprint struct {
	SprintID    int       `db:"SprintID"`
	SprintTitle string    `db:"SprintTitle"`
	Description string    `db:"Description"`
	CompanyID   int       `db:"CompanyID"`
	StartDate   time.Time `db:"StartDate"`
	EndDate     time.Time `db:"EndDate"`
}
