package model

import "time"

type TaskHourEstimated struct {
	TaskHourEstimatedID	int			`db:"TaskHourEstimatedID"`
	TaskID 				int			`db:"TaskID"`	
	ProjectedHours  	float32 	`db:"ProjectedHours"`
	RemainingHours 		float32		`db:"RemainingHours"`
	PersonID 			int   		`db:"PersonID"`
	CreateDate   		time.Time 	`db:"CreateDate"`
	CreatedPersonID 	int   		`db:"CreatedPersonID"`
}
