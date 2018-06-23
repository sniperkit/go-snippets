/*
 * Copyright(C) 2016 Simon Schmidt
 * 
 * This Source Code Form is subject to the terms of the
 * Mozilla Public License, v. 2.0. If a copy of the MPL
 * was not distributed with this file, You can obtain one at
 * http://mozilla.org/MPL/2.0/.
 */
package main

import "net"

import "golang.org/x/crypto/ssh"

import "fmt"
import "github.com/maxymania/go-system/sshlib"
import "github.com/maxymania/go-system/sshlib/unixssh"

import "os/exec"
import "io/ioutil"

import "github.com/maxymania/go-system/authen"

var S *ssh.ServerConfig

var N = 0

var P *ssh.Permissions

var SC = make(chan *sshlib.ShellSession,100)

func handleSession(sl *sshlib.ShellSession) {
	//unixssh.HandleSess(sl,exec.Command("/bin/bash"))
	unixssh.HandleSess(sl,exec.Command("/bin/su",sl.Permission.CriticalOptions["user"]))
}

func handler(){
	for sc := range SC {
		go handleSession(sc)
	}
}

func handle(nc net.Conn) {
	sshlib.Handle(nc,SC,S)
}

func passwd_auth(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
	e := authen.AuthenticatePassword(conn.User(),password)
	if e!=nil { return nil,e }
	P := new(ssh.Permissions)
	P.CriticalOptions = make(map[string]string)
	P.Extensions = make(map[string]string)
	P.CriticalOptions["user"] = conn.User()
	return P,nil
}

// it is disabled, for now.
func pubk_auth(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
	return P,nil
}

func load_keys() bool {
	keys := []string{
		"rsa.pem",
		"dsa.pem",
		"ecdsa.pem",
	}
	for _,k := range keys {
		b,e := ioutil.ReadFile(k)
		if e!=nil { return true }
		s,e := ssh.ParsePrivateKey(b)
		if e!=nil { return true }
		S.AddHostKey(s)
	}
	return false
}

func main() {
	S = new(ssh.ServerConfig)
	P = new(ssh.Permissions)
	P.CriticalOptions = make(map[string]string)
	P.Extensions = make(map[string]string)
	S.PasswordCallback = passwd_auth
	//S.PublicKeyCallback = pubk_auth
	
	if load_keys() { return }
	
	l,e := net.Listen("tcp",":64022")
	if e!=nil {
		fmt.Println(e)
		return
	}
	go handler()
	for {
		c,e := l.Accept()
		if e!=nil {
			fmt.Println(e)
			break
		}
		go handle(c)
	}
	fmt.Println(l.Close())
}



