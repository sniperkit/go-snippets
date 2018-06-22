package model

import "time"

type Task struct {
	TaskID          int       `db:"TaskID"`
	TaskTitle       string    `db:"TaskTitle"`
	TaskStatusID    int       `db:"TaskStatusID"`
	Description     string    `db:"Description"`
	MilestoneID     int       `db:"MilestoneID"`
	CreateDate      time.Time `db:"CreateDate"`
	CreatedPersonID int       `db:"CreatedPersonID"`
	EditedDate      time.Time `db:"EditedDate"`
	EditedPersonID  int       `db:"EditedPersonID"`
	Completed       bool      `db:"Completed"`
}