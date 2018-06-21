package taskresult

import (
	"fmt"
)

const (
	PreConditionCheckFailed  = 0
	PostConditionCheckFaield = 1
	MainRoutineFailed        = 2
	ClearRoutineFailed       = 3
)

type Result struct {
	Name        string
	Description string
	Success     bool
	Reason      int
	Message     string
}

func (r *Result) String() string {
	if r.Success {
		return fmt.Sprintf("Task %s has been run successfully!", r.Name)
	}

	res := fmt.Sprintf("Task %s has been failed: \n", r.Name)
	switch r.Reason {
	case PreConditionCheckFailed:
		res += fmt.Sprintf("		PreCondition Failed:\n")
	case PostConditionCheckFaield:
		res += fmt.Sprintf("		PostCondition Failed:\n")
	case MainRoutineFailed:
		res += fmt.Sprintf("		MainRoutine Failed:\n")
	case ClearRoutineFailed:
		res += fmt.Sprintf("		ClearRoutine Failed:\n")
	default:
		res += fmt.Sprintf("		UnknownReason!\n")
	}
	res += fmt.Sprintf("         %s\n", r.Message)

	return res
}
