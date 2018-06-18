/*
 * Copyright(C) 2015 Simon Schmidt
 * 
 * This Source Code Form is subject to the terms of the
 * Mozilla Public License, v. 2.0. If a copy of the MPL
 * was not distributed with this file, You can obtain one at
 * http://mozilla.org/MPL/2.0/.
 */

/*
 This package implements the session handling on unix machines.
 */
package unixssh

import "github.com/maxymania/go-system/sshlib"
import "github.com/maxymania/go-system/syscall_x"
import "github.com/kr/pty"

import "os/exec"
import "io"

func handleSessResize(sess *sshlib.ShellSession,fd int, end chan struct{}) {
	for {
		select {
		case <- sess.ChSize: syscall_x.Ioctl_resize(fd,sess.Width,sess.Heigth)
		case <- end: return
		}
	}
}

/*
 Runns a Shell session.
 */
func HandleSess(sess *sshlib.ShellSession, cmd *exec.Cmd) {
	end := make(chan struct{})
	defer close(end)
	defer sess.Ch.Close()
	p,e := pty.Start(cmd)
	if e!=nil { return }
	defer p.Close()
	go io.Copy(p,sess.Ch)
	go io.Copy(sess.Ch,p)
	go handleSessResize(sess,int(p.Fd()),end)
	cmd.Wait()
}


