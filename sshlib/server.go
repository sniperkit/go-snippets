/*
 * Copyright(C) 2015 Simon Schmidt
 * 
 * This Source Code Form is subject to the terms of the
 * Mozilla Public License, v. 2.0. If a copy of the MPL
 * was not distributed with this file, You can obtain one at
 * http://mozilla.org/MPL/2.0/.
 */

/*
 The package "sshlib" is a simple library that makes it easier to work with
 the "golang.org/x/crypto/ssh"-package.
 */
package sshlib

import "net"
import "golang.org/x/crypto/ssh"

import "encoding/binary"

// A shell session.
type ShellSession struct{
	Ch          ssh.Channel
	Permission  *ssh.Permissions
	Term        string
	Width       int
	Heigth      int
	// emits a value, if with or height changes.
	ChSize      <- chan int
	chs         chan int
}
func (s *ShellSession) init() {
	s.chs = make(chan int,1)
	s.ChSize = s.chs
}
func (s *ShellSession) signal() {
	select {
	case s.chs <- 1:
	default:
	}
}

func read32(b []byte) ([]byte,int) {
	if len(b)<4 { return nil,0 }
	i := int(binary.BigEndian.Uint32(b))
	if i<0 { i=0 }
	return b[4:],i
}
func nString(b []byte, n int) ([]byte,string) {
	return b[n:],string(b[:n])
}

func handlePCR_R(n int,r <-chan *ssh.Request,ses ssh.Channel, sc chan *ShellSession,
				p *ssh.Permissions) {
	s := new(ShellSession)
	s.Ch = ses
	s.Permission = p
	s.init()
	defer func() { if ses!=nil { ses.Close() } }()
	for re := range r {
		switch re.Type {
		case "shell":
			if sc!=nil {
				sc <- s
				sc = nil
			}
			if re.WantReply { re.Reply(true,nil) }
		case "pty-req":{
				buf,tl := read32(re.Payload)
				buf,s.Term   = nString(buf,tl)
				buf,s.Width  = read32(buf)
				buf,s.Heigth = read32(buf)
				if re.WantReply { re.Reply(true,nil) }
				s.signal()
			}
		case "window-change":{
				buf := re.Payload
				buf,s.Width  = read32(buf)
				buf,s.Heigth = read32(buf)
				if re.WantReply { re.Reply(true,nil) }
				s.signal()
			}
		default:
			if re.WantReply { re.Reply(false,nil) }
		}
	}
}

func handleCR(n int,ncr <-chan ssh.NewChannel, sc chan *ShellSession, p *ssh.Permissions) {
	for re := range ncr {
		switch re.ChannelType() {
		case "session":
			c,r,e := re.Accept()
			if e!=nil { continue }
			go handlePCR_R(n,r,c,sc,p)
			continue
		}
		//fmt.Println(n,": connected with type",re.ChannelType())
		re.Reject(ssh.UnknownChannelType,"Sorry!")
	}
}
func handleR(n int,r <-chan *ssh.Request) {
	for re := range r {
		//fmt.Println(n,": request",re.Type,re.Payload)
		if re.WantReply { re.Reply(false,nil) }
	}
}

// Handles the SSH protocol and feeds the shell-sessions to sc.
func Handle(nc net.Conn, sc chan *ShellSession, S *ssh.ServerConfig) {
	defer nc.Close()
	s,ncr,r,e := ssh.NewServerConn(nc,S)
	n := 0
	if e!=nil { return }
	go handleCR(n,ncr,sc,s.Permissions)
	go handleR(n,r)
	s.Wait()
	s.Close()
}




