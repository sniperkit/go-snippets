package model

type Milestone struct {
	MilestoneID    int       `db:"MilestoneID"`
	MilestoneTitle string    `db:"MilestoneTitle"`
	Description    string	 `db:"Description"`
	SprintID 	   int		 `db:"SprintID"`
}
