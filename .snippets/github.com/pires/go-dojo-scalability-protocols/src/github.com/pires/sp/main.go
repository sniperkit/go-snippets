package main

import (
	"flag"
	"log"
	"strings"
	"sync"

	"github.com/pires/sp/common"
	"github.com/pires/sp/sp"
)

var (
	ip          = flag.String("ip", "127.0.0.1", "The IP address to bind or connect to. Defaults to 127.0.0.1")
	port        = flag.String("port", "49999", "The IP port to bind or connect to. Defaults to 49999")
	mode        = flag.String("mode", "receiver", "Possible modes: sender, receiver. Defaults to receive")
	numMessages = flag.Uint("num_messages", 1000, "The number of messages to send. Defaults to 1000")

	wg             sync.WaitGroup
	successChan    = make(chan bool, 100)
	successCounter uint
	failureCounter uint
)

func main() {
	flag.Parse()

	*mode = strings.ToLower(*mode)
	if *mode != "sender" && *mode != "receiver" {
		log.Fatalln("Invalid mode [", *mode, "]")
	}

	address := "tcp://" + *ip + ":" + *port

	go count()

	// SP
	if *mode == "sender" {
		log.Println("Sending", *numMessages, "messages to", address)

		// instantiate sender
		spRequester := new(sp.SPSender)

		// instantiate message
		msg := &common.Message{"pires", "jason", "test"}

		// prepare wg and execute
		wg.Add(int(*numMessages))
		for counter := 1; counter <= int(*numMessages); counter++ {
			go func() {
				// remember to notify func is done
				defer wg.Done()

				// send message
				err := spRequester.Send(msg, address)
				successChan <- err == nil
			}()
		}

		// wait for func to send message
		wg.Wait()
	} else if *mode == "receiver" {
		log.Println("Listening for messages at", address)
		spReceiver := new(sp.SPReceiver)
		if err := spReceiver.Receive(address); err != nil {
			log.Fatalln(err)
		}
	}

	// terminate success counters and print results
	close(successChan)
	log.Println("Total successes:", successCounter)
	log.Println("Total failures:", failureCounter)

	log.Println("Exiting..")
}

// count successes and errors
func count() {
	for {
		select {
		case success := <-successChan:
			if success == false {
				failureCounter++
			} else {
				successCounter++
			}
			break
		}
	}
}
