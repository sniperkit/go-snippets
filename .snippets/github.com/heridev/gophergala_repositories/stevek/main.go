package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/antage/eventsource"
	"github.com/gophergala/stevek/device"
	"github.com/gophergala/stevek/event"
)

const (
	// These are explict rules to most tightly define the initial example
	// Note that as-is they are not completely defined/implementable rules, just snippets
	// These should eventually be read in from files (well, io.Readers perhaps)
	ruleIptables string = `-A INPUT -s 192.0.2.2/32 -p tcp -m tcp --dport 22 -j ACCEPT`
	ruleCiscoAsa string = `access-list 100 extended permit tcp any host 192.0.2.2 eq ssh`
)

func processEvents(es eventsource.EventSource) {
	for {
		event := <-event.Events
		es.SendEventMessage(event.ToJSON(), "", strconv.Itoa(time.Now().Nanosecond()))
	}
}

func main() {

	//Events = make(chan event.Event)
	//var wg sync.WaitGroup

	// serverSetup
	es := eventsource.New(nil, nil)
	defer es.Close()
	http.Handle("/events", es)
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	go processEvents(es)
	//ip, network, err := net.ParseCIDR("192.0.2.2/32")
	//if err != nil {
	//	fmt.Println(err)
	//}
	dev1 := device.New("Generic host", "")
	dev2 := device.New("Cisco ASA", ruleCiscoAsa)
	dev3 := device.New("iptables", ruleIptables)

	// Implement basic chan -> goroutine (noop) -> chan
	// Set up the pipeline and consume the output.
	//h1.link(h2)
	dev1.Transit = device.Gen(21, 22, 44, 222, 555, 2000)
	dev2.Transit = dev2.Filter(dev1.Transit)
	dev3.Transit = dev3.Filter(dev2.Transit)
	//wg.Add(1)
	go sink(dev3)
	//wg.Wait()
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func sink(dev device.Device) {
	//defer wg.Done()
	id := 1
	for msg := range dev.Transit {
		msg.Action = "Sink"
		msg.Device = "Sink" // TODO - make method like Filter, same as need to do with Gen
		event.Events <- event.Event{Data: msg, Type: "Transit event"}
		time.Sleep(2 * time.Second)
		id++
	}
	return
}
