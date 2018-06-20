package scheduler

import "github.com/RichardKnop/machinery/v1/tasks"

func SchedulerJob (jobId string) (*tasks.Signature){
	return &tasks.Signature{
		Name: "consumption",
		Args: []tasks.Arg{
			{
				Type:  "uint",
				Value: jobId,
			},
		},
	}
}
