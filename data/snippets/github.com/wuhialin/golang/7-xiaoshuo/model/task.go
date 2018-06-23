package model

type Task struct {
	Id       int
	TypeId   int
	Title    string
	State    int
	CreateAt int64
	StartAt  int64
	EndAt    int64
}
