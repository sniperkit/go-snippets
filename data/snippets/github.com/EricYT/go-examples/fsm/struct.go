package main

import (
	"fmt"

	"github.com/looplab/fsm"
)

type Door struct {
	To  string
	FSM *fsm.FSM
}

func NewDoor(to string) *Door {
	d := &Door{
		To: to,
	}

	d.FSM = fsm.NewFSM(
		"closed",
		fsm.Events{
			{Name: "open", Src: []string{"closed"}, Dst: "open"},
			{Name: "close", Src: []string{"open"}, Dst: "closed"},
		},
		fsm.Callbacks{
			//			"enter_state":  func(e *fsm.Event) { d.enterState(e) },
			"before_close": func(e *fsm.Event) { fmt.Println("before closed") },
			"enter_closed": func(e *fsm.Event) { d.enterClose(e) },
			//"close":       close,
			"leave_closed": func(e *fsm.Event) { fmt.Println("leave closed") },
			"after_close":  func(e *fsm.Event) { d.afterClose(e) },
		},
	)

	return d
}

func (d *Door) enterState(e *fsm.Event) {
	fmt.Printf("The door to %s is %s\n", d.To, e.Dst)
}

func (d *Door) afterClose(e *fsm.Event) {
	fmt.Printf("after the door to %s is %s\n", d.To, e.Dst)
}

func (d *Door) enterClose(e *fsm.Event) {
	fmt.Printf("enter the door to %s is %s\n", d.To, e.Dst)
}

func close(e *fsm.Event) {
	fmt.Println("colse call")
}

func main() {
	door := NewDoor("heaven")

	err := door.FSM.Event("open")
	if err != nil {
		fmt.Println(err)
	}

	err = door.FSM.Event("close")
	if err != nil {
		fmt.Println(err)
	}

}
