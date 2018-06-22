package main

import (
"net"
"fmt"
"os"
"runtime"
"time"
"strconv"
)
func loop(startport, endport int, inport chan int) {
	for i := startport; i <= endport; i++{
		inport <- i
	}
}

func scanner(inport, outport, out chan int, ip net.IP, endport int){
	for{
		in := <- inport
		tcpaddr := &net.TCPAddr{IP:ip,Port:in}
		conn, err := net.DialTCP("tcp", nil, tcpaddr)
		if err != nil {
			outport <- 0
		}else{
			outport <- in
			conn.Close()
		}
		if in == endport{
			out <- in
		}
	}
}

func main() {
	runtime.GOMAXPROCS(4)
	inport := make(chan int)
	starttime := time.Now().Unix()
	outport := make(chan int)
	out := make(chan int)
	collect := []int{}
	if len(os.Args) != 4 {
		fmt.Println("Usage: scanner.exe IP startport endport")
		fmt.Println("Endport must be larger than startport")
		os.Exit(0)
	}
	ip := net.ParseIP(os.Args[1])
	if(os.Args[3] < os.Args[2]) {
		fmt.Println("Usage: scanner IP startport endport")
		fmt.Println("Endport must be larger than startport")
		os.Exit(0)
	}
	fmt.Printf("the ip is %s \r\n",ip)
	startport, _:= strconv.Atoi(os.Args[2])
	endport, _ := strconv.Atoi(os.Args[3])
	fmt.Printf("%d---------%d\r\n",startport,endport)
	go loop(startport, endport, inport)
	for{
		select {
		case <-out:
			fmt.Println(collect)
			endtime := time.Now().Unix()
			fmt.Println("The scan process has spent ",endtime - starttime,"second")
			os.Exit(0)
		default:
			go scanner(inport, outport, out, ip, endport)
			port := <- outport
			if port != 0{
				collect = append(collect, port)
			}
		}
	}
}


