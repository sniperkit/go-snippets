package main

import (
	"encoding/json"
	"fmt"
	"github.com/gophergala/GopherKombat/common/game"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type ContestantProcess struct {
	dir string
	cmd *exec.Cmd

	stdin  io.WriteCloser
	stdout io.ReadCloser

	encoder *json.Encoder
	decoder *json.Decoder
}

func NewContestantProcess(contestant *Contestant) (*ContestantProcess, error) {
	var err error
	cp := &ContestantProcess{}

	// Create directory and import AI code
	cp.dir, err = ioutil.TempDir("", "sandbox")
	if err != nil {
		return nil, err
	}

	ai := filepath.Join(cp.dir, "main.go")
	if err := ioutil.WriteFile(ai, []byte(contestant.Code), 0400); err != nil {
		return nil, fmt.Errorf("error creating temp file %q: %v", ai, err)
	}

	exe := filepath.Join(cp.dir, "a.out")
	cmd := exec.Command("go", "build", "-o", exe, ai)
	// Not using NaCl for now because it is printing some bytes to the
	// stdout
	//cmd.Env = []string{"GOOS=nacl", "GOARCH=amd64p32", "GOPATH=/go"}
	if out, err := cmd.CombinedOutput(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			// Error compiling AI
			return nil, fmt.Errorf("Error compiling AI: %s", string(out))
		}
		return nil, fmt.Errorf("error building go source: %v", err)
	}

	// Prepare AI to receive requests
	// Not using NaCl for now because it is printing some bytes to the
	// stdout
	//cp.cmd = exec.Command("sel_ldr_x86_64", "-l", "/dev/null", "-S", "-e", exe)
	cp.cmd = exec.Command(exe)
	cp.stdin, err = cp.cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("error opening stdin: %v", err)
	}
	cp.stdout, err = cp.cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("error opening stderr: %v", err)
	}
	cp.encoder = json.NewEncoder(cp.stdin)
	cp.decoder = json.NewDecoder(cp.stdout)
	err = cp.cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("error starting AI: %v", err)
	}

	return cp, nil
}

func (cp *ContestantProcess) Turn(state *game.State) (*game.Action, error) {
	respc := make(chan *game.Action, 1)
	errc := make(chan error, 1)

	go func() {
		// Send state to AI
		err := cp.encoder.Encode(state)
		if err != nil {
			errc <- err
			return
		}

		//buff := make([]byte, 12)
		//cp.stdout.Read(buff)

		// Read action from AI
		var action game.Action
		err = cp.decoder.Decode(&action)
		if err != nil {
			errc <- err
			return
		}
		respc <- &action
		return
	}()

	t := time.NewTimer(time.Second)
	select {
	case err := <-errc:
		t.Stop()
		return nil, err
	case resp := <-respc:
		t.Stop()
		return resp, nil
	case <-t.C:
		cp.cmd.Process.Kill()
		return nil, fmt.Errorf("timeout")
	}
}

func (cp *ContestantProcess) Close() {
	os.RemoveAll(cp.dir)
	cp.cmd.Process.Kill()
}
