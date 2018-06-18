/*
 * Copyright(C) 2017 Simon Schmidt
 * 
 * This Source Code Form is subject to the terms of the
 * Mozilla Public License, v. 2.0. If a copy of the MPL
 * was not distributed with this file, You can obtain one at
 * http://mozilla.org/MPL/2.0/.
 */
package main

import "os"
import "net"
import "flag"
import "fmt"
import "time"

var is4 = flag.Bool("4",false,"Explicitely use IPv4")
var is6 = flag.Bool("6",false,"Explicitely use IPv6")
var lca = flag.String("l","","Local Address (ip:port)")
var rca = flag.String("r","","Remote Address (ip:port)")
var pdata = flag.String("D","Hello World!","Data to send")

func getNet() string {
	if *is4 { return "udp4" }
	if *is6 { return "udp6" }
	return "udp"
}

func getLocal(n string) *net.UDPAddr {
	s := *lca;
	if s!="" {
		a,_ := net.ResolveUDPAddr(n,s)
		return a
	}
	return nil
}

func getRemote(n string) *net.UDPAddr {
	s := *rca;
	if s!="" {
		a,_ := net.ResolveUDPAddr(n,s)
		return a
	}
	es := flag.Args()
	if len(es)>0 {
		a,_ := net.ResolveUDPAddr(n,es[0])
		return a
	}
	flag.PrintDefaults()
	os.Exit(1)
	return nil
}

func main(){
	flag.Parse()
	data := []byte(*pdata)
	n := getNet()
	conn,err := net.DialUDP(n,getLocal(n),getRemote(n))
	if err!=nil {
		fmt.Println(err)
		os.Exit(1)
	}
	tk := time.Tick(time.Second)
	go func(){
		for{
			buf := make([]byte,4500)
			i,e := conn.Read(buf)
			if e!=nil { i=0 }
			if i>0 {
				fmt.Println("received ",i," bytes")
			}
		}
	}()
	for d := range tk {
		fmt.Println("sending data(",d,")")
		conn.Write(data)
	}
}

