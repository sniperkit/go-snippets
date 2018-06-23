package main

import (
	"errors"
	"log"

	"github.com/coreos/go-workflow"
)

func main() {
	log.Println("test work-flow")
	var testVars [4]bool

	one := &workflow.Step{
		Label: "modify testVar 1",
		Run: func(c workflow.Context) error {
			testVars[1] = true
			return nil
		},
	}

	three := &workflow.Step{
		Label: "modify testVar 3",
		Run: func(c workflow.Context) error {
			testVars[3] = true
			return nil
		},
	}

	two := &workflow.Step{
		Label:     "modify testVar 2",
		DependsOn: []*workflow.Step{three},
		Run: func(c workflow.Context) error {
			log.Println("step 2")
			return errors.New("xxxx")
			if !testVars[3] {
			}
			testVars[2] = true
			return nil
		},
	}

	base := &workflow.Step{
		Label:     "modify testVar 0",
		DependsOn: []*workflow.Step{one, two},
		Run: func(c workflow.Context) error {
			if !testVars[1] || !testVars[2] {
			}
			testVars[0] = true
			return nil
		},
	}

	w := workflow.New()
	w.Start = base
	w.OnFailure = workflow.RetryFailure(3)
	err := w.Run()
	if err != nil {
		log.Println(err)
	}
}
