package model

import "time"

type TaskHour struct {
	TaskHourID      int       `db:"TaskHourID"`
	TaskID          int       `db:"TaskID"`
	TaskHourTypeID  int       `db:"TaskHourTypeID"`
	PersonID        int       `db:"PersonID"`
	CreateDate      time.Time `db:"CreateDate"`
	CreatedPersonID int       `db:"CreatedPersonID"`
	TaskHourWorked  time.Time `db:"TaskHourWorked"`
	TaskHourNote    string    `db:"TaskHourNote"`
}
