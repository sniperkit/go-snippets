// Must install ssh library: 'go get golang.org/x/crypto/ssh'
package main

import (
	"bytes"
	"fmt"
	"log"
	"runtime"
	"time"

	"golang.org/x/crypto/ssh"
)

var bufsize = 4096

type remote struct {
	user     string
	host     string
	password string
	cmd      string
}

type output struct {
	id  string
	raw string
}

var server remote = remote{"tester", "octopus", "funtests", "/home/graham/Programming/go/workspace/bin/pregelTestServer"}

var workers map[string]remote = map[string]remote{
	"id0": remote{"tester", "octopus", "funtests", "/home/graham/Programming/go/workspace/bin/pregelTestWorker"},
}

var client string = "pregelTestClient"

func main() {
	num_remotes := len(workers) + 1 // workers plus server

	stop := make(chan bool, num_remotes)
	results := make(chan output, num_remotes)

	go startServer(stop, results)
	go startWorkers(stop, results)
	clientTime := runClient() // runs in loop til it connects

	// Kill remotes
	for i := 0; i < num_remotes; i++ {
		stop <- true
	}

	// Collect remote outputs
	outputs := map[string]string{}
	for i := 0; i < num_remotes; i++ {
		res := <-results
		outputs[res.id] = res.raw
	}

	analyzeOut(outputs, workers)
	compareResult(clientTime)

	// compare outputs to expected
	// compare result file to expected
}

func runRemoteSSH(r remote, rid string, stop chan bool, results chan output) {
	config := &ssh.ClientConfig{
		User: r.user,
		Auth: []ssh.AuthMethod{
			ssh.Password(r.password),
		},
	}
	client, err := ssh.Dial("tcp", r.host+":22", config)
	checkError(err)
	session, err := client.NewSession()
	checkError(err)
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	err = session.Start(fmt.Sprintf("%s %s", r.cmd, getTestCase(0, rid)))
	checkError(err)

	<-stop
	session.Signal(ssh.SIGINT)
	//checkError(err) // Commented out until we have nonstopping testcode
	results <- output{rid, b.String()}

}

func getTestCase(num int, remoteId string) string {
	// TODO
	if remoteId != "server" {
		return remoteId
	}
	return ""
}

func startServer(stop chan bool, results chan output) {
	runRemoteSSH(server, "server", stop, results)
}

func startWorkers(stop chan bool, results chan output) {
	for wid, worker := range workers {
		go runRemoteSSH(worker, wid, stop, results)
	}
}

func runClient() time.Duration {
	// TODO
	t := 10 * time.Second
	time.Sleep(t)
	return t
}

func analyzeOut(outputs map[string]string, workers map[string]remote) {
	// TODO
	fmt.Println(outputs)
}

func compareResult(clientTime time.Duration) {
	// TODO
	fmt.Printf("Failed to produce valid output in %v\n", clientTime)
}

func checkError(err error) {
	_, _, line, _ := runtime.Caller(1)
	if err != nil {
		log.Fatalf("Line %d: %v", line, err)
	}

}

/*
   - start server and workers in separate goroutines
   - they ssh in and start separate scripts
   - then they wait for a message
   - start client (time it)
   - when client completes, send message to all waiting to retrieve logs and kill
   - compare logs to expected results and report


*/
