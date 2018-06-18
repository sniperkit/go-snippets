package toolkit

import (
	//"io/ioutil"
	"os/exec"
	//"time"
)

func RunCommand(cmd string, parm ...string) (string, error) {
	cmdOut, err := exec.Command(cmd, parm...).Output()
	if err == nil {
		return string(cmdOut), nil
	} else {
		return "", err
	}
}

type AsyncCommand struct {
	Txt   string
	Parms []string

	Output []byte
	Error  error

	channelInput     chan string
	channelStopInput chan bool
}

/*
func (a *AsyncCommand) Run() {
	var e error

	c := exec.Command(a.Txt, a.Parms...)

	cIn, _ := c.StdinPipe()
	cOut, _ := c.StdoutPipe()

	if a.channelInput == nil {
		a.channelInput = make(chan string)
	}

	if a.channelStopInput == nil {
		a.channelStopInput = make(chan bool)
	}

	a.Output = nil
	a.Error = nil
	e = c.Start()
	if e != nil {
		a.Error = "Unable start " + e.Error()
		return
	}

	// writing to stdin via channel
	go func() {
		closeInput := false
		for !closeInput {
			select {
			case <-a.channelStopInput:
				closeInput = true
				close(a.channelInput)
				cIn.Close()

			case in <- a.channelInput:
				cIn.Write([]byte(in))

			default:
				time.Sleep(1 * time.Millisecond)
			}
		}
	}()

	outs, _ := ioutil.ReadAll(cOut)
	e = c.Wait()
	if e != nil {
		a.Error = "Unable to wait - " + e.Error()
		return
	}

	a.Error = nil
	a.Output = outs
}
*/
