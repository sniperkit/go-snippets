package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/miekg/dns"
	"gopkg.in/yaml.v2"
)

type ednsConfig struct {
	DNSServer   string `json:"dns_server" yaml:"dns_server"`
	Port        string `json:"port" yaml:"port"`
	RedisServer string `json:"redis_server" yaml:"redis_server"`
}

// type Zone struct {
// 	Host    string
// 	Type    string
// 	Records map[string][]string
// }
//
// type ZonesConf struct {
// 	Name  string
// 	Zones []Zone
// }

func handleRequests(w dns.ResponseWriter, r *dns.Msg) {
	fmt.Println("handleRequests")
}

func main() {

	fmt.Println("Starting DNS Server...")

	conf := flag.String("conf", "", "Configuration file")
	local := flag.Bool("local", false, "Use /etc/resolv.conf file")
	flag.Parse()

	if *conf != "" && *local == true {
		fmt.Println("You have to choose one of -local or -conf parameters.")
		return
	}

	var server string
	var port string = "53"

	switch {
	case *local == true:
		c, _ := dns.ClientConfigFromFile("/etc/resolv.conf")
		server = c.Servers[0]
		port = c.Port
		fmt.Println(server, port)
	case conf != nil:
		filename, _ := filepath.Abs(*conf)
		yamlFile, err := ioutil.ReadFile(filename)
		// fmt.Println(yamlFile)
		if err != nil {
			fmt.Println(err)
		}
		var c ednsConfig
		err = yaml.Unmarshal(yamlFile, &c)
		fmt.Println(c)
		if err != nil {
			fmt.Println(err)
		}
		server = c.DNSServer
		if c.Port != "" {
			port = c.Port
		}
	}

	fmt.Println(server, port)

	dnsServer := &dns.Server{Addr: ":53", Net: "udp"}
	dns.HandleFunc(".", handleRequests)
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go dnsServer.ListenAndServe()

forever:
	for {
		select {
		case s := <-sig:
			fmt.Printf("Signal (%d) received, stopping\n", s)
			break forever
		}
	}
	//
	// client := new(dns.Client)
	// m := new(dns.Msg)
	// fmt.Println(flag.Args()[0])
	// m.SetQuestion(dns.Fqdn(flag.Args()[0]), dns.TypeMX)
	// m.RecursionDesired = true
	//
	// r, _, err := client.Exchange(m, net.JoinHostPort(server, port))
	//
	// if r == nil {
	// 	log.Fatalf("*** error: %s\n", err.Error())
	// }
	// if r.Rcode != dns.RcodeSuccess {
	// 	log.Fatalf("*** invalid answer name: %s after MX query for %s\n", os.Args[1], os.Args[1])
	// }
	//
	// for _, a := range r.Answer {
	// 	fmt.Printf("%v\n", a)
	// }
}
